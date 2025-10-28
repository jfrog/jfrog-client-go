package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jfrog/gofrog/parallel"

	"github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	httputils "github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DirectDownloadService struct {
	client      *jfroghttpclient.JfrogHttpClient
	Progress    clientio.ProgressMgr
	artDetails  *auth.ServiceDetails
	DryRun      bool
	Threads     int
	saveSummary bool
	// A ContentWriter of FileTransferDetails structs. Used only if saveSummary is set to true.
	filesTransfersWriter *content.ContentWriter
	// A ContentWriter of ArtifactDetails structs. Used only if saveSummary is set to true.
	artifactsDetailsWriter *content.ContentWriter
}

func NewDirectDownloadService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *DirectDownloadService {
	return &DirectDownloadService{artDetails: &artDetails, client: client}
}

func (dds *DirectDownloadService) GetArtifactoryDetails() auth.ServiceDetails {
	return *dds.artDetails
}

func (dds *DirectDownloadService) IsDryRun() bool {
	return dds.DryRun
}

func (dds *DirectDownloadService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return dds.client
}

func (dds *DirectDownloadService) GetThreads() int {
	return dds.Threads
}

func (dds *DirectDownloadService) SetThreads(threads int) {
	dds.Threads = threads
}

func (dds *DirectDownloadService) SetDryRun(isDryRun bool) {
	dds.DryRun = isDryRun
}

func (dds *DirectDownloadService) SetSaveSummary(saveSummary bool) {
	dds.saveSummary = saveSummary
}

func (dds *DirectDownloadService) getOperationSummary(totalSucceeded, totalFailed int) *utils.OperationSummary {
	operationSummary := &utils.OperationSummary{
		TotalSucceeded: totalSucceeded,
		TotalFailed:    totalFailed,
	}
	if dds.saveSummary {
		operationSummary.TransferDetailsReader = content.NewContentReader(dds.filesTransfersWriter.GetFilePath(), content.DefaultKey)
		operationSummary.ArtifactsDetailsReader = content.NewContentReader(dds.artifactsDetailsWriter.GetFilePath(), content.DefaultKey)
	}
	return operationSummary
}

// DirectDownloadFiles downloads files using the direct download API (main entry point)
func (dds *DirectDownloadService) DirectDownloadFiles(downloadParams ...DirectDownloadParams) (int, int, error) {
	summary, err := dds.directDownloadFiles(downloadParams...)
	if err != nil {
		return 0, 0, err
	}
	return summary.TotalSucceeded, summary.TotalFailed, nil
}

// DirectDownloadFilesWithSummary downloads files and returns detailed summary
func (dds *DirectDownloadService) DirectDownloadFilesWithSummary(downloadParams ...DirectDownloadParams) (*utils.OperationSummary, error) {
	dds.SetSaveSummary(true)
	return dds.directDownloadFiles(downloadParams...)
}

// directDownloadFiles is the main download logic - follows the same pattern as DownloadFiles
func (dds *DirectDownloadService) directDownloadFiles(downloadParams ...DirectDownloadParams) (operationSummary *utils.OperationSummary, err error) {
	producerConsumer := parallel.NewRunner(dds.GetThreads(), 20000, false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	expectedChan := make(chan int, 1)
	successCounters := make([]int, dds.GetThreads())

	if dds.saveSummary {
		dds.filesTransfersWriter, err = content.NewContentWriter(content.DefaultKey, true, false)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = errors.Join(err, dds.filesTransfersWriter.Close())
		}()
		dds.artifactsDetailsWriter, err = content.NewContentWriter(content.DefaultKey, true, false)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = errors.Join(err, dds.artifactsDetailsWriter.Close())
		}()
	}

	dds.prepareTasks(producerConsumer, expectedChan, successCounters, errorsQueue, downloadParams...)

	err = dds.performTasks(producerConsumer, errorsQueue)
	totalSuccess := 0
	for _, v := range successCounters {
		totalSuccess += v
	}
	operationSummary = dds.getOperationSummary(totalSuccess, <-expectedChan-totalSuccess)
	return
}

// prepareTasks prepares download tasks for parallel execution
func (dds *DirectDownloadService) prepareTasks(producer parallel.Runner, expectedChan chan int, successCounters []int, errorsQueue *clientutils.ErrorsQueue, downloadParamsSlice ...DirectDownloadParams) {
	go func() {
		defer producer.Done()
		defer close(expectedChan)
		totalTasks := 0
		defer func() {
			expectedChan <- totalTasks
		}()

		// Iterate over download params and produce tasks
		for _, downloadParams := range downloadParamsSlice {
			// Handle build-based downloads
			if downloadParams.Build != "" {
				tasks, err := dds.createBuildDownloadTasks(downloadParams, successCounters, errorsQueue)
				if err != nil {
					log.Error(err)
					errorsQueue.AddError(err)
					continue
				}
				totalTasks += dds.produceTasks(tasks, producer, errorsQueue)
				continue
			}

			// Handle regular pattern-based downloads
			tasks, err := dds.createPatternDownloadTasks(downloadParams, successCounters, errorsQueue)
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
				continue
			}
			totalTasks += dds.produceTasks(tasks, producer, errorsQueue)
		}
	}()
}

// performTasks executes all tasks in parallel
func (dds *DirectDownloadService) performTasks(consumer parallel.Runner, errorsQueue *clientutils.ErrorsQueue) error {
	// Blocked until finish consuming
	consumer.Run()
	return errorsQueue.GetError()
}

// produceTasks adds tasks to the producer
func (dds *DirectDownloadService) produceTasks(tasks []parallel.TaskFunc, producer parallel.Runner, errorsQueue *clientutils.ErrorsQueue) int {
	count := 0
	for _, task := range tasks {
		_, err := producer.AddTaskWithError(task, errorsQueue.AddError)
		if err != nil {
			errorsQueue.AddError(err)
		} else {
			count++
		}
	}
	return count
}

// createBuildDownloadTasks creates download tasks for build artifacts
func (dds *DirectDownloadService) createBuildDownloadTasks(params DirectDownloadParams, successCounters []int, errorsQueue *clientutils.ErrorsQueue) ([]parallel.TaskFunc, error) {
	// When build is specified without flags, download all artifacts (not dependencies)
	if !params.IsExcludeArtifacts() && !params.IsIncludeDeps() {
		params.ExcludeArtifacts = false
		params.IncludeDeps = false
	}

	artifacts, err := dds.getArtifactsFromBuild(&params)
	if err != nil {
		return nil, err
	}

	log.Debug(fmt.Sprintf("Found %d artifacts from build %s", len(artifacts), params.Build))

	var tasks []parallel.TaskFunc
	for _, artifactPath := range artifacts {
		// Parse the artifact path
		parts := strings.SplitN(artifactPath, "/", 2)
		var repo, path string
		if len(parts) == 2 {
			repo = parts[0]
			path = parts[1]
		} else {
			// If no repo in path, use the pattern's repo
			patternRepo, _, _ := dds.parsePattern(params.GetPattern())
			repo = patternRepo
			path = artifactPath
		}

		if repo == "" || path == "" {
			log.Warn("Skipping invalid artifact path:", artifactPath)
			continue
		}

		// Check exclusions
		if dds.isExcluded(path, params.GetExclusions()) {
			log.Debug("Artifact excluded by pattern:", path)
			continue
		}

		// Create download task
		task := dds.createSingleDownloadTask(repo, path, &params, successCounters)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// createPatternDownloadTasks creates download tasks for pattern-based downloads
func (dds *DirectDownloadService) createPatternDownloadTasks(params DirectDownloadParams, successCounters []int, errorsQueue *clientutils.ErrorsQueue) ([]parallel.TaskFunc, error) {
	repo, artifactPath, err := dds.parsePattern(params.GetPattern())
	if err != nil {
		return nil, err
	}

	var tasks []parallel.TaskFunc

	// Check if pattern ends with "/" or if we should treat it as a directory
	isDirectoryPattern := strings.HasSuffix(artifactPath, "/") || strings.HasSuffix(params.GetPattern(), "/")

	if isDirectoryPattern {
		// Handle directory download
		// Remove trailing slash for directory operations
		dirPath := strings.TrimSuffix(artifactPath, "/")

		// List files in the directory
		filesToDownload, err := dds.getFilesFromDirectory(repo, dirPath, &params)
		if err != nil {
			return nil, err
		}

		for _, filePath := range filesToDownload {
			if !dds.isExcluded(filePath, params.GetExclusions()) {
				task := dds.createSingleDownloadTask(repo, filePath, &params, successCounters)
				tasks = append(tasks, task)
			}
		}
	} else if dds.containsWildcards(artifactPath) {
		// Handle wildcard patterns
		filesToDownload, err := dds.getFilesMatchingPattern(repo, artifactPath, &params)
		if err != nil {
			return nil, err
		}

		for _, filePath := range filesToDownload {
			task := dds.createSingleDownloadTask(repo, filePath, &params, successCounters)
			tasks = append(tasks, task)
		}
	} else {
		// Single file download
		if !dds.isExcluded(artifactPath, params.GetExclusions()) {
			task := dds.createSingleDownloadTask(repo, artifactPath, &params, successCounters)
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// createSingleDownloadTask creates a task function for downloading a single file
func (dds *DirectDownloadService) createSingleDownloadTask(repo, artifactPath string, params *DirectDownloadParams, successCounters []int) parallel.TaskFunc {
	return func(threadId int) error {
		success, err := dds.downloadSingleFile(repo, artifactPath, params)
		if err != nil {
			return err
		}
		if success {
			successCounters[threadId]++
		}
		return nil
	}
}

// getFilesFromDirectory returns all files in a directory based on recursive flag
func (dds *DirectDownloadService) getFilesFromDirectory(repo, dirPath string, params *DirectDownloadParams) ([]string, error) {
	var filesToDownload []string

	if params.IsRecursive() {
		// For recursive downloads, collect all files in subdirectories
		err := dds.collectAllFilesRecursively(repo, dirPath, params, &filesToDownload)
		return filesToDownload, err
	}

	// For non-recursive, just get files in the immediate directory
	files, err := dds.listDirectoryFiles(repo, dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name)
		filesToDownload = append(filesToDownload, filePath)
	}

	return filesToDownload, nil
}

// collectAllFilesRecursively collects all files in a directory and its subdirectories
func (dds *DirectDownloadService) collectAllFilesRecursively(repo, basePath string, params *DirectDownloadParams, result *[]string) error {
	// Stack for iterative directory traversal
	dirsToProcess := []string{basePath}

	for len(dirsToProcess) > 0 {
		currentDir := dirsToProcess[len(dirsToProcess)-1]
		dirsToProcess = dirsToProcess[:len(dirsToProcess)-1]

		items, err := dds.listDirectoryItems(repo, currentDir)
		if err != nil {
			log.Error("Failed to list directory:", currentDir, err)
			continue
		}

		for _, item := range items {
			itemPath := filepath.Join(currentDir, item.Name)
			if item.Folder {
				// Add subdirectory to process
				dirsToProcess = append(dirsToProcess, itemPath)
			} else {
				// Add file to results
				*result = append(*result, itemPath)
			}
		}
	}

	return nil
}

// getFilesMatchingPattern returns all files matching the given wildcard pattern
func (dds *DirectDownloadService) getFilesMatchingPattern(repo, pattern string, params *DirectDownloadParams) ([]string, error) {
	var filesToDownload []string

	if params.IsRecursive() && filepath.Dir(pattern) != "." {
		// Recursive search
		err := dds.collectFilesRecursively(repo, filepath.Dir(pattern), filepath.Base(pattern), params, &filesToDownload)
		return filesToDownload, err
	}

	// Non-recursive search
	dir := filepath.Dir(pattern)
	filePattern := filepath.Base(pattern)

	files, err := dds.listDirectoryFiles(repo, dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if matched, _ := filepath.Match(filePattern, file.Name); matched {
			filePath := filepath.Join(dir, file.Name)
			if !dds.isExcluded(filePath, params.GetExclusions()) {
				filesToDownload = append(filesToDownload, filePath)
			}
		}
	}

	return filesToDownload, nil
}

// collectFilesRecursively collects files matching pattern recursively
func (dds *DirectDownloadService) collectFilesRecursively(repo, basePath, filePattern string, params *DirectDownloadParams, result *[]string) error {
	// Stack for iterative directory traversal
	dirsToProcess := []string{basePath}

	for len(dirsToProcess) > 0 {
		currentDir := dirsToProcess[len(dirsToProcess)-1]
		dirsToProcess = dirsToProcess[:len(dirsToProcess)-1]

		items, err := dds.listDirectoryItems(repo, currentDir)
		if err != nil {
			log.Error("Failed to list directory:", currentDir, err)
			continue
		}

		for _, item := range items {
			if item.Folder {
				// Add subdirectory to process
				dirsToProcess = append(dirsToProcess, filepath.Join(currentDir, item.Name))
			} else {
				// Check if file matches pattern
				if matched, _ := filepath.Match(filePattern, item.Name); matched {
					filePath := filepath.Join(currentDir, item.Name)
					if !dds.isExcluded(filePath, params.GetExclusions()) {
						*result = append(*result, filePath)
					}
				}
			}
		}
	}

	return nil
}

// DirectoryItem represents a file or folder in directory listing
type DirectoryItem struct {
	Name   string
	Folder bool
}

// listDirectoryItems lists all items (files and folders) in a directory
func (dds *DirectDownloadService) listDirectoryItems(repo, path string) ([]DirectoryItem, error) {
	storagePath := fmt.Sprintf("api/storage/%s/%s", repo, path)
	listUrl, err := clientutils.BuildUrl((*dds.artDetails).GetUrl(), storagePath, make(map[string]string))
	if err != nil {
		return nil, err
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(listUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckErrorf("Failed to list directory %s: Status %d", listUrl, resp.StatusCode)
	}

	var storageInfo struct {
		Children []struct {
			Uri    string `json:"uri"`
			Folder bool   `json:"folder"`
		} `json:"children"`
	}

	if err := json.Unmarshal(body, &storageInfo); err != nil {
		return nil, err
	}

	var items []DirectoryItem
	for _, child := range storageInfo.Children {
		items = append(items, DirectoryItem{
			Name:   strings.TrimPrefix(child.Uri, "/"),
			Folder: child.Folder,
		})
	}

	return items, nil
}

// listDirectoryFiles lists only files in a directory (not folders)
func (dds *DirectDownloadService) listDirectoryFiles(repo, path string) ([]DirectoryItem, error) {
	items, err := dds.listDirectoryItems(repo, path)
	if err != nil {
		return nil, err
	}

	var files []DirectoryItem
	for _, item := range items {
		if !item.Folder {
			files = append(files, item)
		}
	}

	return files, nil
}

// downloadSingleFile downloads a single file from Artifactory
func (dds *DirectDownloadService) downloadSingleFile(repo, artifactPath string, params *DirectDownloadParams) (bool, error) {
	downloadPath := fmt.Sprintf("%s/%s", repo, artifactPath)
	downloadUrl, err := clientutils.BuildUrl((*dds.artDetails).GetUrl(), downloadPath, make(map[string]string))
	if err != nil {
		return false, err
	}

	targetPath := params.GetTarget()
	if targetPath == "" {
		targetPath = "./"
	}

	var localPath string
	if params.IsFlat() {
		localPath = filepath.Join(targetPath, filepath.Base(artifactPath))
	} else {
		localPath = filepath.Join(targetPath, artifactPath)
	}

	if dds.DryRun {
		log.Info("[Dry run] Would download:", downloadUrl, "to", localPath)
		return true, nil
	}

	// Create local directory if needed
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return false, errorutils.CheckError(err)
	}

	// Perform the download
	resp, err := dds.performFileDownload(downloadUrl, localPath, params)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return false, errorutils.CheckErrorf("Failed to download %s: HTTP %d", downloadUrl, resp.StatusCode)
	}

	// Handle post-download operations
	if err := dds.handlePostDownload(localPath, params, resp); err != nil {
		return false, err
	}

	// Save summary if needed
	if dds.saveSummary {
		dds.saveDownloadSummary(downloadUrl, localPath, repo, artifactPath, resp)
	}

	log.Info("Downloaded:", downloadUrl, "to", localPath)
	return true, nil
}

// performFileDownload handles the actual file download with split support
func (dds *DirectDownloadService) performFileDownload(downloadUrl, localPath string, params *DirectDownloadParams) (*http.Response, error) {
	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()

	// Check if we should use concurrent download
	shouldUseConcurrent, fileSize := dds.shouldUseConcurrentDownload(downloadUrl, params, &httpClientsDetails)

	if shouldUseConcurrent {
		// Use concurrent download for large files
		return dds.downloadFileConcurrently(downloadUrl, localPath, fileSize, params, &httpClientsDetails)
	}

	// Use regular download for small files
	return dds.downloadFileRegularly(downloadUrl, localPath, params, &httpClientsDetails)
}

// shouldUseConcurrentDownload determines if concurrent download should be used
func (dds *DirectDownloadService) shouldUseConcurrentDownload(downloadUrl string, params *DirectDownloadParams, httpClientsDetails *httputils.HttpClientDetails) (bool, int64) {
	if params.SplitCount <= 0 || params.MinSplitSize < 0 || dds.DryRun {
		return false, 0
	}

	// Get file details
	fileDetails, resp, err := dds.client.GetRemoteFileDetails(downloadUrl, httpClientsDetails)
	if err != nil || fileDetails == nil {
		// Try storage API as fallback
		if fileInfo, err := dds.getFileInfo(downloadUrl); err == nil && fileInfo != nil && fileInfo.Size != "" {
			if size, err := strconv.ParseInt(fileInfo.Size, 10, 64); err == nil {
				return size > params.MinSplitSize*1024*1024, size
			}
		}
		return false, 0
	}

	// Check if server supports range requests
	acceptsRange := resp.Header.Get("Accept-Ranges") == "bytes"
	if !acceptsRange {
		return false, 0
	}

	return fileDetails.Size > params.MinSplitSize*1024*1024, fileDetails.Size
}

// downloadFileConcurrently downloads a file using concurrent chunks
func (dds *DirectDownloadService) downloadFileConcurrently(downloadUrl, localPath string, fileSize int64, params *DirectDownloadParams, httpClientsDetails *httputils.HttpClientDetails) (*http.Response, error) {
	log.Info(fmt.Sprintf("Downloading %s (%s) using %d parallel chunks",
		filepath.Base(localPath),
		formatFileSize(fileSize),
		params.SplitCount))

	concurrentDownloadFlags := httpclient.ConcurrentDownloadFlags{
		DownloadPath:  downloadUrl,
		FileName:      filepath.Base(localPath),
		LocalPath:     filepath.Dir(localPath),
		LocalFileName: filepath.Base(localPath),
		FileSize:      fileSize,
		SplitCount:    params.SplitCount,
		SkipChecksum:  params.IsSkipChecksum(),
	}

	return dds.client.DownloadFileConcurrently(concurrentDownloadFlags, "", httpClientsDetails, dds.Progress)
}

// downloadFileRegularly downloads a file in a single stream
func (dds *DirectDownloadService) downloadFileRegularly(downloadUrl, localPath string, params *DirectDownloadParams, httpClientsDetails *httputils.HttpClientDetails) (*http.Response, error) {
	downloadFileDetails := &httpclient.DownloadFileDetails{
		DownloadPath:  downloadUrl,
		LocalPath:     filepath.Dir(localPath),
		LocalFileName: filepath.Base(localPath),
		SkipChecksum:  params.IsSkipChecksum(),
	}

	// The 4th parameter is isExplode, 5th is bypassArchiveInspection
	return dds.client.DownloadFile(downloadFileDetails, "", httpClientsDetails, params.IsExplode(), params.IsBypassArchiveInspection())
}

// handlePostDownload handles post-download operations like checksum validation, symlinks, and archive extraction
func (dds *DirectDownloadService) handlePostDownload(localPath string, params *DirectDownloadParams, resp *http.Response) error {
	// Validate checksums if needed
	if !params.IsSkipChecksum() && resp != nil {
		if err := dds.validateChecksumFromHeaders(localPath, resp); err != nil {
			log.Warn("Checksum validation failed for", localPath, ":", err)
		}
	}

	// Handle symlinks
	if params.IsSymlink() {
		if err := dds.handleSymlink(localPath, params); err != nil {
			return err
		}
	}

	// Extract archive if needed
	if params.IsExplode() {
		if err := dds.extractArchive(localPath, params); err != nil {
			return err
		}
	}

	return nil
}

// validateChecksumFromHeaders validates file checksums using response headers
func (dds *DirectDownloadService) validateChecksumFromHeaders(localPath string, resp *http.Response) error {
	md5 := resp.Header.Get("X-Checksum-Md5")
	sha1 := resp.Header.Get("X-Checksum-Sha1")
	sha256 := resp.Header.Get("X-Checksum-Sha256")

	if md5 == "" && sha1 == "" && sha256 == "" {
		return nil // No checksums to validate
	}

	fileInfo := &utils.FileInfo{
		Checksums: struct {
			Sha1   string `json:"sha1,omitempty"`
			Sha256 string `json:"sha256,omitempty"`
			Md5    string `json:"md5,omitempty"`
		}{
			Md5:    md5,
			Sha1:   sha1,
			Sha256: sha256,
		},
	}

	return dds.validateChecksum(fileInfo, localPath)
}

// handleSymlink handles symlink creation if the downloaded file is a symlink placeholder
func (dds *DirectDownloadService) handleSymlink(localPath string, params *DirectDownloadParams) error {
	content, err := os.ReadFile(localPath)
	if err != nil || len(content) == 0 {
		return nil
	}

	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "symlink:") {
		return nil
	}

	// This is a symlink placeholder
	target := strings.TrimSpace(strings.TrimPrefix(contentStr, "symlink:"))
	log.Debug("Detected symlink placeholder. Target:", target)

	cleanTarget := filepath.Clean(target)
	// Prevent path traversal attacks
	if strings.Contains(cleanTarget, "..") {
		return errorutils.CheckErrorf("Security: Symlink target contains path traversal: %s", target)
	}

	// Validate symlink if requested
	if params.ValidateSymlinks() && !fileutils.IsPathExists(target, false) {
		return errorutils.CheckErrorf("Symlink validation failed, target doesn't exist: %s", target)
	}

	// Remove placeholder and create symlink
	if err := os.Remove(localPath); err != nil {
		return errorutils.CheckErrorf("Failed to remove symlink placeholder: %s", err.Error())
	}

	if err := os.Symlink(target, localPath); err != nil {
		return errorutils.CheckErrorf("Failed to create symlink: %s", err.Error())
	}

	log.Info("Created symlink:", localPath, "->", target)
	return nil
}

// extractArchive extracts an archive file
func (dds *DirectDownloadService) extractArchive(localPath string, params *DirectDownloadParams) error {
	localDir := filepath.Dir(localPath)
	localFileName := filepath.Base(localPath)

	if err := clientutils.ExtractArchive(localDir, localFileName, localFileName, "", params.IsBypassArchiveInspection()); err != nil {
		return errorutils.CheckErrorf("Failed to extract archive %s: %s", localPath, err.Error())
	}

	return nil
}

// saveDownloadSummary saves download details for summary reporting
func (dds *DirectDownloadService) saveDownloadSummary(downloadUrl, localPath, repo, artifactPath string, resp *http.Response) {
	if dds.filesTransfersWriter == nil {
		return
	}

	sha256 := ""
	sha1 := ""
	md5 := ""
	if resp != nil && resp.Header != nil {
		sha256 = resp.Header.Get("X-Checksum-Sha256")
		sha1 = resp.Header.Get("X-Checksum-Sha1")
		md5 = resp.Header.Get("X-Checksum-Md5")
	}

	rtUrl := strings.TrimSuffix((*dds.artDetails).GetUrl(), "/")
	sourcePath := downloadUrl
	if strings.HasPrefix(sourcePath, rtUrl) {
		sourcePath = strings.TrimPrefix(sourcePath, rtUrl)
		if !strings.HasPrefix(sourcePath, "/") {
			sourcePath = "/" + sourcePath
		}
	}

	fileTransferDetails := clientutils.FileTransferDetails{
		SourcePath: sourcePath,
		TargetPath: localPath,
		RtUrl:      rtUrl,
		Sha256:     sha256,
	}
	dds.filesTransfersWriter.Write(fileTransferDetails)

	// Also write ArtifactDetails
	if dds.artifactsDetailsWriter != nil {
		artifactoryPath := fmt.Sprintf("%s/%s", repo, artifactPath)
		artifactDetails := utils.ArtifactDetails{
			ArtifactoryPath: artifactoryPath,
			Checksums: entities.Checksum{
				Sha1:   sha1,
				Md5:    md5,
				Sha256: sha256,
			},
		}
		dds.artifactsDetailsWriter.Write(artifactDetails)
	}
}

func (dds *DirectDownloadService) parsePattern(pattern string) (string, string, error) {
	parts := strings.SplitN(pattern, "/", 2)
	if len(parts) < 2 {
		return "", "", errorutils.CheckErrorf("Invalid pattern format: %s. Should be 'repo/path/to/artifact'", pattern)
	}
	return parts[0], parts[1], nil
}

func (dds *DirectDownloadService) containsWildcards(path string) bool {
	return strings.ContainsAny(path, "*?")
}

func (dds *DirectDownloadService) isExcluded(path string, exclusions []string) bool {
	log.Debug(fmt.Sprintf("Checking exclusions for path: %s", path))
	for _, exclusion := range exclusions {
		// Handle ** pattern for recursive directory matching
		if strings.Contains(exclusion, "**") {
			// Convert ** pattern to a more flexible check
			pattern := strings.ReplaceAll(exclusion, "**", "*")
			// Check if any part of the path matches
			pathParts := strings.Split(path, string(filepath.Separator))
			for i := range pathParts {
				subPath := strings.Join(pathParts[i:], string(filepath.Separator))
				if matched, _ := filepath.Match(pattern, subPath); matched {
					log.Debug(fmt.Sprintf("Path %s excluded by pattern %s", path, exclusion))
					return true
				}
			}
		}

		// Check against full path
		if matched, _ := filepath.Match(exclusion, path); matched {
			log.Debug(fmt.Sprintf("Path %s excluded by pattern %s", path, exclusion))
			return true
		}
		// Also check against just the filename
		if matched, _ := filepath.Match(exclusion, filepath.Base(path)); matched {
			log.Debug(fmt.Sprintf("Path %s excluded by pattern %s", path, exclusion))
			return true
		}
	}
	return false
}

// getFileInfo fetches file information from Artifactory Storage API
func (dds *DirectDownloadService) getFileInfo(downloadUrl string) (*utils.FileInfo, error) {
	artUrl := (*dds.artDetails).GetUrl()
	repoPath := strings.TrimPrefix(downloadUrl, artUrl)
	repoPath = strings.TrimPrefix(repoPath, "/")

	storagePath := fmt.Sprintf("api/storage/%s", repoPath)
	storageUrl, err := clientutils.BuildUrl(artUrl, storagePath, make(map[string]string))
	if err != nil {
		return nil, err
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(storageUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckErrorf("Failed to get file info: HTTP %d", resp.StatusCode)
	}

	var fileInfo utils.FileInfo
	if err := json.Unmarshal(body, &fileInfo); err != nil {
		return nil, errorutils.CheckError(err)
	}

	return &fileInfo, nil
}

// validateChecksum validates the downloaded file's checksums
func (dds *DirectDownloadService) validateChecksum(fileInfo *utils.FileInfo, localPath string) error {
	localFileDetails, err := fileutils.GetFileDetails(localPath, true)
	if err != nil {
		return err
	}

	if localFileDetails.Checksum.Md5 != fileInfo.Checksums.Md5 {
		return errorutils.CheckErrorf("MD5 checksum mismatch for %s. Expected: %s, Got: %s",
			localPath, fileInfo.Checksums.Md5, localFileDetails.Checksum.Md5)
	}

	if localFileDetails.Checksum.Sha1 != fileInfo.Checksums.Sha1 {
		return errorutils.CheckErrorf("SHA1 checksum mismatch for %s. Expected: %s, Got: %s",
			localPath, fileInfo.Checksums.Sha1, localFileDetails.Checksum.Sha1)
	}

	if localFileDetails.Checksum.Sha256 != fileInfo.Checksums.Sha256 {
		return errorutils.CheckErrorf("SHA256 checksum mismatch for %s. Expected: %s, Got: %s",
			localPath, fileInfo.Checksums.Sha256, localFileDetails.Checksum.Sha256)
	}

	log.Debug("Checksum validation passed for:", localPath)
	return nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// getBuildInfo fetches build information using the Build API
func (dds *DirectDownloadService) getBuildInfo(buildName, buildNumber, project string) (*entities.BuildInfo, error) {
	// Parse build identifier if needed
	if buildNumber == "" && strings.Contains(buildName, "/") {
		parts := strings.Split(buildName, "/")
		buildName = parts[0]
		if len(parts) > 1 {
			buildNumber = parts[1]
		}
	}

	// Construct build info URL
	buildUrl := fmt.Sprintf("%s/api/build/%s/%s",
		strings.TrimSuffix((*dds.artDetails).GetUrl(), "/"),
		buildName,
		buildNumber)

	if project != "" {
		buildUrl += "?project=" + project
	}

	log.Debug("Fetching build info from:", buildUrl)

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(buildUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckErrorf("Failed to get build info: HTTP %d", resp.StatusCode)
	}

	log.Debug("Build API response status:", resp.StatusCode)

	// Parse the response - the build info is wrapped in a "buildInfo" field
	var buildResponse struct {
		BuildInfo entities.BuildInfo `json:"buildInfo"`
	}
	if err := json.Unmarshal(body, &buildResponse); err != nil {
		return nil, errorutils.CheckError(err)
	}

	log.Debug(fmt.Sprintf("Parsed build info: Name=%s, Number=%s, Modules=%d",
		buildResponse.BuildInfo.Name,
		buildResponse.BuildInfo.Number,
		len(buildResponse.BuildInfo.Modules)))

	return &buildResponse.BuildInfo, nil
}

// getArtifactsFromBuild retrieves artifacts and dependencies from build info
func (dds *DirectDownloadService) getArtifactsFromBuild(params *DirectDownloadParams) ([]string, error) {
	// Extract build name and number
	buildName := params.Build
	buildNumber := ""

	if strings.Contains(buildName, "/") {
		parts := strings.SplitN(buildName, "/", 2)
		buildName = parts[0]
		if len(parts) > 1 {
			buildNumber = parts[1]
		}
	}

	if buildName == "" {
		return nil, errorutils.CheckErrorf("Build name cannot be empty")
	}

	if buildNumber == "" {
		return nil, errorutils.CheckErrorf("Build number is required. Use format: buildName/buildNumber")
	}

	buildInfo, err := dds.getBuildInfo(buildName, buildNumber, params.Project)
	if err != nil {
		return nil, err
	}

	var artifacts []string

	log.Debug(fmt.Sprintf("Build info has %d modules", len(buildInfo.Modules)))

	// Process modules to extract artifacts and dependencies
	for i, module := range buildInfo.Modules {
		log.Debug(fmt.Sprintf("Module %d: ID=%s, Type=%s, Artifacts=%d, Dependencies=%d",
			i, module.Id, module.Type, len(module.Artifacts), len(module.Dependencies)))
		// Include build artifacts unless excluded
		if !params.ExcludeArtifacts && len(module.Artifacts) > 0 {
			for j, artifact := range module.Artifacts {
				// Build the full artifact path including repository
				var fullPath string

				log.Debug(fmt.Sprintf("  Artifact %d: Name=%s, Path=%s, Type=%s, OriginalDeploymentRepo=%s",
					j, artifact.Name, artifact.Path, artifact.Type, artifact.OriginalDeploymentRepo))

				// First, get the basic path
				artifactPath := artifact.Path
				if artifactPath == "" && artifact.Name != "" {
					artifactPath = artifact.Name
				}

				// If we have OriginalDeploymentRepo, prepend it
				if artifact.OriginalDeploymentRepo != "" && artifactPath != "" {
					// Check if path already includes repo
					if !strings.HasPrefix(artifactPath, artifact.OriginalDeploymentRepo+"/") {
						fullPath = artifact.OriginalDeploymentRepo + "/" + artifactPath
					} else {
						fullPath = artifactPath
					}
				} else {
					fullPath = artifactPath
				}

				if fullPath != "" {
					artifacts = append(artifacts, fullPath)
					log.Debug("Added artifact from build:", fullPath)
				} else {
					log.Warn("Skipping artifact with empty path")
				}
			}
		}

		// Include dependencies if requested
		if params.IncludeDeps && len(module.Dependencies) > 0 {
			for _, dep := range module.Dependencies {
				if dep.Id != "" {
					artifacts = append(artifacts, dep.Id)
				}
			}
		}
	}

	if len(artifacts) == 0 {
		log.Warn("No artifacts found in build:", buildName+"/"+buildNumber)
	}

	return artifacts, nil
}

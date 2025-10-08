package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DirectDownloadService struct {
	client               *jfroghttpclient.JfrogHttpClient
	Progress             clientio.ProgressMgr
	artDetails           *auth.ServiceDetails
	DryRun               bool
	Threads              int
	saveSummary          bool
	filesTransfersWriter *content.ContentWriter
}

func NewDirectDownloadService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *DirectDownloadService {
	return &DirectDownloadService{artDetails: &artDetails, client: client}
}

func (dds *DirectDownloadService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return dds.client
}

func (dds *DirectDownloadService) IsDryRun() bool {
	return dds.DryRun
}

func (dds *DirectDownloadService) SetDryRun(isDryRun bool) {
	dds.DryRun = isDryRun
}

func (dds *DirectDownloadService) GetThreads() int {
	return dds.Threads
}

func (dds *DirectDownloadService) SetThreads(threads int) {
	dds.Threads = threads
}

func (dds *DirectDownloadService) SetArtifactoryDetails(artDetails auth.ServiceDetails) {
	dds.artDetails = &artDetails
}

func (dds *DirectDownloadService) SetSaveSummary(saveSummary bool) {
	dds.saveSummary = saveSummary
}

// DirectDownloadFiles downloads files using Artifactory's native resolution order without AQL
func (dds *DirectDownloadService) DirectDownloadFiles(downloadParams ...DirectDownloadParams) (int, int, error) {
	summary, err := dds.performDirectDownload(downloadParams...)
	if err != nil {
		return 0, 0, err
	}
	return summary.TotalSucceeded, summary.TotalFailed, nil
}

// DirectDownloadFilesWithSummary downloads files directly from Artifactory without using AQL
// and gives detailed information about each file transfer.
func (dds *DirectDownloadService) DirectDownloadFilesWithSummary(downloadParams ...DirectDownloadParams) (operationSummary *utils.OperationSummary, err error) {
	return dds.performDirectDownload(downloadParams...)
}

func (dds *DirectDownloadService) performDirectDownload(downloadParams ...DirectDownloadParams) (summary *utils.OperationSummary, err error) {
	summary = &utils.OperationSummary{}

	if dds.saveSummary {
		dds.filesTransfersWriter, err = content.NewContentWriter(content.DefaultKey, true, false)
		if err != nil {
			return nil, err
		}
	}

	for _, params := range downloadParams {
		repo, artifactPath, err := dds.parsePattern(params.GetPattern())
		if err != nil {
			log.Error(err)
			summary.TotalFailed++
			continue
		}

		if dds.isExcluded(artifactPath, params.GetExclusions()) {
			continue
		}

		if dds.containsWildcards(artifactPath) {
			// Handle patterns like "*.zip" or "test?.jar" by listing the directory
			// and downloading matching files one by one
			count, failed, err := dds.handleWildcardDownload(repo, artifactPath, &params)
			summary.TotalSucceeded += count
			summary.TotalFailed += failed
			if err != nil {
				log.Error(err)
			}
		} else {
			// Directly download the file in case wildcards are not present
			success, err := dds.downloadSingleFile(repo, artifactPath, &params)
			if err != nil {
				log.Error(err)
				summary.TotalFailed++
			}

			if success {
				summary.TotalSucceeded++
			} else {
				summary.TotalFailed++
			}
		}
	}

	if dds.saveSummary && dds.filesTransfersWriter != nil {
		if err = dds.filesTransfersWriter.Close(); err != nil {
			return nil, err
		}
		filePath := dds.filesTransfersWriter.GetFilePath()
		log.Debug("Creating content reader from file:", filePath)

		summary.TransferDetailsReader = content.NewContentReader(filePath, content.DefaultKey)
	}

	return summary, nil
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
	for _, exclusion := range exclusions {
		if matched, _ := filepath.Match(exclusion, path); matched {
			log.Debug("Artifact excluded by pattern:", path, "matches", exclusion)
			return true
		}
	}
	return false
}

// downloadSingleFile downloads a single file using artifactory's direct API.
// It validates checksums using response headers which artifactory always provides for downloaded artifacts.
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

	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return false, errorutils.CheckError(err)
	}

	if dds.DryRun {
		log.Info("[Dry run] Would download:", downloadUrl, "to", localPath)
		return true, nil
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()

	// Get file info to determine if we should use concurrent download
	var (
		fileSize     int64
		resp         *http.Response
		acceptsRange bool
	)

	// First, check if we need to get file size for split download decision
	shouldCheckSplit := params.SplitCount > 0 && params.MinSplitSizeMB >= 0
	log.Debug(fmt.Sprintf("Split download check - SplitCount: %d, MinSplitSizeMB: %d", params.SplitCount, params.MinSplitSizeMB))

	if shouldCheckSplit && !dds.DryRun {
		// Try to get file details using existing GetRemoteFileDetails method which uses HEAD request
		log.Debug("Getting file info using HEAD request")
		fileDetails, headResp, err := dds.client.GetRemoteFileDetails(downloadUrl, &httpClientsDetails)
		if err == nil && fileDetails != nil {
			fileSize = fileDetails.Size
			// Check if server supports range requests using existing IsAcceptRanges
			acceptsRange = headResp.Header.Get("Accept-Ranges") == "bytes"
			log.Debug(fmt.Sprintf("File size: %d bytes, MinSplitSize: %d bytes, Accepts ranges: %v",
				fileSize, params.MinSplitSizeMB*1024*1024, acceptsRange))
		} else {
			// Fallback to Storage API if HEAD request fails
			log.Debug("HEAD request failed, falling back to Storage API:", err)
			fileInfo, fsErr := dds.getFileInfo(downloadUrl)
			if fsErr == nil && fileInfo != nil && fileInfo.Size != "" {
				fileSize, fsErr = strconv.ParseInt(fileInfo.Size, 10, 64)
				if fsErr != nil {
					log.Debug("Failed to parse file size:", fsErr)
					fileSize = 0
				} else {
					// Assume range support for Storage API fallback (existing behavior)
					acceptsRange = true
					log.Debug(fmt.Sprintf("File size from Storage API: %d bytes", fileSize))
				}
			}
		}
	}

	// Decide whether to use concurrent download based on file size, parameters, and server capabilities
	useConcurrentDownload := shouldCheckSplit && fileSize > 0 && fileSize > params.MinSplitSizeMB*1024*1024 && acceptsRange
	log.Debug(fmt.Sprintf("Use concurrent download decision: %v (shouldCheckSplit: %v, fileSize: %d, threshold: %d, acceptsRange: %v)",
		useConcurrentDownload, shouldCheckSplit, fileSize, params.MinSplitSizeMB*1024*1024, acceptsRange))

	if useConcurrentDownload {
		// Use concurrent download for large files
		log.Debug(fmt.Sprintf("Using concurrent download for %s (size: %d bytes, split count: %d)", downloadUrl, fileSize, params.SplitCount))
		// Show user-friendly progress message
		log.Info(fmt.Sprintf("Downloading %s (%s) using %d parallel chunks",
			filepath.Base(localPath),
			formatFileSize(fileSize),
			params.SplitCount))

		concurrentDownloadFlags := httpclient.ConcurrentDownloadFlags{
			DownloadPath:  downloadUrl,
			FileName:      filepath.Base(localPath),
			LocalPath:     filepath.Dir(localPath),
			LocalFileName: filepath.Base(localPath),
			RelativePath:  artifactPath,
			FileSize:      fileSize,
			SplitCount:    params.SplitCount,
			SkipChecksum:  params.IsSkipChecksum(),
		}

		resp, err = dds.client.DownloadFileConcurrently(concurrentDownloadFlags, "", &httpClientsDetails, dds.Progress)
		if err != nil {
			return false, err
		}
		// For concurrent downloads, we expect StatusPartialContent for successful chunks
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			return false, errorutils.CheckErrorf("Failed to download %s: HTTP %d", downloadUrl, resp.StatusCode)
		}
	} else {
		// Use regular download for small files or when split is disabled
		if fileSize > 0 {
			log.Info(fmt.Sprintf("Downloading %s (%s) in single stream",
				filepath.Base(localPath),
				formatFileSize(fileSize)))
		}
		downloadFileDetails := &httpclient.DownloadFileDetails{
			DownloadPath:  downloadUrl,
			LocalPath:     filepath.Dir(localPath),
			LocalFileName: filepath.Base(localPath),
			SkipChecksum:  params.IsSkipChecksum(),
		}

		resp, err = dds.client.DownloadFile(downloadFileDetails, "", &httpClientsDetails, false, false)
		if err != nil {
			return false, err
		}

		if resp.StatusCode != http.StatusOK {
			return false, errorutils.CheckErrorf("Failed to download %s: HTTP %d", downloadUrl, resp.StatusCode)
		}
	}

	var sha256 string
	if !params.IsSkipChecksum() {
		// Get checksums from response headers
		md5FromHeader := resp.Header.Get("X-Checksum-Md5")
		sha1FromHeader := resp.Header.Get("X-Checksum-Sha1")
		sha256FromHeader := resp.Header.Get("X-Checksum-Sha256")
		sha256 = sha256FromHeader

		// Log that we're using checksums from headers
		if md5FromHeader != "" || sha1FromHeader != "" || sha256FromHeader != "" {
			log.Debug("Using checksums from response headers - MD5:", md5FromHeader != "", "SHA1:", sha1FromHeader != "", "SHA256:", sha256FromHeader != "")
		}

		// Validate using checksums from headers
		fileInfo := &utils.FileInfo{
			Checksums: struct {
				Sha1   string `json:"sha1,omitempty"`
				Sha256 string `json:"sha256,omitempty"`
				Md5    string `json:"md5,omitempty"`
			}{
				Md5:    md5FromHeader,
				Sha1:   sha1FromHeader,
				Sha256: sha256FromHeader,
			},
		}
		if err := dds.validateChecksum(fileInfo, localPath); err != nil {
			log.Warn("Checksum validation failed for", localPath, ":", err)
		}
	} else if dds.saveSummary {
		// Skip validation but still need SHA256 for summary
		sha256 = resp.Header.Get("X-Checksum-Sha256")
	}

	log.Info("Downloaded:", downloadUrl, "to", localPath)

	if dds.saveSummary && dds.filesTransfersWriter != nil {
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
		log.Debug("Writing file transfer details - Source:", sourcePath, "Target:", localPath, "RtUrl:", rtUrl, "SHA256:", sha256)
		dds.filesTransfersWriter.Write(fileTransferDetails)
	}

	return true, nil
}

// handleWildcardDownload deals with patterns like "*.zip" or "logs/2024*.txt" while downloading the artifacts.
func (dds *DirectDownloadService) handleWildcardDownload(repo, pattern string, params *DirectDownloadParams) (int, int, error) {
	dir := filepath.Dir(pattern)
	filePattern := filepath.Base(pattern)

	storagePath := fmt.Sprintf("api/storage/%s/%s", repo, dir)
	listUrl, err := clientutils.BuildUrl((*dds.artDetails).GetUrl(), storagePath, make(map[string]string))
	if err != nil {
		return 0, 0, err
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(listUrl, true, &httpClientsDetails)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				log.Warn("Failed to close response body:", err)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, errorutils.CheckErrorf("Failed to list directory %s: Status %d", listUrl, resp.StatusCode)
	}

	var storageInfo struct {
		Children []struct {
			Uri    string `json:"uri"`
			Folder bool   `json:"folder"`
		} `json:"children"`
	}

	if err := json.Unmarshal(body, &storageInfo); err != nil {
		return 0, 0, err
	}

	// Prepare a list of files to download
	var filesToDownload []string
	for _, child := range storageInfo.Children {
		if child.Folder {
			continue
		}

		fileName := strings.TrimPrefix(child.Uri, "/")
		matched, err := filepath.Match(filePattern, fileName)
		if err != nil {
			return 0, 0, err
		}

		if matched {
			// Check if the matched file is not excluded
			filePath := filepath.Join(dir, fileName)
			if !dds.isExcluded(filePath, params.GetExclusions()) {
				filesToDownload = append(filesToDownload, filePath)
			}
		}
	}

	if len(filesToDownload) == 0 {
		return 0, 0, nil
	}

	// Use goroutines for parallel downloads
	threads := dds.Threads
	if threads <= 0 {
		threads = 3
	}

	type downloadResult struct {
		success bool
		err     error
	}

	workChan := make(chan string, len(filesToDownload))
	resultChan := make(chan downloadResult, len(filesToDownload))

	for _, filePath := range filesToDownload {
		workChan <- filePath
	}
	close(workChan)

	var wg sync.WaitGroup
	for i := 0; i < threads && i < len(filesToDownload); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range workChan {
				success, err := dds.downloadSingleFile(repo, filePath, params)
				if err != nil {
					log.Error("Failed to download", filePath, ":", err)
				}
				resultChan <- downloadResult{success: success, err: err}
			}
		}()
	}

	// Wait for all downloads to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	downloadCount := 0
	failCount := 0
	for result := range resultChan {
		if result.err != nil || !result.success {
			failCount++
		} else {
			downloadCount++
		}
	}

	return downloadCount, failCount, nil
}

// getFileInfo fetches the details about a file, including its checksums
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
			if err := resp.Body.Close(); err != nil {
				log.Warn("Failed to close response body:", err)
			}
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

// formatFileSize formats bytes into human-readable format
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

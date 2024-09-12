package services

import (
	"errors"
	"fmt"
	"github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/gofrog/crypto"
	ioutils "github.com/jfrog/gofrog/io"
	"github.com/jfrog/gofrog/version"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jfrog/jfrog-client-go/http/httpclient"

	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DownloadService struct {
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
	// This map is used for validating that a downloaded release bundle is signed with a given GPG public key. This is done for security reasons.
	// The key is the release bundle name and version separated by "/" and the value is it's RbGpgValidator.
	rbGpgValidationMap map[string]*utils.RbGpgValidator
}

func NewDownloadService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *DownloadService {
	rbGpgValidationMap := make(map[string]*utils.RbGpgValidator)
	return &DownloadService{artDetails: &artDetails, client: client, rbGpgValidationMap: rbGpgValidationMap}
}

func (ds *DownloadService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ds.artDetails
}

func (ds *DownloadService) IsDryRun() bool {
	return ds.DryRun
}

func (ds *DownloadService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ds.client
}

func (ds *DownloadService) GetThreads() int {
	return ds.Threads
}

func (ds *DownloadService) SetThreads(threads int) {
	ds.Threads = threads
}

func (ds *DownloadService) SetDryRun(isDryRun bool) {
	ds.DryRun = isDryRun
}

func (ds *DownloadService) SetSaveSummary(saveSummary bool) {
	ds.saveSummary = saveSummary
}

func (ds *DownloadService) getOperationSummary(totalSucceeded, totalFailed int) *utils.OperationSummary {
	operationSummary := &utils.OperationSummary{
		TotalSucceeded: totalSucceeded,
		TotalFailed:    totalFailed,
	}
	if ds.saveSummary {
		operationSummary.TransferDetailsReader = content.NewContentReader(ds.filesTransfersWriter.GetFilePath(), content.DefaultKey)
		operationSummary.ArtifactsDetailsReader = content.NewContentReader(ds.artifactsDetailsWriter.GetFilePath(), content.DefaultKey)
	}
	return operationSummary
}

func (ds *DownloadService) DownloadFiles(downloadParams ...DownloadParams) (operationSummary *utils.OperationSummary, err error) {
	producerConsumer := parallel.NewRunner(ds.GetThreads(), 20000, false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	expectedChan := make(chan int, 1)
	successCounters := make([]int, ds.GetThreads())
	if ds.saveSummary {
		ds.filesTransfersWriter, err = content.NewContentWriter(content.DefaultKey, true, false)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = errors.Join(err, ds.filesTransfersWriter.Close())
		}()
		ds.artifactsDetailsWriter, err = content.NewContentWriter(content.DefaultKey, true, false)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = errors.Join(err, ds.artifactsDetailsWriter.Close())
		}()
	}
	ds.prepareTasks(producerConsumer, expectedChan, successCounters, errorsQueue, downloadParams...)

	err = ds.performTasks(producerConsumer, errorsQueue)
	totalSuccess := 0
	for _, v := range successCounters {
		totalSuccess += v
	}
	operationSummary = ds.getOperationSummary(totalSuccess, <-expectedChan-totalSuccess)
	return
}

func (ds *DownloadService) gpgValidateReleaseBundle(bundleParam, publicKeyFilePath string) error {
	// Check if the release bundle has already been validated.
	if ds.rbGpgValidationMap[bundleParam] != nil {
		return nil
	}
	bundleName, bundleVersion, err := utils.ParseNameAndVersion(bundleParam, false)
	if bundleName == "" || err != nil {
		return err
	}
	gpgValidator := utils.NewRbGpgValidator()
	gpgValidator.SetRbName(bundleName).SetRbVersion(bundleVersion).SetClient(ds.client).SetAtrifactoryDetails(ds.artDetails).SetPublicKey(publicKeyFilePath)
	if err = gpgValidator.Validate(); err != nil {
		return err
	}
	// Add the validated release bundle to the map.
	ds.rbGpgValidationMap[bundleParam] = gpgValidator
	return nil
}

func (ds *DownloadService) prepareTasks(producer parallel.Runner, expectedChan chan int, successCounters []int, errorsQueue *clientutils.ErrorsQueue, downloadParamsSlice ...DownloadParams) {
	go func() {
		defer producer.Done()
		defer close(expectedChan)
		totalTasks := 0
		defer func() {
			expectedChan <- totalTasks
		}()
		artifactoryVersionStr, err := ds.GetArtifactoryDetails().GetVersion()
		if err != nil {
			log.Error(err)
			errorsQueue.AddError(err)
			return
		}
		artifactoryVersion := version.NewVersion(artifactoryVersionStr)
		// Iterate over file-spec groups and produce download tasks.
		// When encountering an error, log and move to next group.
		for _, downloadParams := range downloadParamsSlice {
			utils.DisableTransitiveSearchIfNotAllowed(downloadParams.CommonParams, artifactoryVersion)
			if downloadParams.PublicGpgKey != "" {
				if err = ds.gpgValidateReleaseBundle(downloadParams.GetBundle(), downloadParams.GetPublicGpgKey()); err != nil {
					errorsQueue.AddError(err)
				}
			}

			var reader *content.ContentReader
			// Create handler function for the current group.
			fileHandlerFunc := ds.createFileHandlerFunc(downloadParams, successCounters)
			// Check if we can avoid using AQL to get the file's info.
			avoidAql, err := isFieldsProvidedToAvoidAql(downloadParams)
			// Check for search errors.
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
				continue
			}
			if avoidAql {
				reader, err = createResultsItemWithoutAql(downloadParams)
			} else {
				// Search items using AQL and get their details (size/checksum/etc.) from Artifactory.
				switch downloadParams.GetSpecType() {
				case utils.WILDCARD:
					reader, err = utils.SearchBySpecWithPattern(downloadParams.GetFile(), ds, utils.SYMLINK)
				case utils.BUILD:
					reader, err = utils.SearchBySpecWithBuild(downloadParams.GetFile(), ds)
				case utils.AQL:
					reader, err = utils.SearchBySpecWithAql(downloadParams.GetFile(), ds, utils.SYMLINK)
				}
			}
			// Check for search errors.
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
				continue
			}
			if ds.Progress != nil {
				total, _ := reader.Length()
				ds.Progress.IncGeneralProgressTotalBy(int64(total))
			}
			// Produce download tasks for the download consumers.
			totalTasks += ds.produceTasks(reader, downloadParams, producer, fileHandlerFunc, errorsQueue)
			if err = reader.Close(); err != nil {
				errorsQueue.AddError(err)
				return
			}
		}
	}()
}

func isFieldsProvidedToAvoidAql(downloadParams DownloadParams) (bool, error) {
	if downloadParams.Sha256 != "" && downloadParams.Size != nil {
		// If sha256 and size is provided, we can avoid using AQL to get the file's info.
		return true, nil
	} else if downloadParams.Sha256 == "" && downloadParams.Size == nil {
		// If sha256 and size is missing, we can't avoid using AQL to get the file's info.
		return false, nil
	}
	// If only one of the fields is provided, return an error.
	return false, errors.New("both sha256 and size must be provided in order to avoid using AQL")
}

func createResultsItemWithoutAql(downloadParams DownloadParams) (*content.ContentReader, error) {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer ioutils.Close(writer, &err)
	repo, path, name, err := breakFileDownloadPathToParts(downloadParams.GetPattern())
	if err != nil {
		return nil, err
	}
	resultItem := &utils.ResultItem{
		Type:   string(utils.File),
		Repo:   repo,
		Path:   path,
		Name:   name,
		Size:   *downloadParams.Size,
		Sha256: downloadParams.Sha256,
	}
	writer.Write(*resultItem)
	return content.NewContentReader(writer.GetFilePath(), writer.GetArrayKey()), nil
}

func breakFileDownloadPathToParts(downloadPath string) (repo, path, name string, err error) {
	if utils.IsWildcardPattern(downloadPath) {
		return "", "", "", errorutils.CheckErrorf("downloading without AQL is not supported for the provided wildcard pattern: " + downloadPath)
	}
	parts := strings.Split(downloadPath, "/")
	repo = parts[0]
	path = strings.Join(parts[1:len(parts)-1], "/")
	name = parts[len(parts)-1]
	return
}

func (ds *DownloadService) produceTasks(reader *content.ContentReader, downloadParams DownloadParams, producer parallel.Runner, fileHandler fileHandlerFunc, errorsQueue *clientutils.ErrorsQueue) int {
	flat := downloadParams.IsFlat()
	// Collect all folders path which might be needed to create.
	// key = folder path, value = the necessary data for producing create folder task.
	directoriesData := make(map[string]DownloadData)
	// Store all the paths which was created implicitly due to file upload.
	alreadyCreatedDirs := make(map[string]bool)
	// Store all the keys of directoriesData as an array.
	var directoriesDataKeys []string
	// Task counter
	var tasksCount int

	// A function that gets a ResultItem from the reader and returns a key. The reader will be sorted according to the keys returned from this function.
	// The key in our case is the local path.
	getSortKeyFunc := func(result interface{}) (string, error) {
		resultItem := new(utils.ResultItem)
		err := content.ConvertToStruct(result, &resultItem)
		if err != nil {
			return "", err
		}
		target, placeholdersUsed, err := clientutils.BuildTargetPath(downloadParams.GetPattern(), resultItem.GetItemRelativePath(), downloadParams.GetTarget(), true)
		if err != nil {
			return "", err
		}
		localPath, localFileName := fileutils.GetLocalPathAndFile(resultItem.Name, resultItem.Path, target, flat, placeholdersUsed)
		return filepath.Join(localPath, localFileName), nil
	}
	// The sort process omits results with local path that is identical to previous results.
	// We do it to avoid downloading a file and then download another file to the same path and override it.
	sortedReader, err := content.SortContentReaderByCalculatedKey(reader, getSortKeyFunc, true)
	if err != nil {
		errorsQueue.AddError(err)
		return tasksCount
	}
	defer func() {
		if err = sortedReader.Close(); err != nil {
			log.Warn("Could not close sortedReader. Error: " + err.Error())
		}
	}()
	for resultItem := new(utils.ResultItem); sortedReader.NextRecord(resultItem) == nil; resultItem = new(utils.ResultItem) {
		tempData := DownloadData{
			Dependency:   *resultItem,
			DownloadPath: downloadParams.GetPattern(),
			Target:       downloadParams.GetTarget(),
			Flat:         flat,
		}
		if resultItem.Type != string(utils.Folder) {
			if len(ds.rbGpgValidationMap) != 0 {
				// Gpg validation to the downloaded artifact
				err = rbGpgValidate(ds.rbGpgValidationMap, downloadParams.GetBundle(), resultItem)
				if err != nil {
					errorsQueue.AddError(err)
					return tasksCount
				}
			}
			// Add a task. A task is a function of type TaskFunc which later on will be executed by other go routine, the communication is done using channels.
			// The second argument is an error handling func in case the taskFunc return an error.
			tasksCount++
			_, _ = producer.AddTaskWithError(fileHandler(tempData), errorsQueue.AddError)
			// We don't want to create directories which are created explicitly by download files when CommonParams.IncludeDirs is used.
			alreadyCreatedDirs[resultItem.Path] = true
		} else {
			directoriesData, directoriesDataKeys = collectDirPathsToCreate(*resultItem, directoriesData, tempData, directoriesDataKeys)
		}
	}
	if err = sortedReader.GetError(); err != nil {
		errorsQueue.AddError(errorutils.CheckError(err))
		return tasksCount
	}
	addCreateDirsTasks(directoriesDataKeys, alreadyCreatedDirs, producer, fileHandler, directoriesData, errorsQueue, flat)
	return tasksCount
}

func rbGpgValidate(rbGpgValidationMap map[string]*utils.RbGpgValidator, bundle string, resultItem *utils.ResultItem) error {
	artifactPath := path.Join(resultItem.Repo, resultItem.Path, resultItem.Name)
	rbGpgValidator := rbGpgValidationMap[bundle]
	if rbGpgValidator == nil {
		return errorutils.CheckErrorf("release bundle validator for '%s' was not found unexpectedly. This may be caused by a bug", artifactPath)
	}
	if err := rbGpgValidator.VerifyArtifact(artifactPath, resultItem.Sha256); err != nil {
		return err
	}
	return nil
}

// Extract for the aqlResultItem the directory path, store the path the directoriesDataKeys and in the directoriesData map.
// In addition, directoriesData holds the correlate DownloadData for each key, later on this DownloadData will be used to create a create dir tasks if needed.
// This function append the new data to directoriesDataKeys and to directoriesData and return the new map and the new []string
// We are storing all the keys of directoriesData in additional array(directoriesDataKeys) so we could sort the keys and access the maps in the sorted order.
func collectDirPathsToCreate(aqlResultItem utils.ResultItem, directoriesData map[string]DownloadData, tempData DownloadData, directoriesDataKeys []string) (map[string]DownloadData, []string) {
	key := aqlResultItem.Name
	if aqlResultItem.Path != "." {
		key = path.Join(aqlResultItem.Path, aqlResultItem.Name)
	}
	directoriesData[key] = tempData
	directoriesDataKeys = append(directoriesDataKeys, key)
	return directoriesData, directoriesDataKeys
}

func addCreateDirsTasks(directoriesDataKeys []string, alreadyCreatedDirs map[string]bool, producer parallel.Runner, fileHandler fileHandlerFunc, directoriesData map[string]DownloadData, errorsQueue *clientutils.ErrorsQueue, isFlat bool) {
	// Longest path first
	// We are going to create the longest path first by doing so all sub paths of the longest path will be created implicitly.
	sort.Sort(sort.Reverse(sort.StringSlice(directoriesDataKeys)))
	for index, v := range directoriesDataKeys {
		// In order to avoid duplication we need to check the path wasn't already created by the previous action.
		if v != "." && // For some files the returned path can be the root path, ".", in that case we don't need to create any directory.
			(index == 0 || !utils.IsSubPath(directoriesDataKeys, index, "/")) { // directoriesDataKeys store all the path which might needed to be created, that's include duplicated paths.
			// By sorting the directoriesDataKeys we can assure that the longest path was created and therefore no need to create all it's sub paths.

			// Some directories were created due to file download when we aren't in flat download flow.
			if isFlat {
				_, _ = producer.AddTaskWithError(fileHandler(directoriesData[v]), errorsQueue.AddError)
			} else if !alreadyCreatedDirs[v] {
				_, _ = producer.AddTaskWithError(fileHandler(directoriesData[v]), errorsQueue.AddError)
			}
		}
	}
}

func (ds *DownloadService) performTasks(consumer parallel.Runner, errorsQueue *clientutils.ErrorsQueue) error {
	// Blocked until finish consuming
	consumer.Run()
	return errorsQueue.GetError()
}

func (ds *DownloadService) addToResults(resultItem *utils.ResultItem, rtUrl, localPath, localFileName string) {
	if ds.saveSummary {
		transferDetails := createDependencyTransferDetails(rtUrl, resultItem.GetItemRelativePath(), localPath, localFileName)
		ds.filesTransfersWriter.Write(transferDetails)
		artifactDetails := createDependencyArtifactDetails(*resultItem)
		ds.artifactsDetailsWriter.Write(artifactDetails)
	}
}

func createDependencyTransferDetails(rtUrl, relativeDownloadPath, localPath, localFileName string) clientutils.FileTransferDetails {
	fileInfo := clientutils.FileTransferDetails{
		SourcePath: relativeDownloadPath,
		RtUrl:      rtUrl,
		TargetPath: filepath.Join(localPath, localFileName),
	}
	return fileInfo
}

func createDependencyArtifactDetails(resultItem utils.ResultItem) utils.ArtifactDetails {
	fileInfo := utils.ArtifactDetails{
		ArtifactoryPath: resultItem.GetItemRelativePath(),
		Checksums: entities.Checksum{
			Sha1: resultItem.Actual_Sha1,
			Md5:  resultItem.Actual_Md5,
		},
	}
	return fileInfo
}

func createDownloadFileDetails(downloadPath, localPath, localFileName string, downloadData DownloadData, skipChecksum bool) (details *httpclient.DownloadFileDetails) {
	details = &httpclient.DownloadFileDetails{
		FileName:      downloadData.Dependency.Name,
		DownloadPath:  downloadPath,
		RelativePath:  downloadData.Dependency.GetItemRelativePath(),
		LocalPath:     localPath,
		LocalFileName: localFileName,
		Size:          downloadData.Dependency.Size,
		ExpectedSha1:  downloadData.Dependency.Actual_Sha1,
		SkipChecksum:  skipChecksum}
	return
}

func (ds *DownloadService) downloadFile(downloadFileDetails *httpclient.DownloadFileDetails, logMsgPrefix string, downloadParams DownloadParams) error {
	if ds.Progress != nil {
		ds.Progress.IncrementGeneralProgress()
	}
	httpClientsDetails := ds.GetArtifactoryDetails().CreateHttpClientDetails()
	bulkDownload := downloadParams.SplitCount == 0 || downloadParams.MinSplitSize < 0 || downloadParams.MinSplitSize*1000 > downloadFileDetails.Size
	if !bulkDownload {
		acceptRange, err := ds.isFileAcceptRange(downloadFileDetails)
		if err != nil {
			return err
		}
		bulkDownload = !acceptRange
	}
	if bulkDownload {
		var resp *http.Response
		resp, err := ds.client.DownloadFileWithProgress(downloadFileDetails, logMsgPrefix, &httpClientsDetails,
			downloadParams.IsExplode(), downloadParams.IsBypassArchiveInspection(), ds.Progress)
		if err != nil {
			return err
		}

		log.Debug(logMsgPrefix+"Artifactory response:", resp.Status)
		return errorutils.CheckResponseStatus(resp, http.StatusOK)
	}

	concurrentDownloadFlags := httpclient.ConcurrentDownloadFlags{
		FileName:                downloadFileDetails.FileName,
		DownloadPath:            downloadFileDetails.DownloadPath,
		RelativePath:            downloadFileDetails.RelativePath,
		LocalFileName:           downloadFileDetails.LocalFileName,
		LocalPath:               downloadFileDetails.LocalPath,
		ExpectedSha1:            downloadFileDetails.ExpectedSha1,
		ExpectedSha256:          downloadFileDetails.ExpectedSha256,
		FileSize:                downloadFileDetails.Size,
		SplitCount:              downloadParams.SplitCount,
		Explode:                 downloadParams.Explode,
		BypassArchiveInspection: downloadParams.BypassArchiveInspection,
		SkipChecksum:            downloadParams.SkipChecksum}

	resp, err := ds.client.DownloadFileConcurrently(concurrentDownloadFlags, logMsgPrefix, &httpClientsDetails, ds.Progress)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatus(resp, http.StatusPartialContent)
}

func (ds *DownloadService) isFileAcceptRange(downloadFileDetails *httpclient.DownloadFileDetails) (bool, error) {
	httpClientsDetails := ds.GetArtifactoryDetails().CreateHttpClientDetails()
	isAcceptRange, resp, err := ds.client.IsAcceptRanges(downloadFileDetails.DownloadPath, &httpClientsDetails)
	if err != nil {
		return false, err
	}
	err = errorutils.CheckResponseStatus(resp, http.StatusOK)
	if err != nil {
		return false, err
	}
	return isAcceptRange, err
}

func removeIfSymlink(localSymlinkPath string) error {
	if fileutils.IsPathSymlink(localSymlinkPath) {
		if err := os.Remove(localSymlinkPath); errorutils.CheckError(err) != nil {
			return err
		}
	}
	return nil
}

func createLocalSymlink(localPath, localFileName, symlinkArtifact string, symlinkChecksum bool, symlinkContentChecksum string, logMsgPrefix string) (err error) {
	if symlinkChecksum && symlinkContentChecksum != "" {
		if !fileutils.IsPathExists(symlinkArtifact, false) {
			return errorutils.CheckErrorf("symlink validation failed, target doesn't exist: " + symlinkArtifact)
		}
		var checksums map[crypto.Algorithm]string
		if checksums, err = crypto.GetFileChecksums(symlinkArtifact, crypto.SHA1); err != nil {
			return errorutils.CheckError(err)
		}
		if checksums[crypto.SHA1] != symlinkContentChecksum {
			return errorutils.CheckErrorf("symlink validation failed for target: " + symlinkArtifact)
		}
	}
	localSymlinkPath := filepath.Join(localPath, localFileName)
	isFileExists, err := fileutils.IsFileExists(localSymlinkPath, false)
	if err != nil {
		return err
	}
	// We can't create symlink in case a file with the same name already exist, we must remove the file before creating the symlink
	if isFileExists {
		if err = os.Remove(localSymlinkPath); err != nil {
			return err
		}
	}
	// Need to prepare the directories hierarchy
	_, err = fileutils.CreateFilePath(localPath, localFileName)
	if err != nil {
		return err
	}
	err = os.Symlink(symlinkArtifact, localSymlinkPath)
	if errorutils.CheckError(err) != nil {
		return err
	}
	log.Debug(fmt.Sprintf("%sCreated symlink: %q -> %q", logMsgPrefix, localSymlinkPath, symlinkArtifact))
	return nil
}

func getArtifactPropertyByKey(properties []utils.Property, key string) string {
	for _, v := range properties {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func getArtifactSymlinkPath(properties []utils.Property) string {
	return getArtifactPropertyByKey(properties, utils.ArtifactorySymlink)
}

func getArtifactSymlinkChecksum(properties []utils.Property) string {
	return getArtifactPropertyByKey(properties, utils.SymlinkSha1)
}

type fileHandlerFunc func(DownloadData) parallel.TaskFunc

func (ds *DownloadService) createFileHandlerFunc(downloadParams DownloadParams, successCounters []int) fileHandlerFunc {
	return func(downloadData DownloadData) parallel.TaskFunc {
		return func(threadId int) error {
			logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, ds.DryRun)
			downloadPath, err := clientutils.BuildUrl(ds.GetArtifactoryDetails().GetUrl(), downloadData.Dependency.GetItemRelativePath(), make(map[string]string))
			if err != nil {
				return err
			}
			if ds.DryRun {
				successCounters[threadId]++
				return nil
			}
			target, placeholdersUsed, err := clientutils.BuildTargetPath(downloadData.DownloadPath, downloadData.Dependency.GetItemRelativePath(), downloadData.Target, true)
			if err != nil {
				return err
			}
			localPath, localFileName := fileutils.GetLocalPathAndFile(downloadData.Dependency.Name, downloadData.Dependency.Path, target, downloadData.Flat, placeholdersUsed)
			localFullPath := filepath.Join(localPath, localFileName)
			if downloadData.Dependency.Type == string(utils.Folder) {
				return createDir(localFullPath, logMsgPrefix)
			}
			if err = removeIfSymlink(localFullPath); err != nil {
				return err
			}
			if downloadParams.IsSymlink() {
				if isSymlink, err := ds.createSymlinkIfNeeded(ds.GetArtifactoryDetails().GetUrl(), localPath, localFileName, logMsgPrefix, downloadData, successCounters, threadId, downloadParams); isSymlink {
					return err
				}
			}
			log.Info(fmt.Sprintf("%sDownloading %q to %q", logMsgPrefix, downloadData.Dependency.GetItemRelativePath(), localFullPath))
			if err = ds.downloadFileIfNeeded(downloadPath, localPath, localFileName, logMsgPrefix, downloadData, downloadParams); err != nil {
				log.Error(logMsgPrefix + "Received an error: " + err.Error())
				return err
			}
			successCounters[threadId]++
			ds.addToResults(&downloadData.Dependency, ds.GetArtifactoryDetails().GetUrl(), localPath, localFileName)
			return nil
		}
	}
}

func (ds *DownloadService) downloadFileIfNeeded(downloadPath, localPath, localFileName, logMsgPrefix string, downloadData DownloadData, downloadParams DownloadParams) error {
	localFilePath := filepath.Join(localPath, localFileName)
	isEqual, err := fileutils.IsEqualToLocalFile(localFilePath, downloadData.Dependency.Actual_Md5, downloadData.Dependency.Actual_Sha1)
	if err != nil {
		return err
	}
	if isEqual {
		log.Debug(logMsgPrefix+"File already exists locally:", localFilePath)
		if ds.Progress != nil {
			ds.Progress.IncrementGeneralProgress()
		}
		if downloadParams.IsExplode() {
			err = clientutils.ExtractArchive(localPath, localFileName, downloadData.Dependency.Name, logMsgPrefix, downloadParams.IsBypassArchiveInspection())
		}
		return err
	}
	downloadFileDetails := createDownloadFileDetails(downloadPath, localPath, localFileName, downloadData, downloadParams.IsSkipChecksum())
	return ds.downloadFile(downloadFileDetails, logMsgPrefix, downloadParams)
}

func createDir(folderPath, logMsgPrefix string) error {
	if err := fileutils.CreateDirIfNotExist(folderPath); err != nil {
		return err
	}
	log.Info(logMsgPrefix + "Creating folder: " + folderPath)
	return nil
}

func (ds *DownloadService) createSymlinkIfNeeded(rtUrl, localPath, localFileName, logMsgPrefix string, downloadData DownloadData, successCounters []int, threadId int, downloadParams DownloadParams) (bool, error) {
	symlinkArtifact := getArtifactSymlinkPath(downloadData.Dependency.Properties)
	isSymlink := len(symlinkArtifact) > 0
	if isSymlink {
		symlinkChecksum := getArtifactSymlinkChecksum(downloadData.Dependency.Properties)
		if err := createLocalSymlink(localPath, localFileName, symlinkArtifact, downloadParams.ValidateSymlinks(), symlinkChecksum, logMsgPrefix); err != nil {
			return isSymlink, err
		}
		successCounters[threadId]++
		ds.addToResults(&downloadData.Dependency, rtUrl, localPath, localFileName)
		return isSymlink, nil
	}
	return isSymlink, nil
}

type DownloadData struct {
	Dependency   utils.ResultItem
	DownloadPath string
	Target       string
	Flat         bool
}

type DownloadParams struct {
	*utils.CommonParams
	Symlink                 bool
	ValidateSymlink         bool
	Flat                    bool
	Explode                 bool
	BypassArchiveInspection bool
	// Min split size in Kilobytes
	MinSplitSize int64
	SplitCount   int
	PublicGpgKey string
	SkipChecksum bool

	// Optional fields (Sha256,Size) to avoid AQL request:
	Sha256 string
	// Size in bytes
	Size *int64
}

func (ds *DownloadParams) IsFlat() bool {
	return ds.Flat
}

func (ds *DownloadParams) IsExplode() bool {
	return ds.Explode
}

func (ds *DownloadParams) IsBypassArchiveInspection() bool {
	return ds.BypassArchiveInspection
}

func (ds *DownloadParams) GetFile() *utils.CommonParams {
	return ds.CommonParams
}

func (ds *DownloadParams) IsSymlink() bool {
	return ds.Symlink
}

func (ds *DownloadParams) IsSkipChecksum() bool {
	return ds.SkipChecksum
}

func (ds *DownloadParams) ValidateSymlinks() bool {
	return ds.ValidateSymlink
}

func (ds *DownloadParams) GetPublicGpgKey() string {
	return ds.PublicGpgKey
}

func NewDownloadParams() DownloadParams {
	return DownloadParams{CommonParams: &utils.CommonParams{}, MinSplitSize: 5120, SplitCount: 3}
}

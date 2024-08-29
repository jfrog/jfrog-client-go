package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/jfrog/gofrog/crypto"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/artifactory/services/fspatterns"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	// 10 KiB
	DefaultMinChecksumDeploy = utils.SizeKib * 10
	// The default minimum file size for attempting multi-part upload
	defaultUploadMinSplit = utils.SizeMiB * 200
	// The default maximum number of parts that can be concurrently uploaded per file during a multipart upload
	defaultUploadSplitCount = 5
	// Minimal file size to show progress bar (to avoid polluting the terminal with a lot of progress lines)
	minFileSizeForProgressInKb = 500 * utils.SizeKib
)

type UploadService struct {
	client          *jfroghttpclient.JfrogHttpClient
	Progress        ioutils.ProgressMgr
	ArtDetails      auth.ServiceDetails
	MultipartUpload *utils.MultipartUpload
	DryRun          bool
	Threads         int
	saveSummary     bool
	resultsManager  *resultsManager
}

const JfrogCliUploadEmptyArchiveEnv = "JFROG_CLI_UPLOAD_EMPTY_ARCHIVE"

func NewUploadService(client *jfroghttpclient.JfrogHttpClient) *UploadService {
	return &UploadService{client: client}
}

func (us *UploadService) SetThreads(threads int) {
	us.Threads = threads
}

func (us *UploadService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return us.client
}

func (us *UploadService) SetServiceDetails(artDetails auth.ServiceDetails) {
	us.ArtDetails = artDetails
}

func (us *UploadService) SetDryRun(isDryRun bool) {
	us.DryRun = isDryRun
}

func (us *UploadService) SetSaveSummary(saveSummary bool) {
	us.saveSummary = saveSummary
}

func (us *UploadService) getOperationSummary(totalSucceeded, totalFailed int) *utils.OperationSummary {
	if !us.saveSummary {
		return &utils.OperationSummary{
			TotalSucceeded: totalSucceeded,
			TotalFailed:    totalFailed,
		}
	}
	return us.resultsManager.getOperationSummary(totalSucceeded, totalFailed)
}

func (us *UploadService) UploadFiles(uploadParams ...UploadParams) (summary *utils.OperationSummary, err error) {
	// Uploading threads are using this struct to report upload results.
	uploadSummary := utils.NewResult(us.Threads)
	producerConsumer := parallel.NewRunner(us.Threads, 20000, false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	if us.saveSummary {
		us.resultsManager, err = newResultManager()
		if err != nil || us.resultsManager == nil {
			return nil, err
		}
		defer func() {
			err = errors.Join(err, us.resultsManager.close())
		}()
	}
	us.prepareUploadTasks(producerConsumer, errorsQueue, uploadSummary, uploadParams...)
	totalUploaded, totalFailed := us.performUploadTasks(producerConsumer, uploadSummary)
	return us.getOperationSummary(totalUploaded, totalFailed), errorsQueue.GetError()
}

type ArchiveUploadData struct {
	writer       *content.ContentWriter
	uploadParams UploadParams
}

func (aud *ArchiveUploadData) GetWriter() *content.ContentWriter {
	return aud.writer
}

func (aud *ArchiveUploadData) SetWriter(writer *content.ContentWriter) *ArchiveUploadData {
	aud.writer = writer
	return aud
}

func (aud *ArchiveUploadData) SetUploadParams(uploadParams UploadParams) *ArchiveUploadData {
	aud.uploadParams = uploadParams
	return aud
}

func (us *UploadService) prepareUploadTasks(producer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, uploadSummary *utils.Result, uploadParamsSlice ...UploadParams) {
	go func() {
		defer producer.Done()
		// Iterate over file-spec groups and produce upload tasks.
		// When encountering an error, log and move to next group.
		vcsCache := clientutils.NewVcsDetails()
		toArchive := make(map[string]*ArchiveUploadData)
		for _, uploadParams := range uploadParamsSlice {
			var taskHandler UploadDataHandlerFunc

			if uploadParams.Archive == "zip" {
				taskHandler = getSaveTaskInContentWriterFunc(toArchive, uploadParams, errorsQueue)
			} else {
				artifactHandlerFunc := us.createArtifactHandlerFunc(uploadSummary, uploadParams)
				taskHandler = getAddTaskToProducerFunc(producer, errorsQueue, artifactHandlerFunc)
			}

			err := CollectFilesForUpload(uploadParams, us.Progress, vcsCache, taskHandler)
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
			}
		}

		for targetPath, archiveData := range toArchive {
			err := archiveData.writer.Close()
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
			}
			if us.Progress != nil {
				us.Progress.IncGeneralProgressTotalBy(1)
			}
			_, _ = producer.AddTaskWithError(us.CreateUploadAsZipFunc(uploadSummary, targetPath, archiveData, errorsQueue), errorsQueue.AddError)
		}
	}()
}

func (us *UploadService) performUploadTasks(consumer parallel.Runner, uploadSummary *utils.Result) (totalUploaded, totalFailed int) {
	// Blocking until consuming is finished.
	consumer.Run()
	totalUploaded = utils.SumIntArray(uploadSummary.SuccessCount)
	totalUploadAttempted := utils.SumIntArray(uploadSummary.TotalCount)

	log.Debug("Uploaded", strconv.Itoa(totalUploaded), "artifacts.")
	totalFailed = totalUploadAttempted - totalUploaded
	if totalFailed > 0 {
		log.Error("Failed uploading", strconv.Itoa(totalFailed), "artifacts.")
	}
	return
}

// Creates a new Properties' struct with the artifact's props and the symlink props.
func createProperties(artifact clientutils.Artifact, uploadParams UploadParams) (properties *utils.Properties, err error) {
	artifactProps := utils.NewProperties()
	artifactSymlink := artifact.SymlinkTargetPath
	if uploadParams.IsSymlink() && len(artifactSymlink) > 0 {
		fileInfo, err := os.Stat(artifact.LocalPath)
		if err != nil {
			// If error occurred, but not due to nonexistence of Symlink target -> return empty
			if !os.IsNotExist(err) {
				return nil, errorutils.CheckError(err)
			}
			// If Symlink target exists -> get SHA1 if isn't a directory
		} else if !fileInfo.IsDir() {
			var checksums map[crypto.Algorithm]string
			checksums, err := crypto.GetFileChecksums(artifact.LocalPath, crypto.SHA1)
			if err != nil {
				return nil, errorutils.CheckError(err)
			}
			artifactProps.AddProperty(utils.SymlinkSha1, checksums[crypto.SHA1])
		}
		artifactProps.AddProperty(utils.ArtifactorySymlink, artifactSymlink)
	}
	return utils.MergeProperties([]*utils.Properties{uploadParams.GetTargetProps(), artifactProps}), nil
}

type UploadDataHandlerFunc func(data UploadData)

func getAddTaskToProducerFunc(producer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, artifactHandlerFunc artifactContext) UploadDataHandlerFunc {
	return func(data UploadData) {
		taskFunc := artifactHandlerFunc(data)
		_, _ = producer.AddTaskWithError(taskFunc, errorsQueue.AddError)
	}
}

func getSaveTaskInContentWriterFunc(writersMap map[string]*ArchiveUploadData, uploadParams UploadParams, errorsQueue *clientutils.ErrorsQueue) UploadDataHandlerFunc {
	return func(data UploadData) {
		if _, ok := writersMap[data.Artifact.TargetPath]; !ok {
			var err error
			archiveData := ArchiveUploadData{uploadParams: DeepCopyUploadParams(&uploadParams)}
			archiveData.writer, err = content.NewContentWriter("archive", true, false)
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
				return
			}
			writersMap[data.Artifact.TargetPath] = &archiveData
		} else {
			// Merge all the props
			writersMap[data.Artifact.TargetPath].uploadParams.TargetProps = utils.MergeProperties([]*utils.Properties{writersMap[data.Artifact.TargetPath].uploadParams.TargetProps, uploadParams.TargetProps})
		}
		writersMap[data.Artifact.TargetPath].writer.Write(data)
	}
}

func CollectFilesForUpload(uploadParams UploadParams, progressMgr ioutils.ProgressMgr, vcsCache *clientutils.VcsCache, dataHandlerFunc UploadDataHandlerFunc) error {
	// Target Specifies the target path in Artifactory in the following format: <repository name>/<repository path>, so it cannot start with a slash.
	// Remove leading slash if exists
	uploadParams.SetTarget(strings.TrimPrefix(uploadParams.GetTarget(), "/"))

	// Target Specifies the target path in Artifactory in the following format: <repository name>/<repository path>, so it cannot start with a slash.
	// If the received target path has no slashes then we assume that it's '<repository name>/' and we add the missing slash.
	if !strings.Contains(uploadParams.GetTarget(), "/") {
		uploadParams.SetTarget(uploadParams.GetTarget() + "/")
	}
	if uploadParams.Archive != "" && strings.HasSuffix(uploadParams.GetTarget(), "/") {
		return errorutils.CheckErrorf("an archive's target cannot be a directory")
	}
	uploadParams.SetPattern(clientutils.ReplaceTildeWithUserHome(uploadParams.GetPattern()))
	// Save parentheses index in pattern, witch have corresponding placeholder.
	rootPath, err := fspatterns.GetRootPath(uploadParams.GetPattern(), uploadParams.GetTarget(), uploadParams.TargetPathInArchive, uploadParams.GetPatternType(), uploadParams.IsSymlink())
	if err != nil {
		return err
	}

	isDir, err := fileutils.IsDirExists(rootPath, uploadParams.IsSymlink())
	if err != nil {
		return err
	}

	// If the path is a single file (or a symlink while preserving symlinks) upload it and return
	if !isDir || (fileutils.IsPathSymlink(rootPath) && uploadParams.IsSymlink()) {
		artifact, err := fspatterns.GetSingleFileToUpload(rootPath, uploadParams.GetTarget(), uploadParams.IsFlat())
		if err != nil {
			return err
		}
		props, err := createProperties(artifact, uploadParams)
		if err != nil {
			return err
		}
		buildProps := uploadParams.BuildProps
		if uploadParams.IsAddVcsProps() {
			vcsProps, err := getVcsProps(artifact.LocalPath, vcsCache)
			if err != nil {
				return err
			}
			buildProps += vcsProps
		}
		uploadData := UploadData{Artifact: artifact, TargetProps: props, BuildProps: buildProps}
		incGeneralProgressTotal(progressMgr, uploadParams)
		dataHandlerFunc(uploadData)
		return nil
	}
	if uploadParams.Ant {
		convertAntPatternToRegexp(&uploadParams)
	} else {
		convertPatternToRegexp(&uploadParams)
	}
	return scanFilesByPattern(uploadParams, rootPath, progressMgr, vcsCache, dataHandlerFunc)
}

// convertAntPatternToRegexp converts a given Ant pattern to a regular expression.
// To convert Ant patterns to regexps, we manually add parenthesis and other special characters to the pattern.
// Thus, we need to escape parentheses before converting.
func convertAntPatternToRegexp(uploadParams *UploadParams) {
	uploadParams.SetPattern(addEscapingParenthesesForUpload(uploadParams.GetPattern(), uploadParams.GetTarget(), uploadParams.TargetPathInArchive))
	uploadParams.SetPattern(clientutils.ConvertLocalPatternToRegexp(uploadParams.GetPattern(), uploadParams.GetPatternType()))
}

// convertPatternToRegexp converts a given pattern to a regular expression.
// When converting we have 2 options:
// 1. 'regexp' is true - clients are responsible for escaping parentheses that represent literal characters in the pattern - no additional treatment is required.
// 2. 'regexp' is false - it is necessary to manually escape parentheses that represent literal characters (and not placeholders).
func convertPatternToRegexp(uploadParams *UploadParams) {
	uploadParams.SetPattern(clientutils.ConvertLocalPatternToRegexp(uploadParams.GetPattern(), uploadParams.GetPatternType()))
	if !uploadParams.Regexp {
		uploadParams.SetPattern(addEscapingParenthesesForUpload(uploadParams.GetPattern(), uploadParams.GetTarget(), uploadParams.TargetPathInArchive))
	}
}

// addEscapingParenthesesForUpload escapes parentheses with no corresponding placeholder.
func addEscapingParenthesesForUpload(pattern, target, targetPathInArchive string) string {
	return clientutils.AddEscapingParentheses(pattern, target, targetPathInArchive)
}

func scanFilesByPattern(uploadParams UploadParams, rootPath string, progressMgr ioutils.ProgressMgr, vcsCache *clientutils.VcsCache, dataHandlerFunc UploadDataHandlerFunc) error {
	excludePathPattern := fspatterns.PrepareExcludePathPattern(uploadParams.Exclusions, uploadParams.GetPatternType(), uploadParams.IsRecursive())
	patternRegex, err := clientutils.GetRegExp(uploadParams.GetPattern())
	if err != nil {
		return err
	}

	paths, err := fspatterns.ListFilesFilterPatternAndSize(rootPath, uploadParams.IsRecursive(), uploadParams.IsIncludeDirs(), false, uploadParams.IsSymlink(), excludePathPattern, uploadParams.GetSizeLimit())
	if err != nil {
		return err
	}
	// Longest files path first
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))

	// 'uploadedDirs' is in use only when we need to upload folders with flat=true.
	// 'uploadedDirs' will contain only local directories paths that have been uploaded to Artifactory.
	var uploadedTargets, uploadedDirs []string

	if shouldUploadAnEmptyArchive(uploadParams.Archive, paths) {
		log.Info("All files were filtered out by the exclusion pattern, but the archive flag is set together with JFROG_CLI_UPLOAD_EMPTY_ARCHIVE environment variable. " +
			"Proceeding with an empty archive.")
		paths = []string{""}
	}

	for _, path := range paths {
		matches, isDir, err := fspatterns.SearchPatterns(path, uploadParams.IsSymlink(), uploadParams.IsIncludeDirs(), patternRegex)
		if err != nil {
			return err
		}
		if len(matches) > 0 {
			taskData, err := newUploadTaskData(path, isDir, matches, uploadParams, vcsCache, uploadParams.Regexp)
			if err != nil {
				return err
			}
			if isDir {
				if skipDirUpload(uploadedTargets, uploadedDirs, taskData.target, path, uploadParams.IsIncludeDirs()) {
					continue
				}
				uploadedDirs = append(uploadedDirs, path)
			}
			// Update progress
			incGeneralProgressTotal(progressMgr, uploadParams)
			// Create upload task
			err = createUploadTask(taskData, dataHandlerFunc, uploadParams.Regexp)
			if err != nil {
				return err
			}
			uploadedTargets = append(uploadedTargets, taskData.target)
		}
	}
	return nil
}

func shouldUploadAnEmptyArchive(archive string, paths []string) bool {
	return len(paths) == 0 &&
		archive != "" &&
		strings.ToLower(os.Getenv(JfrogCliUploadEmptyArchiveEnv)) == "true"
}

// targetFiles - Paths in Artifactory of the files that were uploaded.
// sourceDirs - Paths of the local dirs that have already been uploaded to Artifactory. (Longest files path first).
// targetDir - The directory target path to be uploaded.
// sourceDir - The directory source path to be uploaded.
func skipDirUpload(targetFiles, sourceDirs []string, targetDir, sourceDir string, includeDirs bool) bool {
	// Check that the dir is not already created in Artifactory following an upload of a child file.
	if utils.HasPrefix(targetFiles, targetDir+"/") {
		return true
	}
	// Check that the source dir is a bottom-chain dir as includeDirs expect it to be.
	if includeDirs && len(sourceDirs) > 0 && utils.HasPrefix(sourceDirs, sourceDir+fileutils.GetFileSeparator()) {
		return true
	}
	return false
}

func incGeneralProgressTotal(progressMgr ioutils.ProgressMgr, uploadParams UploadParams) {
	if progressMgr != nil {
		if uploadParams.Archive != "" {
			progressMgr.IncGeneralProgressTotalBy(2)
		} else {
			progressMgr.IncGeneralProgressTotalBy(1)
		}
	}
}

type uploadTaskData struct {
	target        string
	path          string
	isDir         bool
	isSymlinkFlow bool
	groups        []string
	size          int
	uploadParams  UploadParams
	vcsCache      *clientutils.VcsCache
}

func newUploadTaskData(path string, isDir bool, groups []string, uploadParams UploadParams, vcsCache *clientutils.VcsCache, isRegexp bool) (*uploadTaskData, error) {
	target, placeholdersUsed, err := clientutils.ReplacePlaceHolders(groups, uploadParams.GetTarget(), isRegexp)
	if err != nil {
		return nil, err
	}
	// Get symlink target (returns empty string if regular file) - Used in upload name / symlinks properties
	symlinkPath, err := fspatterns.GetFileSymlinkPath(path)
	if err != nil {
		return nil, err
	}
	// If preserving symlinks or symlink target is empty, use root path name for upload (symlink itself / regular file)
	if uploadParams.IsSymlink() || symlinkPath == "" {
		target = getUploadTarget(path, target, uploadParams.IsFlat(), placeholdersUsed)
	} else {
		target = getUploadTarget(symlinkPath, target, uploadParams.IsFlat(), placeholdersUsed)
	}
	return &uploadTaskData{target: target, path: path, isDir: isDir,
		groups: groups, size: len(groups), uploadParams: uploadParams,
		vcsCache: vcsCache,
	}, nil
}

func createUploadTask(taskData *uploadTaskData, dataHandlerFunc UploadDataHandlerFunc, isRegexp bool) error {
	symlinkPath, err := fspatterns.GetFileSymlinkPath(taskData.path)
	if err != nil {
		return err
	}
	// When using the 'archive' option for upload, we can control the target path inside the uploaded archive using placeholders.
	// This operation replace the placeholders with the relevant value.
	targetPathInArchive, _, err := clientutils.ReplacePlaceHolders(taskData.groups, taskData.uploadParams.TargetPathInArchive, isRegexp)
	if err != nil {
		return err
	}
	artifact := clientutils.Artifact{LocalPath: taskData.path, TargetPath: taskData.target, SymlinkTargetPath: symlinkPath, TargetPathInArchive: targetPathInArchive}
	props, err := createProperties(artifact, taskData.uploadParams)
	if err != nil {
		return err
	}
	buildProps := taskData.uploadParams.BuildProps
	if taskData.uploadParams.IsAddVcsProps() {
		vcsProps, err := getVcsProps(taskData.path, taskData.vcsCache)
		if err != nil {
			return err
		}
		buildProps += vcsProps
	}
	uploadData := UploadData{Artifact: artifact, TargetProps: props, BuildProps: buildProps, IsDir: taskData.isDir, IsSymlinkFlow: taskData.isSymlinkFlow}
	dataHandlerFunc(uploadData)
	return nil
}

// Construct the target path while taking `flat` flag into account.
func getUploadTarget(rootPath, target string, isFlat, placeholdersUsed bool) string {
	if strings.HasSuffix(target, "/") {
		// When placeholders are used, the file path shouldn't be taken into account (or in other words, flat = true).
		if isFlat || placeholdersUsed {
			fileName, _ := fileutils.GetFileAndDirFromPath(rootPath)
			target += fileName
		} else {
			target += clientutils.TrimPath(rootPath)
		}
	}
	return target
}

// Uploads the file in the specified local path to the specified target path.
// Returns true if the file was successfully uploaded.
func (us *UploadService) uploadFile(artifact UploadData, uploadParams UploadParams, logMsgPrefix string) (*fileutils.FileDetails, bool, error) {
	var checksumDeployed = false
	var resp *http.Response
	var details *fileutils.FileDetails
	var body []byte
	targetPathWithProps, err := buildUploadUrls(us.ArtDetails.GetUrl(), artifact.Artifact.TargetPath, artifact.BuildProps, uploadParams.GetDebian(), artifact.TargetProps)
	if err != nil {
		return nil, false, err
	}
	fileInfo, err := os.Lstat(artifact.Artifact.LocalPath)
	if errorutils.CheckError(err) != nil {
		return nil, false, err
	}
	httpClientsDetails := us.ArtDetails.CreateHttpClientDetails()
	if uploadParams.IsSymlink() && fileutils.IsFileSymlink(fileInfo) {
		resp, details, body, err = us.uploadSymlink(targetPathWithProps, logMsgPrefix, httpClientsDetails, uploadParams)
	} else {
		resp, details, body, checksumDeployed, err = us.doUpload(artifact, targetPathWithProps, logMsgPrefix, httpClientsDetails, fileInfo, uploadParams)
	}
	if err != nil {
		return nil, false, err
	}
	details.Checksum.Sha256, err = clientutils.ExtractSha256FromResponseBody(body)
	if err != nil {
		return nil, false, err
	}
	logUploadResponse(logMsgPrefix, resp, body, checksumDeployed, us.DryRun)
	return details, us.DryRun || checksumDeployed || isSuccessfulUploadStatusCode(resp.StatusCode), nil
}

func (us *UploadService) shouldTryChecksumDeploy(fileSize int64, uploadParams UploadParams) bool {
	return uploadParams.ChecksumsCalcEnabled && fileSize >= uploadParams.MinChecksumDeploy && !uploadParams.IsExplodeArchive()
}

func (us *UploadService) shouldDoMultipartUpload(fileSize int64, uploadParams UploadParams) (bool, error) {
	if uploadParams.SplitCount == 0 || fileSize < uploadParams.MinSplitSize {
		return false, nil
	}
	if fileSize > utils.MaxMultipartUploadFileSize {
		log.Debug(fmt.Sprintf("Max file size for multipart upload exceeded: %d>%d", fileSize, utils.MaxMultipartUploadFileSize))
		return false, nil
	}
	if uploadParams.IsExplodeArchive() {
		// Explode archives is not supported in multipart uploads
		return false, nil
	}
	return us.MultipartUpload.IsSupported(us.ArtDetails)
}

// Reads a file from a Reader that is given from a function (getReaderFunc) and uploads it to the specified target path.
// getReaderFunc is called only if checksum deploy was successful.
// Returns true if the file was successfully uploaded.
func (us *UploadService) uploadFileFromReader(getReaderFunc func() (io.Reader, error), targetUrlWithProps string, uploadParams UploadParams, logMsgPrefix string, details *fileutils.FileDetails) (bool, error) {
	var resp *http.Response
	var body []byte
	var checksumDeployed = false
	var err error
	httpClientsDetails := us.ArtDetails.CreateHttpClientDetails()
	if !us.DryRun {
		if us.shouldTryChecksumDeploy(details.Size, uploadParams) {
			resp, body, err = us.doChecksumDeploy(details, targetUrlWithProps, httpClientsDetails, us.client)
			if err != nil {
				return false, err
			}
			checksumDeployed = isSuccessfulUploadStatusCode(resp.StatusCode)
		}

		if !checksumDeployed {
			retryExecutor := clientutils.RetryExecutor{
				MaxRetries:               us.client.GetHttpClient().GetRetries(),
				RetriesIntervalMilliSecs: us.client.GetHttpClient().GetRetryWaitTime(),
				ErrorMessage:             fmt.Sprintf("Failure occurred while uploading to %s", targetUrlWithProps),
				LogMsgPrefix:             logMsgPrefix,
				ExecutionHandler: func() (bool, error) {
					uploadZipReader, e := getReaderFunc()
					if e != nil {
						return false, e
					}
					resp, details, body, e = us.doUploadFromReader(uploadZipReader, targetUrlWithProps, httpClientsDetails, uploadParams, details)
					if e != nil {
						return true, e
					}
					// Response must not be nil
					if resp == nil {
						return false, errorutils.CheckErrorf("%sReceived empty response from file upload", logMsgPrefix)
					}
					// If response-code < 500, should not retry
					if resp.StatusCode < 500 {
						return false, nil
					}
					// Perform retry
					log.Warn(fmt.Sprintf("%sThe server response: %s\n%s", logMsgPrefix, resp.Status, clientutils.IndentJson(body)))
					return true, nil
				},
			}

			err = retryExecutor.Execute()
			if err != nil {
				return false, err
			}
		}
	}
	logUploadResponse(logMsgPrefix, resp, body, checksumDeployed, us.DryRun)
	uploaded := us.DryRun || checksumDeployed || isSuccessfulUploadStatusCode(resp.StatusCode)
	return uploaded, nil
}

func (us *UploadService) uploadSymlink(targetPath, logMsgPrefix string, httpClientsDetails httputils.HttpClientDetails, uploadParams UploadParams) (resp *http.Response, details *fileutils.FileDetails, body []byte, err error) {
	details, err = fspatterns.CreateSymlinkFileDetails()
	if err != nil {
		return
	}
	resp, body, err = utils.UploadFile("", targetPath, logMsgPrefix, &us.ArtDetails, details, httpClientsDetails, us.client, uploadParams.ChecksumsCalcEnabled, us.Progress)
	return
}

func (us *UploadService) doUpload(artifact UploadData, targetUrlWithProps, logMsgPrefix string, httpClientsDetails httputils.HttpClientDetails, fileInfo os.FileInfo, uploadParams UploadParams) (
	resp *http.Response, details *fileutils.FileDetails, body []byte, checksumDeployed bool, err error) {
	// Get local file details
	details, err = fileutils.GetFileDetails(artifact.Artifact.LocalPath, uploadParams.ChecksumsCalcEnabled)
	if err != nil {
		return
	}

	// Return if dry run
	if us.DryRun {
		return
	}

	// Try checksum deploy
	if us.shouldTryChecksumDeploy(fileInfo.Size(), uploadParams) {
		resp, body, err = us.doChecksumDeploy(details, targetUrlWithProps, httpClientsDetails, us.client)
		if err != nil {
			return resp, details, body, checksumDeployed, err
		}
		if isSuccessfulUploadStatusCode(resp.StatusCode) {
			checksumDeployed = true
			if us.Progress != nil {
				us.Progress.IncrementGeneralProgress()
			}
			return
		}
	}

	// Try multipart upload
	var shouldTryMultipart bool
	if shouldTryMultipart, err = us.shouldDoMultipartUpload(fileInfo.Size(), uploadParams); err != nil {
		return
	}
	if shouldTryMultipart {
		var checksumToken string
		if checksumToken, err = us.MultipartUpload.UploadFileConcurrently(artifact.Artifact.LocalPath, artifact.Artifact.TargetPath,
			fileInfo.Size(), details.Checksum.Sha1, us.Progress, uploadParams.SplitCount, uploadParams.ChunkSize); err != nil {
			return
		}
		// Once the file is uploaded to the storage, we finalize the multipart upload by performing a checksum deployment to save the file in Artifactory.
		utils.AddChecksumTokenHeader(httpClientsDetails.Headers, checksumToken)
		resp, body, err = us.doChecksumDeploy(details, targetUrlWithProps, httpClientsDetails, us.client)
		return
	}

	// Do regular upload
	addExplodeHeader(&httpClientsDetails, uploadParams.IsExplodeArchive())
	resp, body, err = utils.UploadFile(artifact.Artifact.LocalPath, targetUrlWithProps, logMsgPrefix, &us.ArtDetails, details,
		httpClientsDetails, us.client, uploadParams.ChecksumsCalcEnabled, us.Progress)
	return
}

func (us *UploadService) doUploadFromReader(fileReader io.Reader, targetUrlWithProps string, httpClientsDetails httputils.HttpClientDetails, uploadParams UploadParams, details *fileutils.FileDetails) (*http.Response, *fileutils.FileDetails, []byte, error) {
	var resp *http.Response
	var body []byte
	var err error
	var reader io.Reader
	addExplodeHeader(&httpClientsDetails, uploadParams.IsExplodeArchive())
	if us.Progress != nil {
		progressReader := us.Progress.NewProgressReader(details.Size, "Uploading", targetUrlWithProps)
		reader = progressReader.ActionWithProgress(fileReader)
		progressId := progressReader.GetId()
		defer us.Progress.RemoveProgress(progressId)
	} else {
		reader = fileReader
	}
	resp, body, err = utils.UploadFileFromReader(reader, targetUrlWithProps, &us.ArtDetails, details,
		httpClientsDetails, us.client)
	return resp, details, body, err
}

func logUploadResponse(logMsgPrefix string, resp *http.Response, body []byte, checksumDeployed, isDryRun bool) {
	if resp != nil && !isSuccessfulUploadStatusCode(resp.StatusCode) {
		log.Error(logMsgPrefix + "Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
		return
	}
	if !isDryRun {
		var strChecksumDeployed string
		if checksumDeployed {
			strChecksumDeployed = " (Checksum deploy)"
		} else {
			strChecksumDeployed = ""
		}
		log.Debug(logMsgPrefix+"Artifactory response:", resp.Status, strChecksumDeployed)
	}
}

func addExplodeHeader(httpClientsDetails *httputils.HttpClientDetails, isExplode bool) {
	if isExplode {
		httpClientsDetails.AddHeader("X-Explode-Archive", "true")
	}
}

func (us *UploadService) doChecksumDeploy(details *fileutils.FileDetails, targetPath string, httpClientsDetails httputils.HttpClientDetails,
	client *jfroghttpclient.JfrogHttpClient) (resp *http.Response, body []byte, err error) {
	clonedHttpClientsDetails := httpClientsDetails.Clone()
	clonedHttpClientsDetails.AddHeader("X-Checksum-Deploy", "true")
	utils.AddChecksumHeaders(clonedHttpClientsDetails.Headers, details)
	utils.AddAuthHeaders(clonedHttpClientsDetails.Headers, us.ArtDetails)

	resp, body, err = client.SendPut(targetPath, nil, clonedHttpClientsDetails)
	return
}

func getDebianProps(debianPropsStr string) string {
	if debianPropsStr == "" {
		return ""
	}
	result := ""
	debProps := clientutils.SplitWithEscape(debianPropsStr, '/')
	for k, v := range []string{"deb.distribution", "deb.component", "deb.architecture"} {
		debProp := strings.Join([]string{v, debProps[k]}, "=")
		result = strings.Join([]string{result, debProp}, ";")
	}
	return result
}

type UploadParams struct {
	*utils.CommonParams
	Deb                  string
	BuildProps           string
	Symlink              bool
	ExplodeArchive       bool
	Flat                 bool
	AddVcsProps          bool
	MinChecksumDeploy    int64
	MinSplitSize         int64
	SplitCount           int
	ChunkSize            int64
	ChecksumsCalcEnabled bool
	Archive              string
	// When using the 'archive' option for upload, we can control the target path inside the uploaded archive using placeholders. This operation determines the TargetPathInArchive value.
	TargetPathInArchive string
	// Size limit for files to be uploaded.
	SizeLimit *fspatterns.SizeThreshold
}

func NewUploadParams() UploadParams {
	return UploadParams{CommonParams: &utils.CommonParams{}, MinChecksumDeploy: DefaultMinChecksumDeploy,
		ChecksumsCalcEnabled: true, MinSplitSize: defaultUploadMinSplit, SplitCount: defaultUploadSplitCount, ChunkSize: utils.DefaultUploadChunkSize}
}

func DeepCopyUploadParams(params *UploadParams) UploadParams {
	paramsCopy := *params
	paramsCopy.CommonParams = new(utils.CommonParams)
	*paramsCopy.CommonParams = *params.CommonParams
	return paramsCopy
}

func (up *UploadParams) IsFlat() bool {
	return up.Flat
}

func (up *UploadParams) IsSymlink() bool {
	return up.Symlink
}

func (up *UploadParams) IsAddVcsProps() bool {
	return up.AddVcsProps
}

func (up *UploadParams) IsExplodeArchive() bool {
	return up.ExplodeArchive
}

func (up *UploadParams) GetDebian() string {
	return up.Deb
}

func (up *UploadParams) GetSizeLimit() *fspatterns.SizeThreshold {
	return up.SizeLimit
}

type UploadData struct {
	Artifact      clientutils.Artifact
	TargetProps   *utils.Properties
	BuildProps    string
	IsDir         bool
	IsSymlinkFlow bool
}

type artifactContext func(UploadData) parallel.TaskFunc

func (us *UploadService) createArtifactHandlerFunc(uploadResult *utils.Result, uploadParams UploadParams) artifactContext {
	return func(artifact UploadData) parallel.TaskFunc {
		return func(threadId int) (err error) {
			uploadResult.TotalCount[threadId]++
			checksums := &entities.Checksum{}
			var uploaded bool
			logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, us.DryRun)
			log.Info(logMsgPrefix+"Uploading:", artifact.Artifact.LocalPath)
			if artifact.IsDir {
				// Upload directory
				err = us.createFolderInArtifactory(artifact)
				if err != nil {
					return
				}
				if us.Progress != nil {
					us.Progress.IncrementGeneralProgress()
				}
				uploaded = true
			} else {
				// Upload file
				var uploadFileDetails *fileutils.FileDetails
				uploadFileDetails, uploaded, err = us.uploadFile(artifact, uploadParams, logMsgPrefix)
				if err != nil {
					return
				}
				checksums = &uploadFileDetails.Checksum
			}
			if uploaded {
				us.postUpload(uploadResult, threadId, artifact, checksums)
			}
			return
		}
	}
}

func (us *UploadService) postUpload(uploadResult *utils.Result, threadId int, artifact UploadData, checksums *entities.Checksum) {
	uploadResult.SuccessCount[threadId]++
	if us.saveSummary {
		us.resultsManager.addFinalResult(artifact.Artifact.LocalPath, artifact.Artifact.TargetPath, us.ArtDetails.GetUrl(), checksums)
	}
}

func (us *UploadService) createFolderInArtifactory(artifact UploadData) error {
	url, err := clientutils.BuildUrl(us.ArtDetails.GetUrl(), artifact.Artifact.TargetPath, make(map[string]string))
	url = clientutils.AddTrailingSlashIfNeeded(url)
	if err != nil {
		return err
	}
	emptyContent := make([]byte, 0)
	httpClientsDetails := us.ArtDetails.CreateHttpClientDetails()
	resp, body, err := us.client.SendPut(url, emptyContent, &httpClientsDetails)
	if err != nil {
		log.Error(resp)
		return err
	}
	logUploadResponse("Uploaded directory:", resp, body, false, us.DryRun)
	return err
}

func (us *UploadService) CreateUploadAsZipFunc(uploadResult *utils.Result, targetPath string, archiveData *ArchiveUploadData, errorsQueue *clientutils.ErrorsQueue) parallel.TaskFunc {
	return func(threadId int) (err error) {
		uploadResult.TotalCount[threadId]++
		logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, us.DryRun)

		archiveDataReader := content.NewContentReader(archiveData.writer.GetFilePath(), archiveData.writer.GetArrayKey())
		defer func() {
			err = errors.Join(err, errorutils.CheckError(archiveDataReader.Close()))
		}()
		targetUrlWithProps, err := buildUploadUrls(us.ArtDetails.GetUrl(), targetPath, archiveData.uploadParams.BuildProps, archiveData.uploadParams.GetDebian(), archiveData.uploadParams.TargetProps)
		if err != nil {
			return
		}
		var saveFilesPathsFunc func(sourcePath string) error
		if us.saveSummary {
			saveFilesPathsFunc = func(localPath string) error {
				return us.resultsManager.addNonFinalResult(localPath, targetPath, us.ArtDetails.GetUrl())
			}
		}
		// Make sure all go routines in readFilesAsZip calls were done.
		var zipReadersWg sync.WaitGroup
		checksumZipReader := us.readFilesAsZip(archiveDataReader, "Calculating size / checksums",
			archiveData.uploadParams.Flat, archiveData.uploadParams.Symlink, saveFilesPathsFunc, errorsQueue, &zipReadersWg)
		details, err := fileutils.GetFileDetailsFromReader(checksumZipReader, archiveData.uploadParams.ChecksumsCalcEnabled)
		if err != nil {
			return
		}
		log.Info(logMsgPrefix+"Uploading artifact:", targetPath)

		getReaderFunc := func() (io.Reader, error) {
			archiveDataReader.Reset()
			return us.readFilesAsZip(archiveDataReader, "Archiving", archiveData.uploadParams.Flat,
				archiveData.uploadParams.Symlink, nil, errorsQueue, &zipReadersWg), nil
		}
		uploaded, err := us.uploadFileFromReader(getReaderFunc, targetUrlWithProps, archiveData.uploadParams, logMsgPrefix, details)

		if uploaded {
			uploadResult.SuccessCount[threadId]++
			if us.saveSummary {
				err = us.resultsManager.finalizeResult(targetPath, &details.Checksum)
			}
		}
		zipReadersWg.Wait()
		return
	}
}

// Reads files and streams them as a ZIP to a Reader.
// archiveDataReader is a ContentReader of UploadData items containing the details of the files to stream.
// saveFilesPathsFunc (optional) is a func that is called for each file that is written into the ZIP, and gets the file's local path as a parameter.
func (us *UploadService) readFilesAsZip(archiveDataReader *content.ContentReader, progressPrefix string, flat, symlink bool,
	saveFilesPathsFunc func(sourcePath string) error, errorsQueue *clientutils.ErrorsQueue, zipReadersWg *sync.WaitGroup) io.Reader {
	pr, pw := io.Pipe()
	zipReadersWg.Add(1)

	go func() {
		defer zipReadersWg.Done()
		var e error
		zipWriter := zip.NewWriter(pw)
		defer func() {
			e = zipWriter.Close()
			if e != nil {
				errorsQueue.AddError(e)
			}
			e = pw.Close()
			if e != nil {
				errorsQueue.AddError(e)
			}
		}()
		for uploadData := new(UploadData); archiveDataReader.NextRecord(uploadData) == nil; uploadData = new(UploadData) {
			if uploadData.Artifact.LocalPath != "" {
				e = us.addFileToZip(&uploadData.Artifact, progressPrefix, flat, symlink, zipWriter)
				if e != nil {
					errorsQueue.AddError(e)
				}
			}
			if saveFilesPathsFunc != nil {
				e = saveFilesPathsFunc(uploadData.Artifact.LocalPath)
				if e != nil {
					errorsQueue.AddError(e)
				}
			}
		}
		if e = archiveDataReader.GetError(); e != nil {
			errorsQueue.AddError(e)
		}
	}()

	return pr
}

func (us *UploadService) addFileToZip(artifact *clientutils.Artifact, progressPrefix string, flat, symlink bool, zipWriter *zip.Writer) (err error) {
	var reader io.Reader
	localPath := artifact.LocalPath
	// In case of a symlink there are 2 options:
	// 1. symlink == true : symlink will be added to zip as a symlink file.
	// 2. symlink == false : the symlink's target will be added to zip.
	if artifact.SymlinkTargetPath != "" && !symlink {
		localPath = artifact.SymlinkTargetPath
	}
	info, err := os.Lstat(localPath)
	if errorutils.CheckError(err) != nil {
		return
	}
	header, err := zip.FileInfoHeader(info)
	if errorutils.CheckError(err) != nil {
		return
	}
	if !flat {
		header.Name = clientutils.TrimPath(localPath)
	}
	if artifact.TargetPathInArchive != "" {
		header.Name = artifact.TargetPathInArchive
	}
	header.Method = zip.Deflate
	header.Modified = info.ModTime()

	// If this is a directory, add it to the writer with a trailing slash.
	if info.IsDir() {
		header.Name += "/"
		_, err = zipWriter.CreateHeader(header)
		return
	}
	writer, err := zipWriter.CreateHeader(header)
	if errorutils.CheckError(err) != nil {
		return
	}
	// Symlink will be written to zip as a symlink and not the symlink target file.
	if artifact.SymlinkTargetPath != "" && symlink {
		// Write symlink's target to writer - file's body for symlinks is the symlink target.
		_, err = writer.Write([]byte(filepath.ToSlash(artifact.SymlinkTargetPath)))
		return
	}
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func() {
		if file != nil {
			err = errors.Join(err, errorutils.CheckError(file.Close()))
		}
	}()
	// Show progress bar only for files larger than 250Kb to avoid polluting the terminal with endless progress bars.
	if us.Progress != nil && info.Size() > minFileSizeForProgressInKb {
		progressReader := us.Progress.NewProgressReader(info.Size(), progressPrefix, localPath)
		reader = progressReader.ActionWithProgress(file)
		progressId := progressReader.GetId()
		defer us.Progress.RemoveProgress(progressId)
	} else {
		reader = file
	}

	_, err = io.Copy(writer, reader)
	if errorutils.CheckError(err) != nil {
		return
	}
	return
}

func buildUploadUrls(artifactoryUrl, targetPath, buildProps, debianConfig string, targetProps *utils.Properties) (targetUrlWithProps string, err error) {
	targetUrl, err := clientutils.BuildUrl(artifactoryUrl, targetPath, make(map[string]string))
	if err != nil {
		return
	}
	targetUrlWithProps, err = addPropsToTargetPath(targetUrl, buildProps, debianConfig, targetProps)
	return
}

func addPropsToTargetPath(targetPath, buildProps, debConfig string, props *utils.Properties) (string, error) {
	pathParts := []string{targetPath}

	encodedTargetProps := props.ToEncodedString(false)
	if len(encodedTargetProps) > 0 {
		pathParts = append(pathParts, encodedTargetProps)
	}

	debianProps, err := utils.ParseProperties(getDebianProps(debConfig))
	if err != nil {
		return "", err
	}
	encodedDebProps := debianProps.ToEncodedString(false)
	if len(encodedDebProps) > 0 {
		pathParts = append(pathParts, encodedDebProps)
	}

	buildProperties, err := utils.ParseProperties(buildProps)
	if err != nil {
		return "", err
	}
	encodedBuildProps := buildProperties.ToEncodedString(true)
	if len(encodedBuildProps) > 0 {
		pathParts = append(pathParts, encodedBuildProps)
	}

	return strings.Join(pathParts, ";"), nil
}

func getVcsProps(path string, vcsCache *clientutils.VcsCache) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	props := ""
	revision, url, branch, err := vcsCache.GetVcsDetails(filepath.Dir(path))
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	if revision != "" {
		props += ";vcs.revision=" + revision
	}
	if url != "" {
		props += ";vcs.url=" + url
	}
	if branch != "" {
		props += ";vcs.branch=" + branch
	}
	return props, nil
}

type resultsManager struct {
	// A ContentWriter of FileTransferDetails structs. Each struct written to this ContentWriter represents a successful file transfer.
	singleFinalTransfersWriter *content.ContentWriter
	// A map for saving file transfers that are not successful (these transfers might still be in progress).
	// If the transfer is completed successfully, then the ContentWriter is closed and its file path is saved in finalTransfersFilesPaths.
	// The key of each record is the target URL of the artifact.
	// The value is a ContentWriter of FileTransferDetails structs, that all have the same target.
	notFinalTransfersWriters map[string]*content.ContentWriter
	// A slice of paths to files containing FileTransferDetails structs that represent successful file transfers.
	// These paths are of files created by ContentWriters that were in notFinalTransfersWriters.
	finalTransfersFilesPaths []string
	// A ContentWriter of ArtifactDetails structs. Each struct written to this ContentWriter represents an artifact in Artifactory
	// that was successfully uploaded in the current operation.
	artifactsDetailsWriter *content.ContentWriter
}

func newResultManager() (*resultsManager, error) {
	singleFinalTransfersWriter, e := content.NewContentWriter(content.DefaultKey, true, false)
	if e != nil {
		return nil, e
	}
	artifactsDetailsWriter, e := content.NewContentWriter(content.DefaultKey, true, false)
	if e != nil {
		return nil, e
	}
	return &resultsManager{
		singleFinalTransfersWriter: singleFinalTransfersWriter,
		notFinalTransfersWriters:   make(map[string]*content.ContentWriter),
		artifactsDetailsWriter:     artifactsDetailsWriter,
	}, nil
}

// Write a result of a successful upload.
// localPath - Path in the local file system
// targetUrl - Path in artifactory (repo-name/my/path/to/artifact)
// rtUrl - Artifactory URL (https://127.0.0.1/artifactory)
func (rm *resultsManager) addFinalResult(localPath, targetPath, rtUrl string, checksums *entities.Checksum) {
	fileTransferDetails := clientutils.FileTransferDetails{
		SourcePath: localPath,
		TargetPath: targetPath,
		RtUrl:      rtUrl,
		Sha256:     checksums.Sha256,
	}
	rm.singleFinalTransfersWriter.Write(fileTransferDetails)
	artifactDetails := utils.ArtifactDetails{
		ArtifactoryPath: targetPath,
		Checksums: entities.Checksum{
			Sha256: checksums.Sha256,
			Sha1:   checksums.Sha1,
			Md5:    checksums.Md5,
		},
	}
	rm.artifactsDetailsWriter.Write(artifactDetails)
}

// Write the details of a file transfer that is not completed yet.
// localPath - Path in the local file system
// targetUrl - Path in artifactory (repo-name/my/path/to/artifact)
// rtUrl - Artifactory URL (https://127.0.0.1/artifactory)
func (rm *resultsManager) addNonFinalResult(localPath, targetUrl, rtUrl string) error {
	if _, ok := rm.notFinalTransfersWriters[targetUrl]; !ok {
		var e error
		rm.notFinalTransfersWriters[targetUrl], e = content.NewContentWriter(content.DefaultKey, true, false)
		if e != nil {
			return e
		}
	}
	fileTransferDetails := clientutils.FileTransferDetails{
		SourcePath: localPath,
		TargetPath: targetUrl,
		RtUrl:      rtUrl,
	}
	rm.notFinalTransfersWriters[targetUrl].Write(fileTransferDetails)
	return nil
}

// Mark all the transfers to a specific target as completed successfully
func (rm *resultsManager) finalizeResult(targetPath string, checksums *entities.Checksum) error {
	writer := rm.notFinalTransfersWriters[targetPath]
	e := writer.Close()
	if e != nil {
		return e
	}
	rm.finalTransfersFilesPaths = append(rm.finalTransfersFilesPaths, writer.GetFilePath())
	delete(rm.notFinalTransfersWriters, targetPath)
	artifactDetails := utils.ArtifactDetails{
		ArtifactoryPath: targetPath,
		Checksums: entities.Checksum{
			Sha256: checksums.Sha256,
			Sha1:   checksums.Sha1,
			Md5:    checksums.Md5,
		},
	}
	rm.artifactsDetailsWriter.Write(artifactDetails)
	return nil
}

// Closes the ContentWriters that were opened by the resultManager
func (rm *resultsManager) close() error {
	err := rm.singleFinalTransfersWriter.Close()
	if err != nil {
		return err
	}
	err = rm.artifactsDetailsWriter.Close()
	if err != nil {
		return err
	}
	for _, writer := range rm.notFinalTransfersWriters {
		err = writer.Close()
		if err != nil {
			return err
		}
		err = writer.RemoveOutputFilePath()
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates an OperationSummary struct with the results. New results should not be written after this method is called.
func (rm *resultsManager) getOperationSummary(totalSucceeded, totalFailed int) *utils.OperationSummary {
	return &utils.OperationSummary{
		TransferDetailsReader:  rm.getTransferDetailsReader(),
		ArtifactsDetailsReader: content.NewContentReader(rm.artifactsDetailsWriter.GetFilePath(), content.DefaultKey),
		TotalSucceeded:         totalSucceeded,
		TotalFailed:            totalFailed,
	}
}

// Creates a ContentReader of FileTransferDetails structs. New results should not be written after this method is called.
func (rm *resultsManager) getTransferDetailsReader() *content.ContentReader {
	writersPaths := rm.finalTransfersFilesPaths
	if !rm.singleFinalTransfersWriter.IsEmpty() {
		writersPaths = append(writersPaths, rm.singleFinalTransfersWriter.GetFilePath())
	}
	return content.NewMultiSourceContentReader(writersPaths, content.DefaultKey)
}

func isSuccessfulUploadStatusCode(statusCode int) bool {
	return statusCode == http.StatusOK || statusCode == http.StatusCreated || statusCode == http.StatusAccepted
}

package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

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
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils/checksum"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type UploadService struct {
	client         *jfroghttpclient.JfrogHttpClient
	Progress       ioutils.ProgressMgr
	ArtDetails     auth.ServiceDetails
	DryRun         bool
	Threads        int
	saveSummary    bool
	resultsManager *resultsManager
}

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

func (us *UploadService) getCommandSummary(totalSucceeded, totalFailed int) *utils.CommandSummary {
	if !us.saveSummary {
		return &utils.CommandSummary{
			TotalSucceeded: totalSucceeded,
			TotalFailed:    totalFailed,
		}
	}
	return us.resultsManager.getCommandSummary(totalSucceeded, totalFailed)
}

func (us *UploadService) UploadFiles(uploadParams ...UploadParams) (*utils.CommandSummary, error) {
	// Uploading threads are using this struct to report upload results.
	var e error
	uploadSummary := utils.NewResult(us.Threads)
	producerConsumer := parallel.NewRunner(us.Threads, 20000, false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	if us.saveSummary {
		us.resultsManager, e = newResultManager()
		if e != nil {
			return nil, e
		}
		defer us.resultsManager.close()
	}
	us.prepareUploadTasks(producerConsumer, errorsQueue, uploadSummary, uploadParams...)
	totalUploaded, totalFailed := us.performUploadTasks(producerConsumer, uploadSummary)
	e = errorsQueue.GetError()
	if e != nil {
		return nil, e
	}
	return us.getCommandSummary(totalUploaded, totalFailed), nil
}

type archiveUploadData struct {
	writer       *content.ContentWriter
	uploadParams UploadParams
}

func (us *UploadService) prepareUploadTasks(producer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, uploadSummary *utils.Result, uploadParamsSlice ...UploadParams) {
	go func() {
		defer producer.Done()
		// Iterate over file-spec groups and produce upload tasks.
		// When encountering an error, log and move to next group.
		vcsCache := clientutils.NewVcsDetals()
		toArchive := make(map[string]*archiveUploadData)
		for _, uploadParams := range uploadParamsSlice {
			var taskHandler uploadDataHandlerFunc

			if uploadParams.Archive == "zip" {
				taskHandler = getSaveTaskInContentWriterFunc(toArchive, uploadParams, errorsQueue)
			} else {
				artifactHandlerFunc := us.createArtifactHandlerFunc(uploadSummary, uploadParams)
				taskHandler = getAddTaskToProducerFunc(producer, errorsQueue, artifactHandlerFunc)
			}

			err := collectFilesForUpload(uploadParams, us.Progress, vcsCache, taskHandler)
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
			}
		}

		for targetPath, archiveData := range toArchive {
			archiveData.writer.Close()
			if us.Progress != nil {
				us.Progress.IncGeneralProgressTotalBy(1)
			}
			producer.AddTaskWithError(us.createUploadAsZipFunc(uploadSummary, targetPath, archiveData, errorsQueue), errorsQueue.AddError)
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

func addSymlinkProps(artifact clientutils.Artifact, uploadParams UploadParams) error {
	artifactProps := ""
	artifactSymlink := artifact.Symlink
	if uploadParams.IsSymlink() && len(artifactSymlink) > 0 {
		sha1Property := ""
		fileInfo, err := os.Stat(artifact.LocalPath)
		if err != nil {
			// If error occurred, but not due to nonexistence of Symlink target -> return empty
			if !os.IsNotExist(err) {
				return errorutils.CheckError(err)
			}
			// If Symlink target exists -> get SHA1 if isn't a directory
		} else if !fileInfo.IsDir() {
			file, err := os.Open(artifact.LocalPath)
			if err != nil {
				return errorutils.CheckError(err)
			}
			defer file.Close()
			checksumInfo, err := checksum.Calc(file, checksum.SHA1)
			if err != nil {
				return err
			}
			sha1 := checksumInfo[checksum.SHA1]
			sha1Property = ";" + utils.SYMLINK_SHA1 + "=" + sha1
		}
		artifactProps += utils.ARTIFACTORY_SYMLINK + "=" + artifactSymlink + sha1Property
	}
	uploadParams.TargetProps = clientutils.AddProps(uploadParams.GetTargetProps(), artifactProps)
	return nil
}

type uploadDataHandlerFunc func(data UploadData)

func getAddTaskToProducerFunc(producer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, artifactHandlerFunc artifactContext) uploadDataHandlerFunc {
	return func(data UploadData) {
		taskFunc := artifactHandlerFunc(data)
		producer.AddTaskWithError(taskFunc, errorsQueue.AddError)
	}
}

func getSaveTaskInContentWriterFunc(writersMap map[string]*archiveUploadData, uploadParams UploadParams, errorsQueue *clientutils.ErrorsQueue) uploadDataHandlerFunc {
	return func(data UploadData) {
		if _, ok := writersMap[data.Artifact.TargetPath]; !ok {
			var err error
			archiveData := archiveUploadData{uploadParams: uploadParams}
			archiveData.writer, err = content.NewContentWriter("archive", true, false)
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
				return
			}
			writersMap[data.Artifact.TargetPath] = &archiveData
		}

		// Merge all of the props
		writersMap[data.Artifact.TargetPath].uploadParams.TargetProps = strings.Join([]string{writersMap[data.Artifact.TargetPath].uploadParams.TargetProps, uploadParams.TargetProps}, ";")
		writersMap[data.Artifact.TargetPath].writer.Write(data)
	}
}

func collectFilesForUpload(uploadParams UploadParams, progressMgr ioutils.ProgressMgr, vcsCache *clientutils.VcsCache, dataHandlerFunc uploadDataHandlerFunc) error {
	if strings.Index(uploadParams.GetTarget(), "/") < 0 {
		uploadParams.SetTarget(uploadParams.GetTarget() + "/")
	}
	if uploadParams.Archive != "" && strings.HasSuffix(uploadParams.GetTarget(), "/") {
		return errors.New("An archive's target cannot be a directory.")
	}
	uploadParams.SetPattern(clientutils.ReplaceTildeWithUserHome(uploadParams.GetPattern()))
	// Save parentheses index in pattern, witch have corresponding placeholder.
	rootPath, err := fspatterns.GetRootPath(uploadParams.GetPattern(), uploadParams.GetTarget(), uploadParams.IsRegexp(), uploadParams.IsSymlink())
	if err != nil {
		return err
	}

	isDir, err := fileutils.IsDirExists(rootPath, uploadParams.IsSymlink())
	if err != nil {
		return err
	}

	// If the path is a single file (or a symlink while preserving symlinks) upload it and return
	if !isDir || (fileutils.IsPathSymlink(rootPath) && uploadParams.IsSymlink()) {
		artifact, err := fspatterns.GetSingleFileToUpload(rootPath, uploadParams.GetTarget(), uploadParams.IsFlat(), uploadParams.IsSymlink())
		if err != nil {
			return err
		}
		if err = addSymlinkProps(artifact, uploadParams); err != nil {
			return err
		}
		if uploadParams.IsAddVcsProps() {
			vcsProps, err := getVcsProps(artifact.LocalPath, vcsCache)
			if err != nil {
				return err
			}
			uploadParams.BuildProps += vcsProps
		}
		uploadData := UploadData{Artifact: artifact, TargetProps: uploadParams.GetTargetProps(), BuildProps: uploadParams.BuildProps}
		if progressMgr != nil {
			if uploadParams.Archive != "" {
				progressMgr.IncGeneralProgressTotalBy(2)
			} else {
				progressMgr.IncGeneralProgressTotalBy(1)
			}
		}
		dataHandlerFunc(uploadData)
		return err
	}
	uploadParams.SetPattern(clientutils.PrepareLocalPathForUpload(uploadParams.GetPattern(), uploadParams.IsRegexp()))
	err = collectPatternMatchingFiles(uploadParams, rootPath, progressMgr, vcsCache, dataHandlerFunc)
	return err
}

func collectPatternMatchingFiles(uploadParams UploadParams, rootPath string, progressMgr ioutils.ProgressMgr, vcsCache *clientutils.VcsCache, dataHandlerFunc uploadDataHandlerFunc) error {
	excludePathPattern := fspatterns.PrepareExcludePathPattern(uploadParams)
	patternRegex, err := regexp.Compile(uploadParams.GetPattern())
	if errorutils.CheckError(err) != nil {
		return err
	}

	paths, err := fspatterns.GetPaths(rootPath, uploadParams.IsRecursive(), uploadParams.IsIncludeDirs(), uploadParams.IsSymlink())
	if err != nil {
		return err
	}
	// Longest paths first
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	// 'foldersPaths' is a subset of the 'paths' array. foldersPaths is in use only when we need to upload folders with flat=true.
	// 'foldersPaths' will contain only the directories paths which are in the 'paths' array.
	var foldersPaths []string
	for index, path := range paths {
		matches, isDir, isSymlinkFlow, err := fspatterns.PrepareAndFilterPaths(path, excludePathPattern, uploadParams.IsSymlink(), uploadParams.IsIncludeDirs(), patternRegex)
		if err != nil {
			return err
		}

		if matches != nil && len(matches) > 0 {
			target := uploadParams.GetTarget()
			tempPaths := paths
			tempIndex := index
			// In case we need to upload directories with flat=true, we want to avoid the creation of unnecessary paths in Artifactory.
			// To achieve this, we need to take into consideration the directories which had already been uploaded, ignoring all files paths.
			// When flat=false we take into consideration folder paths which were created implicitly by file upload
			if uploadParams.IsFlat() && uploadParams.IsIncludeDirs() && isDir {
				foldersPaths = append(foldersPaths, path)
				tempPaths = foldersPaths
				tempIndex = len(foldersPaths) - 1
			}
			taskData := &uploadTaskData{target: target, path: path, isDir: isDir, isSymlinkFlow: isSymlinkFlow,
				paths: tempPaths, groups: matches, index: tempIndex, size: len(matches), uploadParams: uploadParams,
				vcsCache: vcsCache,
			}
			if progressMgr != nil {
				if uploadParams.Archive != "" {
					progressMgr.IncGeneralProgressTotalBy(2)
				} else {
					progressMgr.IncGeneralProgressTotalBy(1)
				}
			}
			createUploadTask(taskData, dataHandlerFunc)
		}
	}
	return nil
}

type uploadTaskData struct {
	target        string
	path          string
	isDir         bool
	isSymlinkFlow bool
	paths         []string
	groups        []string
	index         int
	size          int
	uploadParams  UploadParams
	vcsCache      *clientutils.VcsCache
}

func createUploadTask(taskData *uploadTaskData, dataHandlerFunc uploadDataHandlerFunc) error {
	for i := 1; i < taskData.size; i++ {
		group := strings.Replace(taskData.groups[i], "\\", "/", -1)
		taskData.target = strings.Replace(taskData.target, "{"+strconv.Itoa(i)+"}", group, -1)
	}

	// Get symlink target (returns empty string if regular file) - Used in upload name / symlinks properties
	symlinkPath, err := fspatterns.GetFileSymlinkPath(taskData.path)
	if err != nil {
		return err
	}

	// If preserving symlinks or symlink target is empty, use root path name for upload (symlink itself / regular file)
	if taskData.uploadParams.IsSymlink() || symlinkPath == "" {
		taskData.target = getUploadTarget(taskData.path, taskData.target, taskData.uploadParams.IsFlat())
	} else {
		taskData.target = getUploadTarget(symlinkPath, taskData.target, taskData.uploadParams.IsFlat())
	}

	artifact := clientutils.Artifact{LocalPath: taskData.path, TargetPath: taskData.target, Symlink: symlinkPath}
	if err := addSymlinkProps(artifact, taskData.uploadParams); err != nil {
		return err
	}
	if taskData.uploadParams.IsAddVcsProps() {
		vcsProps, err := getVcsProps(taskData.path, taskData.vcsCache)
		if err != nil {
			return err
		}
		taskData.uploadParams.BuildProps += vcsProps
	}
	uploadData := UploadData{Artifact: artifact, TargetProps: taskData.uploadParams.GetTargetProps(), BuildProps: taskData.uploadParams.BuildProps}
	if taskData.isDir && taskData.uploadParams.IsIncludeDirs() && !taskData.isSymlinkFlow {
		if taskData.path != "." && (taskData.index == 0 || !utils.IsSubPath(taskData.paths, taskData.index, fileutils.GetFileSeparator())) {
			uploadData.IsDir = true
		} else {
			return nil
		}
	}
	dataHandlerFunc(uploadData)
	return nil
}

// Construct the target path while taking `flat` flag into account.
func getUploadTarget(rootPath, target string, isFlat bool) string {
	if strings.HasSuffix(target, "/") {
		if isFlat {
			fileName, _ := fileutils.GetFileAndDirFromPath(rootPath)
			target += fileName
		} else {
			target += clientutils.TrimPath(rootPath)
		}
	}
	return target
}

func addPropsToTargetPath(targetPath, props, buildProps, debConfig string) (string, error) {
	propsStr := strings.Join([]string{props, getDebianProps(debConfig)}, ";")
	properties, err := utils.ParseProperties(propsStr, utils.SplitCommas)
	if err != nil {
		return "", err
	}
	buildProperties, err := utils.ParseProperties(buildProps, utils.JoinCommas)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{targetPath, properties.ToEncodedString(), buildProperties.ToEncodedString()}, ";"), nil
}

func prepareUploadData(localPath, baseTargetPath, props, buildProps string, uploadParams UploadParams, logMsgPrefix string) (fileInfo os.FileInfo, targetPath string, err error) {
	targetPath, err = addPropsToTargetPath(baseTargetPath, props, buildProps, uploadParams.GetDebian())
	if errorutils.CheckError(err) != nil {
		return
	}
	log.Info(logMsgPrefix+"Uploading artifact:", localPath)

	fileInfo, err = os.Lstat(localPath)
	errorutils.CheckError(err)
	return
}

// Uploads the file in the specified local path to the specified target path.
// Returns true if the file was successfully uploaded.
func (us *UploadService) uploadFile(localPath, targetPath, targetUrl, props, buildProps string, uploadParams UploadParams, logMsgPrefix string) (*fileutils.FileDetails, bool, error) {
	fileInfo, targetPathWithProps, err := prepareUploadData(localPath, targetUrl, props, buildProps, uploadParams, logMsgPrefix)
	if err != nil {
		return nil, false, err
	}

	var checksumDeployed = false
	var resp *http.Response
	var details *fileutils.FileDetails
	var body []byte
	httpClientsDetails := us.ArtDetails.CreateHttpClientDetails()
	if errorutils.CheckError(err) != nil {
		return nil, false, err
	}
	if uploadParams.IsSymlink() && fileutils.IsFileSymlink(fileInfo) {
		resp, details, body, err = us.uploadSymlink(targetPathWithProps, logMsgPrefix, httpClientsDetails, uploadParams)
	} else {
		resp, details, body, checksumDeployed, err = us.doUpload(localPath, targetPath, targetPathWithProps, logMsgPrefix, httpClientsDetails, fileInfo, uploadParams)
	}
	if err != nil {
		return nil, false, err
	}
	logUploadResponse(logMsgPrefix, resp, body, checksumDeployed, us.DryRun)
	return details, us.DryRun || checksumDeployed || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK, nil
}

// Reads a file from a Reader that is given from a function (getReaderFunc) and uploads it to the specified target path.
// getReaderFunc is called only if checksum deploy was successful.
// Returns true if the file was successfully uploaded.
func (us *UploadService) uploadFileFromReader(getReaderFunc func() (io.Reader, error), targetPath, targetUrlWithProps string, uploadParams UploadParams, logMsgPrefix string, details *fileutils.FileDetails) (bool, error) {
	var resp *http.Response
	var body []byte
	var checksumDeployed = false
	var e error
	httpClientsDetails := us.ArtDetails.CreateHttpClientDetails()
	if !us.DryRun {
		if details.Size >= uploadParams.MinChecksumDeploy && !uploadParams.IsExplodeArchive() {
			resp, body, e = us.tryChecksumDeploy(details, targetUrlWithProps, httpClientsDetails, us.client)
			if e != nil {
				return false, e
			}
			checksumDeployed = resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK
		}

		if !checksumDeployed {
			retryExecutor := clientutils.RetryExecutor{
				MaxRetries:      uploadParams.Retries,
				RetriesInterval: 0,
				ErrorMessage:    fmt.Sprintf("Failure occurred while uploading to %s", targetUrlWithProps),
				LogMsgPrefix:    logMsgPrefix,
				ExecutionHandler: func() (bool, error) {
					uploadZipReader, e := getReaderFunc()
					if e != nil {
						return false, e
					}
					resp, details, body, e = us.doUploadFromReader(uploadZipReader, targetPath, targetUrlWithProps, httpClientsDetails, uploadParams, details)
					if e != nil {
						return true, e
					}
					// Response must not be nil
					if resp == nil {
						return false, errorutils.CheckError(errors.New(fmt.Sprintf("%sReceived empty response from file upload", logMsgPrefix)))
					}
					// If response-code < 500, should not retry
					if resp.StatusCode < 500 {
						return false, nil
					}
					// Perform retry
					log.Warn(fmt.Sprintf("%sThe server response: %s", logMsgPrefix, resp.Status))
					return true, nil
				},
			}

			e = retryExecutor.Execute()
			if e != nil {
				return false, e
			}
		}
	}
	logUploadResponse(logMsgPrefix, resp, body, checksumDeployed, us.DryRun)
	uploaded := us.DryRun || checksumDeployed || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK
	return uploaded, nil
}

func (us *UploadService) uploadSymlink(targetPath, logMsgPrefix string, httpClientsDetails httputils.HttpClientDetails, uploadParams UploadParams) (resp *http.Response, details *fileutils.FileDetails, body []byte, err error) {
	details, err = fspatterns.CreateSymlinkFileDetails()
	if err != nil {
		return
	}
	resp, body, err = utils.UploadFile("", targetPath, logMsgPrefix, &us.ArtDetails, details, httpClientsDetails, us.client, uploadParams.GetRetries(), nil, "")
	return
}

func (us *UploadService) doUpload(localPath, targetPath, targetUrlWithProps, logMsgPrefix string, httpClientsDetails httputils.HttpClientDetails, fileInfo os.FileInfo, uploadParams UploadParams) (*http.Response, *fileutils.FileDetails, []byte, bool, error) {
	var details *fileutils.FileDetails
	var checksumDeployed bool
	var resp *http.Response
	var body []byte
	var err error
	addExplodeHeader(&httpClientsDetails, uploadParams.IsExplodeArchive())
	if !us.DryRun {
		if fileInfo.Size() >= uploadParams.MinChecksumDeploy && !uploadParams.IsExplodeArchive() {
			details, err = fileutils.GetFileDetails(localPath)
			if err != nil {
				return resp, details, body, checksumDeployed, err
			}
			resp, body, err = us.tryChecksumDeploy(details, targetUrlWithProps, httpClientsDetails, us.client)
			if err != nil {
				return resp, details, body, checksumDeployed, err
			}
			checksumDeployed = resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK
		}
		if !checksumDeployed {
			resp, body, err = utils.UploadFile(localPath, targetUrlWithProps, logMsgPrefix, &us.ArtDetails, details,
				httpClientsDetails, us.client, uploadParams.Retries, us.Progress, targetPath)
			if err != nil {
				return resp, details, body, checksumDeployed, err
			}
		}
	}
	if details == nil {
		details, err = fileutils.GetFileDetails(localPath)
	}
	return resp, details, body, checksumDeployed, err
}

func (us *UploadService) doUploadFromReader(fileReader io.Reader, targetPath, targetUrlWithProps string, httpClientsDetails httputils.HttpClientDetails, uploadParams UploadParams, details *fileutils.FileDetails) (*http.Response, *fileutils.FileDetails, []byte, error) {
	var resp *http.Response
	var body []byte
	var err error
	var reader io.Reader
	addExplodeHeader(&httpClientsDetails, uploadParams.IsExplodeArchive())
	if us.Progress != nil {
		progressReader := us.Progress.NewProgressReader(details.Size, "Uploading", targetPath)
		reader = progressReader.ActionWithProgress(fileReader)
		defer us.Progress.RemoveProgress(progressReader.GetId())
	} else {
		reader = fileReader
	}
	resp, body, err = utils.UploadFileFromReader(reader, targetUrlWithProps, &us.ArtDetails, details,
		httpClientsDetails, us.client)
	return resp, details, body, err
}

func logUploadResponse(logMsgPrefix string, resp *http.Response, body []byte, checksumDeployed, isDryRun bool) {
	if resp != nil && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
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
		log.Debug(logMsgPrefix, "Artifactory response:", resp.Status, strChecksumDeployed)
	}
}

func addExplodeHeader(httpClientsDetails *httputils.HttpClientDetails, isExplode bool) {
	if isExplode {
		utils.AddHeader("X-Explode-Archive", "true", &httpClientsDetails.Headers)
	}
}

func (us *UploadService) tryChecksumDeploy(details *fileutils.FileDetails, targetPath string, httpClientsDetails httputils.HttpClientDetails,
	client *jfroghttpclient.JfrogHttpClient) (resp *http.Response, body []byte, err error) {
	requestClientDetails := httpClientsDetails.Clone()
	utils.AddHeader("X-Checksum-Deploy", "true", &requestClientDetails.Headers)
	utils.AddChecksumHeaders(requestClientDetails.Headers, details)
	utils.AddAuthHeaders(requestClientDetails.Headers, us.ArtDetails)

	resp, body, err = client.SendPut(targetPath, nil, requestClientDetails)
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
	*utils.ArtifactoryCommonParams
	Deb               string
	BuildProps        string
	Symlink           bool
	ExplodeArchive    bool
	Flat              bool
	AddVcsProps       bool
	Retries           int
	MinChecksumDeploy int64
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

func (up *UploadParams) GetRetries() int {
	return up.Retries
}

type UploadData struct {
	Artifact    clientutils.Artifact
	TargetProps string
	BuildProps  string
	IsDir       bool
}

type artifactContext func(UploadData) parallel.TaskFunc

func (us *UploadService) createArtifactHandlerFunc(uploadResult *utils.Result, uploadParams UploadParams) artifactContext {
	return func(artifact UploadData) parallel.TaskFunc {
		return func(threadId int) (e error) {
			if artifact.IsDir {
				us.createFolderInArtifactory(artifact)
				return
			}
			uploadResult.TotalCount[threadId]++
			logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, us.DryRun)
			targetUrl, e := utils.BuildArtifactoryUrl(us.ArtDetails.GetUrl(), artifact.Artifact.TargetPath, make(map[string]string))
			if e != nil {
				return
			}
			uploadFileDetails, uploaded, e := us.uploadFile(artifact.Artifact.LocalPath, artifact.Artifact.TargetPath, targetUrl, artifact.TargetProps, artifact.BuildProps, uploadParams, logMsgPrefix)
			if e != nil {
				return
			}
			if uploaded {
				uploadResult.SuccessCount[threadId]++
				if us.saveSummary {
					us.resultsManager.addFinalResult(artifact.Artifact.LocalPath, targetUrl, &uploadFileDetails.Checksum)
				}
			}
			return
		}
	}
}

func (us *UploadService) createFolderInArtifactory(artifact UploadData) error {
	url, err := utils.BuildArtifactoryUrl(us.ArtDetails.GetUrl(), artifact.Artifact.TargetPath, make(map[string]string))
	url = clientutils.AddTrailingSlashIfNeeded(url)
	if err != nil {
		return err
	}
	content := make([]byte, 0)
	httpClientsDetails := us.ArtDetails.CreateHttpClientDetails()
	resp, body, err := us.client.SendPut(url, content, &httpClientsDetails)
	if err != nil {
		log.Debug(resp)
		return err
	}
	logUploadResponse("Uploaded directory:", resp, body, false, us.DryRun)
	return err
}

func (us *UploadService) createUploadAsZipFunc(uploadResult *utils.Result, targetPath string, archiveData *archiveUploadData, errorsQueue *clientutils.ErrorsQueue) parallel.TaskFunc {
	return func(threadId int) (e error) {
		uploadResult.TotalCount[threadId]++
		logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, us.DryRun)

		archiveDataReader := content.NewContentReader(archiveData.writer.GetFilePath(), archiveData.writer.GetArrayKey())
		defer archiveDataReader.Close()
		targetUrl, e := utils.BuildArtifactoryUrl(us.ArtDetails.GetUrl(), targetPath, make(map[string]string))
		if e != nil {
			return
		}
		var saveFilesPathsFunc func(sourcePath string) error
		if us.saveSummary {
			saveFilesPathsFunc = func(localPath string) error {
				return us.resultsManager.addNotFinalResult(localPath, targetUrl)
			}
		}
		checksumZipReader := us.readFilesAsZip(archiveDataReader, "Calculating checksums", archiveData.uploadParams.Flat, saveFilesPathsFunc, errorsQueue)
		details, e := fileutils.GetFileDetailsFromReader(checksumZipReader)
		if e != nil {
			return
		}
		targetUrlWithProps, e := addPropsToTargetPath(targetUrl, archiveData.uploadParams.TargetProps,
			archiveData.uploadParams.BuildProps, archiveData.uploadParams.GetDebian())
		if e != nil {
			return
		}
		log.Info(logMsgPrefix+"Uploading artifact:", targetPath)

		getReaderFunc := func() (io.Reader, error) {
			archiveDataReader.Reset()
			return us.readFilesAsZip(archiveDataReader, "Archiving", archiveData.uploadParams.Flat, nil, errorsQueue), nil
		}
		uploaded, e := us.uploadFileFromReader(getReaderFunc, targetPath, targetUrlWithProps, archiveData.uploadParams, logMsgPrefix, details)

		if uploaded {
			uploadResult.SuccessCount[threadId]++
			if us.saveSummary {
				e = us.resultsManager.finalizeResult(targetUrl, &details.Checksum)
			}
		}
		return
	}
}

// Reads files and streams them as a ZIP to a Reader.
// archiveDataReader is a ContentReader of UploadData items containing the details of the files to stream.
// saveFilesPathsFunc (optional) is a func that is called for each file that is written into the ZIP, and gets the file's local path as a parameter.
func (us *UploadService) readFilesAsZip(archiveDataReader *content.ContentReader, progressPrefix string, flat bool,
	saveFilesPathsFunc func(sourcePath string) error, errorsQueue *clientutils.ErrorsQueue) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		zipWriter := zip.NewWriter(pw)
		defer pw.Close()
		defer zipWriter.Close()
		for uploadData := new(UploadData); archiveDataReader.NextRecord(uploadData) == nil; uploadData = new(UploadData) {
			var e error
			if uploadData.Artifact.Symlink != "" {
				e = us.addFileToZip(uploadData.Artifact.Symlink, progressPrefix, flat, zipWriter)
			} else {
				e = us.addFileToZip(uploadData.Artifact.LocalPath, progressPrefix, flat, zipWriter)
			}

			if e != nil {
				errorsQueue.AddError(e)
			}
			if saveFilesPathsFunc != nil {
				e = saveFilesPathsFunc(uploadData.Artifact.LocalPath)
				if e != nil {
					errorsQueue.AddError(e)
				}
			}
		}
	}()

	return pr
}

func (us *UploadService) addFileToZip(localPath, progressPrefix string, flat bool, zipWriter *zip.Writer) (e error) {
	var reader io.Reader
	file, e := os.Open(localPath)
	defer file.Close()
	if e != nil {
		return
	}
	info, e := file.Stat()
	if e != nil {
		return
	}
	header, e := zip.FileInfoHeader(info)
	if e != nil {
		return
	}
	if !flat {
		header.Name = clientutils.TrimPath(localPath)
	}
	header.Method = zip.Deflate
	writer, e := zipWriter.CreateHeader(header)
	if e != nil {
		return
	}

	if us.Progress != nil {
		progressReader := us.Progress.NewProgressReader(info.Size(), progressPrefix, localPath)
		reader = progressReader.ActionWithProgress(file)
		defer us.Progress.RemoveProgress(progressReader.GetId())
	} else {
		reader = file
	}

	_, e = io.Copy(writer, reader)
	if e != nil {
		return
	}
	return
}

func NewUploadParams() UploadParams {
	return UploadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{}, MinChecksumDeploy: 10240}
}

func getVcsProps(path string, vcsCache *clientutils.VcsCache) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	props := ""
	revision, url, err := vcsCache.GetVcsDetails(filepath.Dir(path))
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	if revision != "" {
		props += ";vcs.revision=" + revision
	}
	if url != "" {
		props += ";vcs.url=" + url
	}
	return props, nil
}

type resultsManager struct {
	singleFinalTransfersWriter *content.ContentWriter
	notFinalTransfersWriters   map[string]*content.ContentWriter
	finalTransfersFilesPaths   []string
	artifactsDetailsWriter     *content.ContentWriter
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

func (rm *resultsManager) addFinalResult(localPath, targetUrl string, checksums *fileutils.ChecksumDetails) {
	fileTransferDetails := utils.FileTransferDetails{
		LocalPath:       localPath,
		ArtifactoryPath: targetUrl,
	}
	rm.singleFinalTransfersWriter.Write(fileTransferDetails)
	artifactDetails := utils.ArtifactDetails{
		ArtifactoryPath: targetUrl,
		Checksums: utils.Checksums{
			Sha256: checksums.Sha256,
			Sha1:   checksums.Sha1,
			Md5:    checksums.Md5,
		},
	}
	rm.artifactsDetailsWriter.Write(artifactDetails)
}

func (rm *resultsManager) addNotFinalResult(localPath, targetUrl string) error {
	if _, ok := rm.notFinalTransfersWriters[targetUrl]; !ok {
		var e error
		rm.notFinalTransfersWriters[targetUrl], e = content.NewContentWriter(content.DefaultKey, true, false)
		if e != nil {
			return e
		}
	}
	fileTransferDetails := utils.FileTransferDetails{
		LocalPath:       localPath,
		ArtifactoryPath: targetUrl,
	}
	rm.notFinalTransfersWriters[targetUrl].Write(fileTransferDetails)
	return nil
}

func (rm *resultsManager) finalizeResult(targetPath string, checksums *fileutils.ChecksumDetails) error {
	writer := rm.notFinalTransfersWriters[targetPath]
	e := writer.Close()
	if e != nil {
		return e
	}
	rm.finalTransfersFilesPaths = append(rm.finalTransfersFilesPaths, writer.GetFilePath())
	delete(rm.notFinalTransfersWriters, targetPath)
	artifactDetails := utils.ArtifactDetails{
		ArtifactoryPath: targetPath,
		Checksums: utils.Checksums{
			Sha256: checksums.Sha256,
			Sha1:   checksums.Sha1,
			Md5:    checksums.Md5,
		},
	}
	rm.artifactsDetailsWriter.Write(artifactDetails)
	return nil
}

func (rm *resultsManager) close() {
	rm.singleFinalTransfersWriter.Close()
	rm.artifactsDetailsWriter.Close()
}

func (rm *resultsManager) getCommandSummary(totalSucceeded, totalFailed int) *utils.CommandSummary {
	return &utils.CommandSummary{
		TransferDetailsReader:  rm.getTransferDetailsReader(),
		ArtifactsDetailsReader: content.NewContentReader(rm.artifactsDetailsWriter.GetFilePath(), content.DefaultKey),
		TotalSucceeded:         totalSucceeded,
		TotalFailed:            totalFailed,
	}
}

func (rm *resultsManager) getTransferDetailsReader() *content.ContentReader {
	writersPaths := rm.finalTransfersFilesPaths
	if !rm.singleFinalTransfersWriter.IsEmpty() {
		writersPaths = append(writersPaths, rm.singleFinalTransfersWriter.GetFilePath())
	}
	return content.NewCombinedContentReader(writersPaths, content.DefaultKey)
}

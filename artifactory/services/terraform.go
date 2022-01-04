package services

import (
	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

const ArtifactoryMinSupportedVersion = "6.10.0" // change version

// Support for Artifactory 6.10.0 and above API
type TerraformPublishCommand struct {
	artifactoryVersion string
	clientDetails      httputils.HttpClientDetails
	client             *jfroghttpclient.JfrogHttpClient
}

type TerraformService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
	Threads    int
	namespace  string
	provider   string
	tag        string
	targetRepo string
}

func NewTerraformService(client *jfroghttpclient.JfrogHttpClient, artDetails auth.ServiceDetails) *TerraformService {
	return &TerraformService{client: client, ArtDetails: artDetails, Threads: 3}
}

func (ts *TerraformService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ts.client
}

func (ts *TerraformService) SetServiceDetails(artDetails auth.ServiceDetails) {
	ts.ArtDetails = artDetails
}

func (ts *TerraformService) GetServiceDetails() auth.ServiceDetails {
	return ts.ArtDetails
}

func (ts *TerraformService) GetThreads() int {
	return ts.Threads
}

func (ts *TerraformService) TerraformPublish(terraformParams *TerraformParams) (int, int, error) {
	// Uploading threads are using this struct to report upload results.
	var e error
	uploadSummary := utils.NewResult(ts.Threads)
	producerConsumer := parallel.NewRunner(ts.Threads, 20000, false)
	errorsQueue := clientutils.NewErrorsQueue(1)

	ts.prepareTerraformPublishTasks(producerConsumer, errorsQueue, uploadSummary, *terraformParams)
	totalUploaded, totalFailed := ts.performTerraformPublishTasks(producerConsumer, uploadSummary)
	e = errorsQueue.GetError()
	if e != nil {
		return 0, 0, e
	}
	return totalUploaded, totalFailed, nil
}

func (ts *TerraformService) prepareTerraformPublishTasks(producer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, uploadSummary *utils.Result, terraformParams TerraformParams) {
	go func() {
		defer producer.Done()
		// Iterate over file-spec groups and produce upload tasks.
		// When encountering an error, log and move to next group.

		//vcsCache := clientutils.NewVcsDetails()
		toArchive := make(map[string]*archiveUploadData)

		pwd, err := os.Getwd()
		if err != nil {
			log.Error(err)
			errorsQueue.AddError(err)
		}
		filepath.WalkDir(pwd, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
				return err
			}
			pathIinfo, e := os.Lstat(path)
			if e != nil {
				log.Error(e)
				errorsQueue.AddError(e)
				return e
			}
			// Skip files and check only directories.
			if !pathIinfo.IsDir() {
				return nil
			}
			terraformModule, e := isTerraformModule(path)
			if e != nil {
				log.Error(e)
				errorsQueue.AddError(e)
				return e
			}

			if terraformModule {
				moduleName := info.Name()
				target, e := getPublishTarget(moduleName, &terraformParams)
				if e != nil {
					log.Error(e)
					errorsQueue.AddError(e)
					return e
				}
				//artifact, e := fspatterns.GetSingleFileToUpload(path, target, false)
				//artifact.TargetPathInZipFile =
				if e != nil {
					log.Error(e)
					errorsQueue.AddError(e)
					return e
				}
				terraformParams.Pattern = path + string(os.PathSeparator)
				terraformParams.Target = target
				//uploadData := UploadData{Artifact: artifact, TargetProps: utils.NewProperties()}
				uploadParams := deepCopyTerraformToUploadParams(&terraformParams)
				dataHandlerFunc := getSaveTaskInContentWriterFunc(toArchive, uploadParams, errorsQueue)

				err := collectFilesForUpload(uploadParams, nil, nil, dataHandlerFunc)
				if err != nil {
					log.Error(err)
					errorsQueue.AddError(err)
				}
				//incGeneralProgressTotal(progressMgr, uploadParams)
				//dataHandlerFunc(uploadData)
				//return tpc.doDeploy(path, moduleName, tpc.serverDetails)

			}
			return nil
		})
		for targetPath, archiveData := range toArchive {
			err := archiveData.writer.Close()
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
			}
			//if us.Progress != nil {
			//	us.Progress.IncGeneralProgressTotalBy(1)
			//}
			uploadService := NewUploadService(ts.client)
			uploadService.SetServiceDetails(ts.ArtDetails)
			uploadService.SetThreads(ts.Threads)
			producer.AddTaskWithError(uploadService.createUploadAsZipFunc(uploadSummary, targetPath, archiveData, errorsQueue), errorsQueue.AddError)
		}
	}()
}

func (ts *TerraformService) performTerraformPublishTasks(consumer parallel.Runner, uploadSummary *utils.Result) (totalUploaded, totalFailed int) {
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

func deepCopyTerraformToUploadParams(params *TerraformParams) UploadParams {
	uploadParams := NewUploadParams()
	uploadParams.Archive = "zip"
	uploadParams.Recursive = true
	uploadParams.Exclusions = []string{"*.git", "*.DS_Store"}
	uploadParams.CommonParams = new(utils.CommonParams)
	uploadParams.CommonParams.TargetProps = utils.NewProperties()
	uploadParams.Target = params.Target
	uploadParams.Pattern = params.Pattern
	//uploadParams.CommonParams = params.CommonParams
	return uploadParams
}

func getPublishTarget(moduleName string, terraformParams *TerraformParams) (string, error) {
	return filepath.ToSlash(filepath.Join(terraformParams.TargetRepo, terraformParams.Namespace, terraformParams.Provider, moduleName, terraformParams.Tag+".zip")), nil
}

func getTargetPathInZipFile(moduleName string, terraformParams *TerraformParams) (string, error) {

	return filepath.ToSlash(filepath.Join(terraformParams.TargetRepo, terraformParams.Namespace, terraformParams.Provider, moduleName, terraformParams.Tag+".zip")), nil
}

// We identify a terraform module by the existing of '.tf' files inside the directory.
// isTerraformModule search for '.tf' file inside the directory and returns true it founds at least one.
func isTerraformModule(path string) (bool, error) {
	dirname := path + string(filepath.Separator)

	d, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return false, err
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".tf" {
				return true, nil
			}
		}
	}
	return false, nil
}

type TerraformParams struct {
	*utils.CommonParams
	Namespace  string
	Provider   string
	Tag        string
	TargetRepo string
}

func (tp *TerraformParams) GetNamespace() string {
	return tp.Namespace
}

func (tp *TerraformParams) GetProvider() string {
	return tp.Provider
}

func (tp *TerraformParams) GetTag() string {
	return tp.Tag
}

func (tp *TerraformParams) GetTargetRepo() string {
	return tp.TargetRepo
}

func (tp *TerraformParams) SetNamespace(namespace string) *TerraformParams {
	tp.Namespace = namespace
	return tp
}

func (tp *TerraformParams) SetProvider(provider string) *TerraformParams {
	tp.Provider = provider
	return tp
}

func (tp *TerraformParams) SetTag(tag string) *TerraformParams {
	tp.Tag = tag
	return tp
}

func (tp *TerraformParams) SetTargetRepo(repo string) *TerraformParams {
	tp.TargetRepo = repo
	return tp
}

func NewTerraformParams(commonParams *utils.CommonParams) *TerraformParams {
	return &TerraformParams{CommonParams: commonParams}
}

func (tpc *TerraformPublishCommand) verifyCompatibleVersion(artifactoryVersion string) error {
	propertiesApi := ArtifactoryMinSupportedVersion
	version := version.NewVersion(artifactoryVersion)
	tpc.artifactoryVersion = artifactoryVersion
	if !version.AtLeast(propertiesApi) {
		return errorutils.CheckErrorf("Unsupported version of Artifactory: %s\nSupports Artifactory version %s and above", artifactoryVersion, propertiesApi)
	}
	return nil
}

// Creates an OperationSummary struct with the results. New results should not be written after this method is called.
func (rm *resultsManager) getOperationSummaryTerraform(totalSucceeded, totalFailed int) *utils.OperationSummary {
	return &utils.OperationSummary{
		TransferDetailsReader:  rm.getTransferDetailsReader(),
		ArtifactsDetailsReader: content.NewContentReader(rm.artifactsDetailsWriter.GetFilePath(), content.DefaultKey),
		TotalSucceeded:         totalSucceeded,
		TotalFailed:            totalFailed,
	}
}

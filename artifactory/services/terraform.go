package services

import (
	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

type TerraformService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
	Threads    int
	namespace  string
	provider   string
	tag        string
	targetRepo string
}

func NewTerraformService(client *jfroghttpclient.JfrogHttpClient) *TerraformService {
	return &TerraformService{client: client}
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
		toArchive := make(map[string]*archiveUploadData)
		// Walk and upload directories which contain '.tf' files.
		pwd, err := os.Getwd()
		if err != nil {
			log.Error(err)
			errorsQueue.AddError(err)
		}
		err = filepath.WalkDir(pwd, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			pathInfo, e := os.Lstat(path)
			if e != nil {
				return e
			}
			// Skip files and check only directories.
			if !pathInfo.IsDir() {
				return nil
			}
			isTerraformModule, e := checkIfTerraformModule(path)
			if e != nil {
				return e
			}
			if isTerraformModule {
				uploadParams, e := terraformParams.uploadParamsForTerraformPublish(pathInfo.Name(), path)
				if e != nil {
					return e
				}
				dataHandlerFunc := getSaveTaskInContentWriterFunc(toArchive, *uploadParams, errorsQueue)
				return collectFilesForUpload(*uploadParams, nil, nil, dataHandlerFunc)
			}
			return nil
		})
		if err != nil {
			log.Error(err)
			errorsQueue.AddError(err)
		}
		// Upload modules
		for targetPath, archiveData := range toArchive {
			err := archiveData.writer.Close()
			if err != nil {
				log.Error(err)
				errorsQueue.AddError(err)
			}
			// Upload module using upload service
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

func (tp *TerraformParams) uploadParamsForTerraformPublish(moduleName, dirPath string) (*UploadParams, error) {
	uploadParams := NewUploadParams()
	target, e := tp.getPublishTarget(moduleName)
	if e != nil {
		return nil, e
	}
	uploadParams.Target = target
	uploadParams.Pattern = dirPath + "/(*)"
	uploadParams.TargetPathInArchive = "{1}"
	uploadParams.Archive = "zip"
	uploadParams.Recursive = true
	uploadParams.Exclusions = []string{"*.git", "*.DS_Store"}
	uploadParams.CommonParams.TargetProps = utils.NewProperties()

	return &uploadParams, nil
}

// Module's path in terraform repo : namespace/provider/moduleName/tag.zip
func (tp *TerraformParams) getPublishTarget(moduleName string) (string, error) {
	return filepath.ToSlash(filepath.Join(tp.TargetRepo, tp.Namespace, tp.Provider, moduleName, tp.Tag+".zip")), nil
}

// We identify a Terraform module by the existing of a '.tf' file inside the module directory.
// isTerraformModule search for '.tf' file inside and returns true it founds at least one.
func checkIfTerraformModule(path string) (bool, error) {
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

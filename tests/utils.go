package tests

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pipelinesAuth "github.com/jfrog/jfrog-client-go/pipelines/auth"
	pipelinesServices "github.com/jfrog/jfrog-client-go/pipelines/services"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"

	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distributionAuth "github.com/jfrog/jfrog-client-go/distribution/auth"
	distributionServices "github.com/jfrog/jfrog-client-go/distribution/services"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	xrayAuth "github.com/jfrog/jfrog-client-go/xray/auth"
	xrayServices "github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/assert"
)

var TestArtifactory *bool
var TestDistribution *bool
var TestXray *bool
var TestPipelines *bool
var RtUrl *string
var DistUrl *string
var XrayUrl *string
var PipelinesUrl *string
var RtUser *string
var RtPassword *string
var RtSshKeyPath *string
var RtSshPassphrase *string
var RtAccessToken *string
var PipelinesAccessToken *string
var PipelinesVcsToken *string
var PipelinesVcsRepoFullPath *string
var PipelinesVcsBranch *string

// Artifactory services
var testsUploadService *services.UploadService
var testsSearchService *services.SearchService
var testsDeleteService *services.DeleteService
var testsDownloadService *services.DownloadService
var testsSecurityService *services.SecurityService
var testsCreateLocalRepositoryService *services.LocalRepositoryService
var testsCreateRemoteRepositoryService *services.RemoteRepositoryService
var testsCreateVirtualRepositoryService *services.VirtualRepositoryService
var testsUpdateLocalRepositoryService *services.LocalRepositoryService
var testsUpdateRemoteRepositoryService *services.RemoteRepositoryService
var testsUpdateVirtualRepositoryService *services.VirtualRepositoryService
var testsDeleteRepositoryService *services.DeleteRepositoryService
var testsRepositoriesService *services.RepositoriesService
var testsCreateReplicationService *services.CreateReplicationService
var testsUpdateReplicationService *services.UpdateReplicationService
var testsReplicationGetService *services.GetReplicationService
var testsReplicationDeleteService *services.DeleteReplicationService
var testsPermissionTargetService *services.PermissionTargetService
var testUserService *services.UserService
var testGroupService *services.GroupService
var testBuildInfoService *services.BuildInfoService

// Distribution services
var testsBundleSetSigningKeyService *distributionServices.SetSigningKeyService
var testsBundleCreateService *distributionServices.CreateReleaseBundleService
var testsBundleUpdateService *distributionServices.UpdateReleaseBundleService
var testsBundleSignService *distributionServices.SignBundleService
var testsBundleDistributeService *distributionServices.DistributeReleaseBundleService
var testsBundleDistributionStatusService *distributionServices.DistributionStatusService
var testsBundleDeleteLocalService *distributionServices.DeleteLocalReleaseBundleService
var testsBundleDeleteRemoteService *distributionServices.DeleteReleaseBundleService

// Xray Services
var testsXrayWatchService *xrayServices.WatchService
var testsXrayPolicyService *xrayServices.PolicyService
var testXrayBinMgrService *xrayServices.BinMgrService

// Pipelines Services
var testsPipelinesIntegrationsService *pipelinesServices.IntegrationsService
var testsPipelinesSourcesService *pipelinesServices.SourcesService

var timestamp = time.Now().Unix()
var trueValue = true
var falseValue = false

const (
	RtTargetRepo                     = "jfrog-client-tests-repo1/"
	SpecsTestRepositoryConfig        = "specs_test_repository_config.json"
	RepoDetailsUrl                   = "api/repositories/"
	HttpClientCreationFailureMessage = "Failed while attempting to create HttpClient: %s"
	RepoKeyPrefixForRepoServiceTest  = "jf-client-go-test"
)

func init() {
	TestArtifactory = flag.Bool("test.artifactory", false, "Test Artifactory")
	TestDistribution = flag.Bool("test.distribution", false, "Test distribution")
	TestXray = flag.Bool("test.xray", false, "Test xray")
	TestPipelines = flag.Bool("test.pipelines", false, "Test pipelines")
	RtUrl = flag.String("rt.url", "", "Artifactory url")
	DistUrl = flag.String("ds.url", "", "Distribution url")
	XrayUrl = flag.String("xr.url", "", "Xray url")
	PipelinesUrl = flag.String("pipe.url", "", "Pipelines url")
	RtUser = flag.String("rt.user", "", "Artifactory username")
	RtPassword = flag.String("rt.password", "", "Artifactory password")
	RtSshKeyPath = flag.String("rt.sshKeyPath", "", "Ssh key file path")
	RtSshPassphrase = flag.String("rt.sshPassphrase", "", "Ssh key passphrase")
	RtAccessToken = flag.String("rt.accessToken", "", "Artifactory access token")
	PipelinesAccessToken = flag.String("pipe.accessToken", "", "Pipelines access token")
	PipelinesVcsToken = flag.String("pipe.vcsToken", "", "Vcs token for Pipelines tests")
	PipelinesVcsRepoFullPath = flag.String("pipe.vcsRepo", "", "Vcs full repo path for Pipelines tests")
	PipelinesVcsBranch = flag.String("pipe.vcsBranch", "", "Vcs branch for Pipelines tests")
}

func createArtifactorySecurityManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsSecurityService = services.NewSecurityService(client)
	testsSecurityService.ArtDetails = artDetails
}

func createArtifactorySearchManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsSearchService = services.NewSearchService(artDetails, client)
}

func createArtifactoryDeleteManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsDeleteService = services.NewDeleteService(artDetails, client)
	testsDeleteService.SetThreads(3)
}

func createArtifactoryUploadManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsUploadService = services.NewUploadService(client)
	testsUploadService.ArtDetails = artDetails
	testsUploadService.Threads = 3
}

func createArtifactoryUserManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testUserService = services.NewUserService(client)
	testUserService.ArtDetails = artDetails
}

func createArtifactoryGroupManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testGroupService = services.NewGroupService(client)
	testGroupService.ArtDetails = artDetails
}

func createArtifactoryBuildInfoManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testBuildInfoService = services.NewBuildInfoService(artDetails, client)
}

func createArtifactoryDownloadManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsDownloadService = services.NewDownloadService(artDetails, client)
	testsDownloadService.SetThreads(3)
}

func createDistributionManager() {
	distDetails := GetDistDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(distDetails.GetClientCertPath()).
		SetClientCertKeyPath(distDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(distDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testsBundleCreateService = distributionServices.NewCreateReleaseBundleService(client)
	testsBundleUpdateService = distributionServices.NewUpdateReleaseBundleService(client)
	testsBundleSignService = distributionServices.NewSignBundleService(client)
	testsBundleDistributeService = distributionServices.NewDistributeReleaseBundleService(client)
	testsBundleDistributionStatusService = distributionServices.NewDistributionStatusService(client)
	testsBundleDeleteLocalService = distributionServices.NewDeleteLocalDistributionService(client)
	testsBundleSetSigningKeyService = distributionServices.NewSetSigningKeyService(client)
	testsBundleDeleteRemoteService = distributionServices.NewDeleteReleaseBundleService(client)
	testsBundleCreateService.DistDetails = distDetails
	testsBundleUpdateService.DistDetails = distDetails
	testsBundleSignService.DistDetails = distDetails
	testsBundleDistributeService.DistDetails = distDetails
	testsBundleDistributionStatusService.DistDetails = distDetails
	testsBundleDeleteLocalService.DistDetails = distDetails
	testsBundleSetSigningKeyService.DistDetails = distDetails
	testsBundleDeleteRemoteService.DistDetails = distDetails
}

func createArtifactoryCreateLocalRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsCreateLocalRepositoryService = services.NewLocalRepositoryService(client, false)
	testsCreateLocalRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateLocalRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsUpdateLocalRepositoryService = services.NewLocalRepositoryService(client, true)
	testsUpdateLocalRepositoryService.ArtDetails = artDetails
}

func createArtifactoryCreateRemoteRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsCreateRemoteRepositoryService = services.NewRemoteRepositoryService(client, false)
	testsCreateRemoteRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateRemoteRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsUpdateRemoteRepositoryService = services.NewRemoteRepositoryService(client, true)
	testsUpdateRemoteRepositoryService.ArtDetails = artDetails
}

func createArtifactoryCreateVirtualRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsCreateVirtualRepositoryService = services.NewVirtualRepositoryService(client, false)
	testsCreateVirtualRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateVirtualRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsUpdateVirtualRepositoryService = services.NewVirtualRepositoryService(client, true)
	testsUpdateVirtualRepositoryService.ArtDetails = artDetails
}

func createArtifactoryDeleteRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsDeleteRepositoryService = services.NewDeleteRepositoryService(client)
	testsDeleteRepositoryService.ArtDetails = artDetails
}

func createArtifactoryGetRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsRepositoriesService = services.NewRepositoriesService(client)
	testsRepositoriesService.ArtDetails = artDetails
}

func createArtifactoryReplicationCreateManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsCreateReplicationService = services.NewCreateReplicationService(client)
	testsCreateReplicationService.ArtDetails = artDetails
}

func createArtifactoryReplicationUpdateManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsUpdateReplicationService = services.NewUpdateReplicationService(client)
	testsUpdateReplicationService.ArtDetails = artDetails
}

func createArtifactoryReplicationGetManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsReplicationGetService = services.NewGetReplicationService(client)
	testsReplicationGetService.ArtDetails = artDetails
}

func createArtifactoryReplicationDeleteManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsReplicationDeleteService = services.NewDeleteReplicationService(client)
	testsReplicationDeleteService.ArtDetails = artDetails
}

func createArtifactoryPermissionTargetManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsPermissionTargetService = services.NewPermissionTargetService(client)
	testsPermissionTargetService.ArtDetails = artDetails
}

func createJfrogHttpClient(artDetails *auth.ServiceDetails) (*jfroghttpclient.JfrogHttpClient, error) {
	return jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath((*artDetails).GetClientCertPath()).
		SetClientCertKeyPath((*artDetails).GetClientCertKeyPath()).
		AppendPreRequestInterceptor((*artDetails).RunPreRequestFunctions).
		Build()
}

func createXrayWatchManager() {
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testsXrayWatchService = xrayServices.NewWatchService(client)
	testsXrayWatchService.XrayDetails = xrayDetails
}

func createXrayPolicyManager() {
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testsXrayPolicyService = xrayServices.NewPolicyService(client)
	testsXrayPolicyService.XrayDetails = xrayDetails
}

func createXrayBinMgrManager() {
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testXrayBinMgrService = xrayServices.NewBinMgrService(client)
	testXrayBinMgrService.XrayDetails = xrayDetails
}

func createPipelinesIntegrationsManager() {
	pipelinesDetails := GetPipelinesDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(pipelinesDetails.GetClientCertPath()).
		SetClientCertKeyPath(pipelinesDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(pipelinesDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testsPipelinesIntegrationsService = pipelinesServices.NewIntegrationsService(client)
	testsPipelinesIntegrationsService.ServiceDetails = pipelinesDetails
}

func createPipelinesSourcesManager() {
	pipelinesDetails := GetPipelinesDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(pipelinesDetails.GetClientCertPath()).
		SetClientCertKeyPath(pipelinesDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(pipelinesDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testsPipelinesSourcesService = pipelinesServices.NewSourcesService(client)
	testsPipelinesSourcesService.ServiceDetails = pipelinesDetails
}

func failOnHttpClientCreation(err error) {
	if err != nil {
		log.Error(fmt.Sprintf(HttpClientCreationFailureMessage, err.Error()))
		os.Exit(1)
	}
}

func GetRtDetails() auth.ServiceDetails {
	rtDetails := artifactoryAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*RtUrl))
	setAuthenticationDetail(rtDetails)
	return rtDetails
}

func GetDistDetails() auth.ServiceDetails {
	distDetails := distributionAuth.NewDistributionDetails()
	distDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*DistUrl))
	setAuthenticationDetail(distDetails)
	return distDetails
}

func GetXrayDetails() auth.ServiceDetails {
	xrayDetails := xrayAuth.NewXrayDetails()
	xrayDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*XrayUrl))
	setAuthenticationDetail(xrayDetails)
	return xrayDetails
}

func GetPipelinesDetails() auth.ServiceDetails {
	pDetails := pipelinesAuth.NewPipelinesDetails()
	pDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*PipelinesUrl))
	pDetails.SetAccessToken(*PipelinesAccessToken)
	return pDetails
}

func setAuthenticationDetail(details auth.ServiceDetails) {
	if !fileutils.IsSshUrl(details.GetUrl()) {
		if *RtAccessToken != "" {
			details.SetAccessToken(*RtAccessToken)
		} else {
			details.SetUser(*RtUser)
			details.SetPassword(*RtPassword)
		}
		return
	}

	err := details.AuthenticateSsh(*RtSshKeyPath, *RtSshPassphrase)
	if err != nil {
		log.Error("Failed while attempting to authenticate: " + err.Error())
		os.Exit(1)
	}
}

func uploadDummyFile(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := FixWinPath(filepath.Join(workingDir, "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "test/"}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "b.in"}
	up.Flat = true
	summary, err = testsUploadService.UploadFiles(up)
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	archivePath := filepath.Join(workingDir, "c.tar.gz")
	err = archiver.Archive([]string{filepath.Join(workingDir, "out/a.in")}, archivePath)
	if err != nil {
		t.Error(err)
	}
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: archivePath, Recursive: true, Target: RtTargetRepo}
	up.Flat = true
	summary, err = testsUploadService.UploadFiles(up)
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
}

func artifactoryCleanup(t *testing.T) {
	params := &utils.ArtifactoryCommonParams{Pattern: RtTargetRepo}
	toDelete, err := testsDeleteService.GetPathsToDelete(services.DeleteParams{ArtifactoryCommonParams: params})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer toDelete.Close()
	NumberOfItemToDelete, err := toDelete.Length()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testsDeleteService.SetThreads(3)
	deletedCount, err := testsDeleteService.DeleteFiles(toDelete)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if NumberOfItemToDelete != deletedCount {
		t.Errorf("Failed to delete files from Artifactory expected %d items to be deleted got %d.", NumberOfItemToDelete, deletedCount)
	}
}

func createReposIfNeeded() error {
	if !(*TestArtifactory || *TestDistribution || *TestXray) {
		return nil
	}
	var err error
	var repoConfig string
	repo := RtTargetRepo
	if strings.HasSuffix(repo, "/") {
		repo = repo[0:strings.LastIndex(repo, "/")]
	}
	if !isRepoExist(repo) {
		repoConfig = filepath.Join(getTestDataPath(), "reposconfig", SpecsTestRepositoryConfig)
		err = execCreateRepoRest(repoConfig, repo)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func isRepoExist(repoName string) bool {
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	resp, _, _, err := client.SendGet(artDetails.GetUrl()+RepoDetailsUrl+repoName, true, artHttpDetails, "")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusBadRequest {
		return true
	}
	return false
}

func execCreateRepoRest(repoConfig, repoName string) error {
	content, err := ioutil.ReadFile(repoConfig)
	if err != nil {
		return err
	}
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()

	artHttpDetails.Headers = map[string]string{"Content-Type": "application/json"}
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, _, err := client.SendPut(artDetails.GetUrl()+"api/repositories/"+repoName, content, artHttpDetails, "")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("Fail to create repository. Reason local repository with key: " + repoName + " already exist\n")
	}
	log.Info("Repository", repoName, "created.")
	return nil
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "testdata")
}

func FixWinPath(filePath string) string {
	fixedPath := strings.Replace(filePath, "\\", "\\\\", -1)
	return fixedPath
}

func getRepoConfig(repoKey string) (body []byte, err error) {
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return
	}
	resp, body, _, err := client.SendGet(artDetails.GetUrl()+"api/repositories/"+repoKey, false, artHttpDetails, "")
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	return
}

func createRepoConfigValidationFunc(repoKey string, expectedConfig interface{}) clientutils.ExecutionHandlerFunc {
	return func() (shouldRetry bool, err error) {
		config, err := getRepoConfig(repoKey)
		if err != nil || config == nil {
			return true, errors.New("failed reading repository config for " + repoKey)
		}
		var confMap, expectedConfigMap map[string]interface{}
		if err = json.Unmarshal(config, &confMap); err != nil {
			return false, errors.New("failed unmarshalling repository config for " + repoKey)
		}
		tmpJson, err := json.Marshal(expectedConfig)
		if err != nil {
			return false, errors.New("failed marshalling expected config for " + repoKey)
		}
		if err = json.Unmarshal(tmpJson, &expectedConfigMap); err != nil {
			return false, errors.New("failed unmarshalling expected config for " + repoKey)
		}
		for key, expectedValue := range expectedConfigMap {
			if !assert.ObjectsAreEqual(confMap[key], expectedValue) {
				errMsg := fmt.Sprintf("config validation for %s failed. key: %s expected: %s actual: %s", repoKey, key, expectedValue, confMap[key])
				return true, errors.New(errMsg)

			}
		}
		return false, nil
	}
}

func validateRepoConfig(t *testing.T, repoKey string, params interface{}) {
	retryExecutor := &clientutils.RetryExecutor{
		MaxRetries:       12,
		RetriesInterval:  10,
		ErrorMessage:     "Waiting for Artifactory to evaluate repository operation...",
		ExecutionHandler: createRepoConfigValidationFunc(repoKey, params),
	}
	err := retryExecutor.Execute()
	assert.NoError(t, err)
}

func deleteRepo(t *testing.T, repoKey string) {
	err := testsDeleteRepositoryService.Delete(repoKey)
	assert.NoError(t, err, "Failed to delete "+repoKey)
}

func GenerateRepoKeyForRepoServiceTest() string {
	timestamp++
	return fmt.Sprintf("%s-%d", RepoKeyPrefixForRepoServiceTest, timestamp)
}

func getRepo(t *testing.T, repoKey string) *services.RepositoryDetails {
	data := services.RepositoryDetails{}
	err := testsRepositoriesService.Get(repoKey, &data)
	assert.NoError(t, err, "Failed to get "+repoKey+" details")
	return &data
}

func getAllRepos(t *testing.T, repoType, packageType string) *[]services.RepositoryDetails {
	params := services.NewRepositoriesFilterParams()
	params.RepoType = repoType
	params.PackageType = packageType
	data, err := testsRepositoriesService.GetWithFilter(params)
	assert.NoError(t, err, "Failed to get all repositories details")
	return data
}

func createDummyBuild(buildName string) error {
	dataArtifactoryBuild := &buildinfo.BuildInfo{
		Name:    buildName,
		Number:  "1.0.0",
		Started: "2014-09-30T12:00:19.893+0300",
		Modules: []buildinfo.Module{{
			Id: "example-module",
			Artifacts: []buildinfo.Artifact{
				{
					Type: "gz",
					Name: "c.tar.gz",
					Checksum: &buildinfo.Checksum{
						Sha1: "9d4336ff7bc2d2348aee4e27ad55e42110df4a80",
						Md5:  "b4918187cc9b3bf1b0772546d9398d7d",
					},
				},
			},
		}},
	}
	_, err := testBuildInfoService.PublishBuildInfo(dataArtifactoryBuild, "")
	return err
}

func deleteBuild(buildName string) error {
	err := deleteBuildIndex(buildName)
	if err != nil {
		return err
	}

	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}

	resp, _, err := client.SendDelete(artDetails.GetUrl()+"api/build/"+buildName+"?deleteAll=1", nil, artHTTPDetails, "")

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to delete build " + resp.Status)
	}

	return nil
}

func getIndexedBuilds() ([]string, error) {
	xrayDetails := GetXrayDetails()
	artHTTPDetails := xrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &artHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return []string{}, err
	}

	resp, body, _, err := client.SendGet(xrayDetails.GetUrl()+"api/v1/binMgr/default/builds", true, artHTTPDetails, "")
	if err != nil {
		return []string{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []string{}, errors.New("failed to get build index " + resp.Status)
	}

	response := &indexedBuildsPayload{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return []string{}, err
	}

	return response.IndexedBuilds, nil
}

func deleteBuildIndex(buildName string) error {
	// Prepare new indexed builds list
	indexedBuilds, err := getIndexedBuilds()
	if err != nil {
		return err
	}
	buildIndex := indexOf(buildName, indexedBuilds)
	if buildIndex == -1 {
		// Build indexing does not exist
		return nil
	}
	indexedBuilds = append(indexedBuilds[:buildIndex], indexedBuilds[buildIndex+1:]...)

	// Delete build index
	xrayDetails := GetXrayDetails()
	artHTTPDetails := xrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &artHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}

	dataIndexBuild := indexedBuildsPayload{IndexedBuilds: indexedBuilds}
	requestContentIndexBuild, err := json.Marshal(dataIndexBuild)

	resp, _, err := client.SendPut(xrayDetails.GetUrl()+"api/v1/binMgr/default/builds", requestContentIndexBuild, artHTTPDetails, "")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to delete build index " + resp.Status)
	}

	return nil
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

type indexedBuildsPayload struct {
	BinMgrId         string   `json:"bin_mgr_id,omitempty"`
	IndexedBuilds    []string `json:"indexed_builds"`
	NonIndexedBuilds []string `json:"non_indexed_builds,omitempty"`
}

// Verify sha256 is valid (a string size 256 characters) and not an empty string.
func verifyValidSha256(t *testing.T, sha256 string) {
	assert.Equal(t, 64, len(sha256), "Invalid sha256 : \""+sha256+"\"\nexpected length is 64 digit.")
}

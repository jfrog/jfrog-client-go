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

	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"
	"github.com/jfrog/jfrog-client-go/auth"
	distributionAuth "github.com/jfrog/jfrog-client-go/distribution/auth"
	distributionServices "github.com/jfrog/jfrog-client-go/distribution/services"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/mholt/archiver"
	"github.com/stretchr/testify/assert"
)

var RtUrl *string
var DistUrl *string
var RtUser *string
var RtPassword *string
var RtApiKey *string
var RtSshKeyPath *string
var RtSshPassphrase *string
var RtAccessToken *string
var LogLevel *string

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
var testsGetRepositoryService *services.GetRepositoryService
var testsCreateReplicationService *services.CreateReplicationService
var testsUpdateReplicationService *services.UpdateReplicationService
var testsReplicationGetService *services.GetReplicationService
var testsReplicationDeleteService *services.DeleteReplicationService

// Distribution services
var testsBundleSetSigningKeyService *distributionServices.SetSigningKeyService
var testsBundleCreateService *distributionServices.CreateReleaseBundleService
var testsBundleUpdateService *distributionServices.UpdateReleaseBundleService
var testsBundleSignService *distributionServices.SignBundleService
var testsBundleDistributeService *distributionServices.DistributeReleaseBundleService
var testsBundleDeleteLocalService *distributionServices.DeleteLocalReleaseBundleService
var testsBundleDeleteRemoteService *distributionServices.DeleteReleaseBundleService

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
	RtUrl = flag.String("rt.url", "http://localhost:8081/artifactory/", "Artifactory url")
	DistUrl = flag.String("rt.distUrl", "", "Distribution url")
	RtUser = flag.String("rt.user", "admin", "Artifactory username")
	RtPassword = flag.String("rt.password", "password", "Artifactory password")
	RtApiKey = flag.String("rt.apikey", "", "Artifactory user API key")
	RtSshKeyPath = flag.String("rt.sshKeyPath", "", "Ssh key file path")
	RtSshPassphrase = flag.String("rt.sshPassphrase", "", "Ssh key passphrase")
	RtAccessToken = flag.String("rt.accessToken", "", "Artifactory access token")
	LogLevel = flag.String("log-level", "INFO", "Sets the log level")
}

func createArtifactorySecurityManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsSecurityService = services.NewSecurityService(client)
	testsSecurityService.ArtDetails = artDetails
}

func createArtifactorySearchManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsSearchService = services.NewSearchService(client)
	testsSearchService.ArtDetails = artDetails
}

func createArtifactoryDeleteManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsDeleteService = services.NewDeleteService(client)
	testsDeleteService.SetThreads(3)
	testsDeleteService.ArtDetails = artDetails
}

func createArtifactoryUploadManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsUploadService = services.NewUploadService(client)
	testsUploadService.ArtDetails = artDetails
	testsUploadService.Threads = 3
}

func createArtifactoryDownloadManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsDownloadService = services.NewDownloadService(client)
	testsDownloadService.ArtDetails = artDetails
	testsDownloadService.SetThreads(3)
}

func createDistributionManager() {
	distDetails := GetDistDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&distDetails).Build()
	failOnHttpClientCreation(err)
	testsBundleCreateService = distributionServices.NewCreateReleseBundleService(client)
	testsBundleUpdateService = distributionServices.NewUpdateReleaseBundleService(client)
	testsBundleSignService = distributionServices.NewSignBundleService(client)
	testsBundleDistributeService = distributionServices.NewDistributeReleaseBundleService(client)
	testsBundleDeleteLocalService = distributionServices.NewDeleteLocalDistributionService(client)
	testsBundleSetSigningKeyService = distributionServices.NewSetSigningKeyService(client)
	testsBundleDeleteRemoteService = distributionServices.NewDeleteReleaseBundleService(client)
	testsBundleCreateService.DistDetails = distDetails
	testsBundleUpdateService.DistDetails = distDetails
	testsBundleSignService.DistDetails = distDetails
	testsBundleDistributeService.DistDetails = distDetails
	testsBundleDeleteLocalService.DistDetails = distDetails
	testsBundleSetSigningKeyService.DistDetails = distDetails
	testsBundleDeleteRemoteService.DistDetails = distDetails
}

func createArtifactoryCreateLocalRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsCreateLocalRepositoryService = services.NewLocalRepositoryService(client, false)
	testsCreateLocalRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateLocalRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsUpdateLocalRepositoryService = services.NewLocalRepositoryService(client, true)
	testsUpdateLocalRepositoryService.ArtDetails = artDetails
}

func createArtifactoryCreateRemoteRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsCreateRemoteRepositoryService = services.NewRemoteRepositoryService(client, false)
	testsCreateRemoteRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateRemoteRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsUpdateRemoteRepositoryService = services.NewRemoteRepositoryService(client, true)
	testsUpdateRemoteRepositoryService.ArtDetails = artDetails
}

func createArtifactoryCreateVirtualRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsCreateVirtualRepositoryService = services.NewVirtualRepositoryService(client, false)
	testsCreateVirtualRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateVirtualRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsUpdateVirtualRepositoryService = services.NewVirtualRepositoryService(client, true)
	testsUpdateVirtualRepositoryService.ArtDetails = artDetails
}

func createArtifactoryDeleteRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsDeleteRepositoryService = services.NewDeleteRepositoryService(client)
	testsDeleteRepositoryService.ArtDetails = artDetails
}

func createArtifactoryGetRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsGetRepositoryService = services.NewGetRepositoryService(client)
	testsGetRepositoryService.ArtDetails = artDetails
}

func createArtifactoryReplicationCreateManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsCreateReplicationService = services.NewCreateReplicationService(client)
	testsCreateReplicationService.ArtDetails = artDetails
}

func createArtifactoryReplicationUpdateManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsUpdateReplicationService = services.NewUpdateReplicationService(client)
	testsUpdateReplicationService.ArtDetails = artDetails
}

func createArtifactoryReplicationGetManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsReplicationGetService = services.NewGetReplicationService(client)
	testsReplicationGetService.ArtDetails = artDetails
}

func createArtifactoryReplicationDeleteManager() {
	artDetails := GetRtDetails()
	client, err := rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&artDetails).Build()
	failOnHttpClientCreation(err)
	testsReplicationDeleteService = services.NewDeleteReplicationService(client)
	testsReplicationDeleteService.ArtDetails = artDetails
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

func setAuthenticationDetail(details auth.ServiceDetails) {
	if !fileutils.IsSshUrl(details.GetUrl()) {
		if *RtApiKey != "" {
			details.SetApiKey(*RtApiKey)
		} else if *RtAccessToken != "" {
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
	_, uploaded, failed, err := testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "b.in"}
	up.Flat = true
	_, uploaded, failed, err = testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	archivePath := filepath.Join(workingDir, "c.tar.gz")
	err = archiver.TarGz.Make(archivePath, []string{filepath.Join(workingDir, "out/a.in")})
	if err != nil {
		t.Error(err)
	}
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: archivePath, Recursive: true, Target: RtTargetRepo}
	up.Flat = true
	_, uploaded, failed, err = testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
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
	deleteItems := make([]utils.ResultItem, len(toDelete))
	for i, item := range toDelete {
		deleteItems[i] = item
	}
	testsDeleteService.SetThreads(3)
	deletedCount, err := testsDeleteService.DeleteFiles(deleteItems)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(toDelete) != deletedCount {
		t.Errorf("Failed to delete files from Artifactory expected %d items to be deleted got %d.", len(toDelete), deletedCount)
	}
}

func createReposIfNeeded() error {
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
	resp, _, _, err := client.SendGet(artDetails.GetUrl()+RepoDetailsUrl+repoName, true, artHttpDetails)
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
	resp, _, err := client.SendPut(artDetails.GetUrl()+"api/repositories/"+repoName, content, artHttpDetails)
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
	return filepath.Join(dir, "testsdata")
}

func FixWinPath(filePath string) string {
	fixedPath := strings.Replace(filePath, "\\", "\\\\", -1)
	return fixedPath
}

func getRepoConfig(repoKey string) (body []byte) {
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil
	}
	resp, body, _, err := client.SendGet(artDetails.GetUrl()+"api/repositories/"+repoKey, false, artHttpDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	return
}

func validateRepoConfig(t *testing.T, repoKey string, params interface{}) {
	config := getRepoConfig(repoKey)
	if config == nil {
		t.Error("Failed to get repository config")
	}
	var confMap, paramsMap map[string]interface{}
	if err := json.Unmarshal(config, &confMap); err != nil {
		t.Error("Failed to marshal config map")
	}
	tmpJson, err := json.Marshal(params)
	if err != nil {
		t.Error("Failed to marshal repository params")
	}
	if err := json.Unmarshal(tmpJson, &paramsMap); err != nil {
		t.Error("Failed to unmarshal repository params")
	}
	for key, value := range paramsMap {
		assert.Equal(t, confMap[key], value, "faild validating "+key)
	}
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
	data, err := testsGetRepositoryService.Get(repoKey)
	assert.NoError(t, err, "Failed to get "+repoKey+" details")
	return data
}

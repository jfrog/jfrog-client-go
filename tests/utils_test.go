package tests

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientTestUtils "github.com/jfrog/jfrog-client-go/utils/tests"

	buildinfo "github.com/jfrog/build-info-go/entities"

	accessAuth "github.com/jfrog/jfrog-client-go/access/auth"
	accessServices "github.com/jfrog/jfrog-client-go/access/services"
	pipelinesAuth "github.com/jfrog/jfrog-client-go/pipelines/auth"
	pipelinesServices "github.com/jfrog/jfrog-client-go/pipelines/services"

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

var (
	TestArtifactory          *bool
	TestDistribution         *bool
	TestXray                 *bool
	TestPipelines            *bool
	TestAccess               *bool
	TestRepositories         *bool
	RtUrl                    *string
	DistUrl                  *string
	XrayUrl                  *string
	PipelinesUrl             *string
	RtUser                   *string
	RtPassword               *string
	RtApiKey                 *string
	RtSshKeyPath             *string
	RtSshPassphrase          *string
	PipelinesAccessToken     *string
	PipelinesVcsToken        *string
	PipelinesVcsRepoFullPath *string
	PipelinesVcsBranch       *string
	AccessUrl                *string
	AccessToken              *string
	ciRunId                  *string

	// Artifactory services
	testsUploadService                    *services.UploadService
	testsSearchService                    *services.SearchService
	testsDeleteService                    *services.DeleteService
	testsDownloadService                  *services.DownloadService
	testsSecurityService                  *services.SecurityService
	testsCreateLocalRepositoryService     *services.LocalRepositoryService
	testsCreateRemoteRepositoryService    *services.RemoteRepositoryService
	testsCreateVirtualRepositoryService   *services.VirtualRepositoryService
	testsCreateFederatedRepositoryService *services.FederatedRepositoryService
	testsUpdateLocalRepositoryService     *services.LocalRepositoryService
	testsUpdateRemoteRepositoryService    *services.RemoteRepositoryService
	testsUpdateVirtualRepositoryService   *services.VirtualRepositoryService
	testsUpdateFederatedRepositoryService *services.FederatedRepositoryService
	testsDeleteRepositoryService          *services.DeleteRepositoryService
	testsRepositoriesService              *services.RepositoriesService
	testsCreateReplicationService         *services.CreateReplicationService
	testsUpdateReplicationService         *services.UpdateReplicationService
	testsReplicationGetService            *services.GetReplicationService
	testsReplicationDeleteService         *services.DeleteReplicationService
	testsPermissionTargetService          *services.PermissionTargetService
	testUserService                       *services.UserService
	testGroupService                      *services.GroupService
	testBuildInfoService                  *services.BuildInfoService
	testsFederationService                *services.FederationService
	testsSystemService                    *services.SystemService
	testsStorageService                   *services.StorageService

	// Distribution services
	testsBundleSetSigningKeyService      *distributionServices.SetSigningKeyService
	testsBundleCreateService             *distributionServices.CreateReleaseBundleService
	testsBundleUpdateService             *distributionServices.UpdateReleaseBundleService
	testsBundleSignService               *distributionServices.SignBundleService
	testsBundleDistributeService         *distributionServices.DistributeReleaseBundleService
	testsBundleDistributionStatusService *distributionServices.DistributionStatusService
	testsBundleDeleteLocalService        *distributionServices.DeleteLocalReleaseBundleService
	testsBundleDeleteRemoteService       *distributionServices.DeleteReleaseBundleService

	// Xray Services
	testsXrayWatchService  *xrayServices.WatchService
	testsXrayPolicyService *xrayServices.PolicyService
	testXrayBinMgrService  *xrayServices.BinMgrService

	// Pipelines Services
	testsPipelinesIntegrationsService *pipelinesServices.IntegrationsService
	testsPipelinesSourcesService      *pipelinesServices.SourcesService

	// Access Services
	testsAccessProjectService *accessServices.ProjectService
	testsAccessInviteService  *accessServices.InviteService
	testsAccessTokensService  *accessServices.TokenService

	timestamp    = time.Now().Unix()
	timestampStr = strconv.FormatInt(timestamp, 10)
	trueValue    = true
	falseValue   = false

	// Tests configuration
	RtTargetRepo = "client-go"
)

const (
	HttpClientCreationFailureMessage = "Failed while attempting to create HttpClient: %s"
)

func init() {
	ciRunId = flag.String("ci.runId", "", "A unique identifier used as a suffix to create repositories in the tests")
	TestArtifactory = flag.Bool("test.artifactory", false, "Test Artifactory")
	TestDistribution = flag.Bool("test.distribution", false, "Test distribution")
	TestXray = flag.Bool("test.xray", false, "Test xray")
	TestPipelines = flag.Bool("test.pipelines", false, "Test pipelines")
	TestAccess = flag.Bool("test.access", false, "Test access")
	TestRepositories = flag.Bool("test.repositories", false, "Test repositories in Artifactory")
	RtUrl = flag.String("rt.url", "", "Artifactory url")
	DistUrl = flag.String("ds.url", "", "Distribution url")
	XrayUrl = flag.String("xr.url", "", "Xray url")
	PipelinesUrl = flag.String("pipe.url", "", "Pipelines url")
	RtUser = flag.String("rt.user", "", "Artifactory username")
	RtPassword = flag.String("rt.password", "", "Artifactory password")
	RtApiKey = flag.String("rt.apikey", "", "Artifactory user API key")
	RtSshKeyPath = flag.String("rt.sshKeyPath", "", "Ssh key file path")
	RtSshPassphrase = flag.String("rt.sshPassphrase", "", "Ssh key passphrase")
	PipelinesAccessToken = flag.String("pipe.accessToken", "", "Pipelines access token")
	PipelinesVcsToken = flag.String("pipe.vcsToken", "", "Vcs token for Pipelines tests")
	PipelinesVcsRepoFullPath = flag.String("pipe.vcsRepo", "", "Vcs full repo path for Pipelines tests")
	PipelinesVcsBranch = flag.String("pipe.vcsBranch", "", "Vcs branch for Pipelines tests")
	AccessUrl = flag.String("access.url", "", "Access url")
	AccessToken = flag.String("access.token", "", "Access token")
}

func getRtTargetRepoKey() string {
	return RtTargetRepo + "-" + getRunId()
}

func getRtTargetRepo() string {
	return getRtTargetRepoKey() + "/"
}

func getRunId() string {
	if ciRunId != nil && *ciRunId != "" {
		return *ciRunId + "-" + timestampStr
	}
	return timestampStr
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

func createArtifactoryCreateFederatedRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsCreateFederatedRepositoryService = services.NewFederatedRepositoryService(client, false)
	testsCreateFederatedRepositoryService.ArtDetails = artDetails
}

func createArtifactoryUpdateFederatedRepositoryManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsUpdateFederatedRepositoryService = services.NewFederatedRepositoryService(client, true)
	testsUpdateFederatedRepositoryService.ArtDetails = artDetails
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

func createArtifactoryFederationManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsFederationService = services.NewFederationService(client)
	testsFederationService.ArtDetails = artDetails
}

func createArtifactorySystemManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsSystemService = services.NewSystemService(artDetails, client)
}

func createArtifactoryStorageManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsStorageService = services.NewStorageService(artDetails, client)
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
		if *RtApiKey != "" {
			details.SetApiKey(*RtApiKey)
		} else if *AccessToken != "" {
			details.SetAccessToken(*AccessToken)
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
	defer clientTestUtils.RemoveAllAndAssert(t, workingDir)
	pattern := filepath.Join(workingDir, "*")
	up := services.NewUploadParams()
	targetProps, err := utils.ParseProperties("dummy=yes")
	if err != nil {
		t.Error(err)
	}
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo() + "test/", TargetProps: targetProps}
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
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo() + "b.in"}
	up.Flat = true
	summary, err = testsUploadService.UploadFiles(up)
	assert.NoError(t, err)
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
	up.CommonParams = &utils.CommonParams{Pattern: archivePath, Recursive: true, Target: getRtTargetRepo()}
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
	params := &utils.CommonParams{Pattern: getRtTargetRepo()}
	toDelete, err := testsDeleteService.GetPathsToDelete(services.DeleteParams{CommonParams: params})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer func() {
		assert.NoError(t, toDelete.Close())
	}()
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

func createRepo() error {
	if !(*TestArtifactory || *TestDistribution || *TestXray || *TestRepositories) {
		return nil
	}
	var err error
	repoKey := getRtTargetRepoKey()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, true)
	err = testsCreateLocalRepositoryService.Generic(glp)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func teardownIntegrationTests() {
	if !(*TestArtifactory || *TestDistribution || *TestXray || *TestRepositories) {
		return
	}
	repo := getRtTargetRepoKey()
	err := testsDeleteRepositoryService.Delete(repo)
	if err != nil {
		fmt.Printf("teardownIntegrationTests failed for:" + err.Error())
		os.Exit(1)
	}
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "testdata")
}

func setRepositoryBaseParams(params *services.RepositoryBaseParams, isUpdate bool) {
	if !isUpdate {
		params.ProjectKey = ""
		params.Environments = nil
		params.Description = strings.ToTitle(params.PackageType) + " Repo for jfrog-client-go local-repository-test"
		params.Notes = "Repo has been created"
		params.IncludesPattern = "dir1/*"
		params.ExcludesPattern = "dir2/*"
		params.RepoLayoutRef = "simple-default"
	} else {
		params.ProjectKey = ""
		params.Environments = nil
		params.Description += " - Updated"
		params.Notes = "Repo has been updated"
		params.IncludesPattern = ""
		params.ExcludesPattern = ""
		params.RepoLayoutRef = "build-default"
	}
}

func setAdditionalRepositoryBaseParams(params *services.AdditionalRepositoryBaseParams, isUpdate bool) {
	if !isUpdate {
		params.BlackedOut = &trueValue
		params.XrayIndex = &trueValue
		params.PropertySets = []string{"artifactory"}
		params.DownloadRedirect = &trueValue
		params.PriorityResolution = &trueValue
	} else {
		params.BlackedOut = &falseValue
		params.XrayIndex = &falseValue
		params.PropertySets = nil
		params.DownloadRedirect = &falseValue
		params.PriorityResolution = &falseValue
	}
}

func setCargoRepositoryParams(params *services.CargoRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.CargoAnonymousAccess = &trueValue
	} else {
		params.CargoAnonymousAccess = &falseValue
	}
}

func setDebianRepositoryParams(params *services.DebianRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.DebianTrivialLayout = &trueValue
		params.OptionalIndexCompressionFormats = []string{"bz2", "lzma"}
	} else {
		params.DebianTrivialLayout = &falseValue
		params.OptionalIndexCompressionFormats = nil
	}
}

func setDockerRepositoryParams(params *services.DockerRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.DockerApiVersion = "V1"
		params.MaxUniqueTags = 18
		params.BlockPushingSchema1 = &falseValue
		params.DockerTagRetention = 10
	} else {
		params.DockerApiVersion = "V2"
		params.MaxUniqueTags = 36
		params.BlockPushingSchema1 = &trueValue
		params.DockerTagRetention = 0
	}
}

func setJavaPackageManagersRepositoryParams(params *services.JavaPackageManagersRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.MaxUniqueSnapshots = 18
		params.HandleReleases = &trueValue
		params.HandleSnapshots = &trueValue
		params.SuppressPomConsistencyChecks = &trueValue
		params.SnapshotVersionBehavior = "non-unique"
		params.ChecksumPolicyType = "server-generated-checksums"
	} else {
		params.MaxUniqueSnapshots = 36
		params.HandleReleases = &falseValue
		params.HandleSnapshots = &falseValue
		params.SuppressPomConsistencyChecks = &falseValue
		params.SnapshotVersionBehavior = "unique"
		params.ChecksumPolicyType = "client-checksums"
	}
}

func setNugetRepositoryParams(params *services.NugetRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.ForceNugetAuthentication = &trueValue
		params.MaxUniqueSnapshots = 24
	} else {
		params.ForceNugetAuthentication = &falseValue
		params.MaxUniqueSnapshots = 18
	}
}

func setRpmRepositoryParams(params *services.RpmRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.YumRootDepth = 6
		params.CalculateYumMetadata = &trueValue
		params.EnableFileListsIndexing = &trueValue
		params.YumGroupFileNames = "filename"
	} else {
		params.YumRootDepth = 18
		params.CalculateYumMetadata = &falseValue
		params.EnableFileListsIndexing = &falseValue
		params.YumGroupFileNames = ""
	}
}

func getRepoConfig(repoKey string) ([]byte, error) {
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(artDetails.GetUrl()+"api/repositories/"+repoKey, false, artHttpDetails, "")
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return body, nil
}

func isEnterprisePlus() (bool, error) {
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return false, err
	}
	resp, body, _, err := client.SendGet(artDetails.GetUrl()+"api/system/license", false, artHttpDetails, "")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false, err
	}
	value, err := jsonparser.GetString(body, "type")
	if err != nil {
		return false, err
	}
	return value == "Enterprise Plus", nil
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
			// The password field may be encrypted and won't match the value set, need to handle this during validation
			if key == "password" {
				continue
			}
			// Download Redirect is only supported on Enterprise Plus. Expect false otherwise.
			if key == "downloadRedirect" {
				eplus, err := isEnterprisePlus()
				if err != nil {
					return false, err
				}
				if !eplus {
					expectedValue = false
				}
			}
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
		MaxRetries: 12,
		// RetriesIntervalMilliSecs in milliseconds
		RetriesIntervalMilliSecs: 10 * 1000,
		ErrorMessage:             "Waiting for Artifactory to evaluate repository operation...",
		ExecutionHandler:         createRepoConfigValidationFunc(repoKey, params),
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
	return fmt.Sprintf("%s-%d", getRtTargetRepoKey(), timestamp)
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

func isRepoExists(t *testing.T, repoKey string) bool {
	exists, err := testsRepositoriesService.IsExists(repoKey)
	assert.NoError(t, err, "Failed to check if "+repoKey+" exists")
	return exists
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
					Checksum: buildinfo.Checksum{
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
	if err != nil {
		return err
	}

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

func GetAccessDetails() auth.ServiceDetails {
	accessDetails := accessAuth.NewAccessDetails()
	accessDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*AccessUrl))
	accessDetails.SetAccessToken(*AccessToken)
	return accessDetails
}

func createAccessProjectManager() {
	accessDetails := GetAccessDetails()
	client, err := createJfrogHttpClient(&accessDetails)
	failOnHttpClientCreation(err)
	testsAccessProjectService = accessServices.NewProjectService(client)
	testsAccessProjectService.ServiceDetails = accessDetails

	artDetails := GetRtDetails()
	rtclient, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testGroupService = services.NewGroupService(rtclient)
	testGroupService.SetArtifactoryDetails(artDetails)
}

func createAccessInviteManager() {
	accessDetails := GetAccessDetails()
	client, err := createJfrogHttpClient(&accessDetails)
	failOnHttpClientCreation(err)
	testsAccessInviteService = accessServices.NewInviteService(client)
	testsAccessInviteService.ServiceDetails = accessDetails
	// To test "invite" flow we first have to create new "invited user" using ArtifactoryUserManager and Artifactory's API.
	createArtifactoryUserManager()
}

func createAccessTokensManager() {
	accessDetails := GetAccessDetails()
	client, err := createJfrogHttpClient(&accessDetails)
	failOnHttpClientCreation(err)
	testsAccessTokensService = accessServices.NewTokenService(client)
	testsAccessTokensService.ServiceDetails = accessDetails
}

func getUniqueField(prefix string) string {
	return strings.Join([]string{prefix, strconv.FormatInt(time.Now().Unix(), 10), runtime.GOOS}, "-")
}

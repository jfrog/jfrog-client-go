//go:build itest

package tests

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	testUtils "github.com/jfrog/jfrog-client-go/utils/tests"

	buildinfo "github.com/jfrog/build-info-go/entities"

	accessAuth "github.com/jfrog/jfrog-client-go/access/auth"
	accessServices "github.com/jfrog/jfrog-client-go/access/services"
	pipelinesAuth "github.com/jfrog/jfrog-client-go/pipelines/auth"
	pipelinesServices "github.com/jfrog/jfrog-client-go/pipelines/services"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"

	"github.com/jfrog/archiver/v3"
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
	xscAuth "github.com/jfrog/jfrog-client-go/xsc/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	TestArtifactory          *bool
	TestDistribution         *bool
	TestXray                 *bool
	TestCatalog              *bool
	TestXsc                  *bool
	TestPipelines            *bool
	TestAccess               *bool
	TestRepositories         *bool
	TestMultipartUpload      *bool
	TestUnit                 *bool
	PlatformUrl              *string
	RtUrl                    *string
	DistUrl                  *string
	XrayUrl                  *string
	XscUrl                   *string
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
	testsPackageService                   *services.PackageService
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
	testsAqlService                       *services.AqlService

	// Distribution services
	testsBundleSetSigningKeyService      *distributionServices.SetSigningKeyService
	testsBundleCreateService             *distributionServices.CreateReleaseBundleService
	testsBundleUpdateService             *distributionServices.UpdateReleaseBundleService
	testsBundleSignService               *distributionServices.SignBundleService
	testsBundleDistributeService         *distributionServices.DistributeReleaseBundleV1Service
	testsBundleDistributionStatusService *distributionServices.DistributionStatusService
	testsBundleDeleteLocalService        *distributionServices.DeleteLocalReleaseBundleService
	testsBundleDeleteRemoteService       *distributionServices.DeleteReleaseBundleService

	// Xray Services
	testsXrayWatchService      *xrayServices.WatchService
	testsXrayPolicyService     *xrayServices.PolicyService
	testXrayBinMgrService      *xrayServices.BinMgrService
	testsXrayIgnoreRuleService *xrayServices.IgnoreRuleService

	// Pipelines Services
	testsPipelinesIntegrationsService *pipelinesServices.IntegrationsService
	testsPipelinesSourcesService      *pipelinesServices.SourcesService
	testPipelinesRunService           *pipelinesServices.RunService
	testPipelinesSyncService          *pipelinesServices.SyncService
	testPipelinesSyncStatusService    *pipelinesServices.SyncStatusService

	// Access Services
	testsAccessPingService    *accessServices.PingService
	testsAccessProjectService *accessServices.ProjectService
	testsAccessInviteService  *accessServices.InviteService
	testsAccessTokensService  *accessServices.TokenService

	timestamp    = time.Now().Unix()
	timestampStr = strconv.FormatInt(timestamp, 10)

	// Tests configuration
	RtTargetRepo = "client-go"
)

const (
	HttpClientCreationFailureMessage = "Failed while attempting to create HttpClient: %s"
	buildNumber                      = "1.0.0"
	buildTimestamp                   = "1412067619893"
)

func init() {
	ciRunId = flag.String("ci.runId", "", "A unique identifier used as a suffix to create repositories in the tests")
	TestArtifactory = flag.Bool("test.artifactory", false, "Test Artifactory")
	TestDistribution = flag.Bool("test.distribution", false, "Test distribution")
	TestXray = flag.Bool("test.xray", false, "Test xray")
	TestCatalog = flag.Bool("test.catalog", false, "Test catalog")
	TestXsc = flag.Bool("test.xsc", false, "Test xsc")
	TestPipelines = flag.Bool("test.pipelines", false, "Test pipelines")
	TestAccess = flag.Bool("test.access", false, "Test access")
	TestRepositories = flag.Bool("test.repositories", false, "Test repositories in Artifactory")
	TestMultipartUpload = flag.Bool("test.mpu", false, "Test Artifactory multipart upload")
	TestUnit = flag.Bool("test.unit", false, "Run unit tests")
	PlatformUrl = flag.String("platform.url", "http://localhost:8082", "Platform url")
	RtUrl = flag.String("rt.url", "", "Artifactory url")
	DistUrl = flag.String("ds.url", "", "Distribution url")
	XrayUrl = flag.String("xr.url", "", "Xray url")
	XscUrl = flag.String("xsc.url", "", "Xsc url")
	PipelinesUrl = flag.String("pipe.url", "", "Pipelines url")
	AccessUrl = flag.String("access.url", "", "Access url")
	RtUser = flag.String("rt.user", "admin", "Artifactory username")
	RtPassword = flag.String("rt.password", "password", "Artifactory password")
	AccessToken = flag.String("access.token", testUtils.GetLocalArtifactoryTokenIfNeeded(*RtUrl), "Access token")
	RtApiKey = flag.String("rt.apikey", "", "Artifactory user API key")
	RtSshKeyPath = flag.String("rt.sshKeyPath", "", "Ssh key file path")
	RtSshPassphrase = flag.String("rt.sshPassphrase", "", "Ssh key passphrase")
	PipelinesAccessToken = flag.String("pipe.accessToken", "", "Pipelines access token")
	PipelinesVcsToken = flag.String("pipe.vcsToken", "", "Vcs token for Pipelines tests")
	PipelinesVcsRepoFullPath = flag.String("pipe.vcsRepo", "", "Vcs full repo path for Pipelines tests")
	PipelinesVcsBranch = flag.String("pipe.vcsBranch", "", "Vcs branch for Pipelines tests")
}

func checkFlags() {
	platformUrl := strings.TrimSuffix(*PlatformUrl, "/") + "/"

	if *RtUrl == "" {
		*RtUrl = platformUrl + "artifactory"
	}

	if *DistUrl == "" {
		*DistUrl = platformUrl + "distribution"
	}

	if *XrayUrl == "" {
		*XrayUrl = platformUrl + "xray"
	}

	if *XscUrl == "" {
		*XscUrl = platformUrl + "xsc"
	}

	if *PipelinesUrl == "" {
		*PipelinesUrl = platformUrl + "pipelines"
	}

	if *AccessUrl == "" {
		*AccessUrl = platformUrl + "access"
	}

	if *AccessToken == "" {
		*AccessToken = testUtils.GetLocalArtifactoryTokenIfNeeded(*RtUrl)
	}
}

func getRtTargetRepoKey() string {
	return RtTargetRepo + "-" + getRunId()
}

func getRtTargetRepo() string {
	return getRtTargetRepoKey() + "/"
}

// Get a run ID string used in the generated tests resources to prevent using same resources names in the test.
// Examples - Repository names, build-info names, Docker image names, Release Bundle names.
func getRunId() string {
	return getCustomRunId('-')
}

// Get a run ID string using a custom character. We use '-' for most of the resources names and '_' for JFrog Pipelines integration names.
func getCustomRunId(separator rune) string {
	if ciRunId != nil && *ciRunId != "" {
		return *ciRunId + string(separator) + timestampStr
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
	httpClientDetails := testsUploadService.ArtDetails.CreateHttpClientDetails()
	testsUploadService.MultipartUpload = utils.NewMultipartUpload(client, &httpClientDetails, testsUploadService.ArtDetails.GetUrl())
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
	testsBundleDistributeService = distributionServices.NewDistributeReleaseBundleV1Service(client)
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

func createArtifactoryGetPackageManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsPackageService = services.NewPackageService(client)
	testsPackageService.ArtDetails = artDetails
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

func createArtifactoryAqlManager() {
	artDetails := GetRtDetails()
	client, err := createJfrogHttpClient(&artDetails)
	failOnHttpClientCreation(err)
	testsAqlService = services.NewAqlService(artDetails, client)
}

func createJfrogHttpClient(artDetailsPtr *auth.ServiceDetails) (*jfroghttpclient.JfrogHttpClient, error) {
	artDetails := *artDetailsPtr
	return jfroghttpclient.JfrogClientBuilder().
		SetRetries(3).
		SetClientCertPath(artDetails.GetClientCertPath()).
		SetClientCertKeyPath(artDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(artDetails.RunPreRequestFunctions).
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
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testXrayBinMgrService = xrayServices.NewBinMgrService(client)
	testXrayBinMgrService.XrayDetails = xrayDetails
}

func createXrayIgnoreRuleManager() {
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testsXrayIgnoreRuleService = xrayServices.NewIgnoreRuleService(client)
	testsXrayIgnoreRuleService.XrayDetails = xrayDetails
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

func createPipelinesRunManager() {
	pipelinesDetails := GetPipelinesDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(pipelinesDetails.GetClientCertPath()).
		SetClientCertKeyPath(pipelinesDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(pipelinesDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testPipelinesRunService = pipelinesServices.NewRunService(client)
	testPipelinesRunService.ServiceDetails = pipelinesDetails
}

func createPipelinesSyncManager() {
	pipelinesDetails := GetPipelinesDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(pipelinesDetails.GetClientCertPath()).
		SetClientCertKeyPath(pipelinesDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(pipelinesDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testPipelinesSyncService = pipelinesServices.NewSyncService(client)
	testPipelinesSyncService.ServiceDetails = pipelinesDetails
}

func createPipelinesSyncStatusManager() {
	pipelinesDetails := GetPipelinesDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(pipelinesDetails.GetClientCertPath()).
		SetClientCertKeyPath(pipelinesDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(pipelinesDetails.RunPreRequestFunctions).
		Build()
	failOnHttpClientCreation(err)
	testPipelinesSyncStatusService = pipelinesServices.NewSyncStatusService(client)
	testPipelinesSyncStatusService.ServiceDetails = pipelinesDetails
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

// TODO this can be deleted after old Xsc service is deprecated from all services
func GetXscDetails() auth.ServiceDetails {
	xscDetails := xscAuth.NewXscDetails()
	xscDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*XscUrl))
	setAuthenticationDetail(xscDetails)
	return xscDetails
}

func GetPipelinesDetails() auth.ServiceDetails {
	pDetails := pipelinesAuth.NewPipelinesDetails()
	pDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*PipelinesUrl))
	pDetails.SetAccessToken(*PipelinesAccessToken)
	return pDetails
}

func setAuthenticationDetail(details auth.ServiceDetails) {
	if !fileutils.IsSshUrl(details.GetUrl()) {
		switch {
		case *RtApiKey != "":
			details.SetApiKey(*RtApiKey)
		case *AccessToken != "":
			details.SetAccessToken(*AccessToken)
		default:
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
	require.NoError(t, err)

	t.Cleanup(func() {
		testUtils.RemoveAllQuietly(t, workingDir)
	})

	pattern := filepath.Join(workingDir, "*")

	targetProps, err := utils.ParseProperties("dummy=yes")
	require.NoError(t, err)

	doUploadFile(t, pattern, "test/", targetProps)
	doUploadFile(t, pattern, "b.in", nil)

	archivePath := filepath.Join(workingDir, "c.tar.gz")
	err = archiver.Archive([]string{filepath.Join(workingDir, "out/a.in")}, archivePath)
	require.NoError(t, err)

	doUploadFile(t, archivePath, "", nil)
}

func doUploadFile(t *testing.T, pattern string, relativeTarget string, props *utils.Properties) {
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo() + relativeTarget, TargetProps: props}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	require.NoError(t, err)
	require.Equalf(t, 1, summary.TotalSucceeded, "Expected to upload 1 file.")
	require.Equalf(t, 0, summary.TotalFailed, "Failed to upload %d files.", summary.TotalFailed)
}

func artifactoryCleanup(t *testing.T) {
	params := &utils.CommonParams{Pattern: getRtTargetRepo()}

	toDelete, err := testsDeleteService.GetPathsToDelete(services.DeleteParams{CommonParams: params})
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to get paths to delete: %+v", err))
		return
	}

	defer testUtils.CloseQuietly(t, toDelete)

	numberOfItemToDelete, err := toDelete.Length()
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to get length of paths to delete: %+v", err))
		return
	}

	testsDeleteService.SetThreads(3)

	deletedCount, err := testsDeleteService.DeleteFiles(toDelete)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to delete files: %+v", err))
		return
	}

	if numberOfItemToDelete != deletedCount {
		log.Warn(fmt.Sprintf("Failed to delete files from Artifactory expected %d items to be deleted got %d.", numberOfItemToDelete, deletedCount))
	}
}

func createRepo(t *testing.T) {
	if !(*TestArtifactory || *TestDistribution || *TestXray || *TestRepositories || *TestMultipartUpload) {
		return
	}

	repoKey := getRtTargetRepoKey()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey

	exists, err := testsRepositoriesService.IsExists(repoKey)
	require.NoError(t, err)

	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, exists)

	err = testsCreateLocalRepositoryService.Generic(glp)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := testsDeleteRepositoryService.Delete(repoKey)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to delete repository %s: %+v", repoKey, err))
		}
	})
}

func getTestDataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "testdata")
}

func setRepositoryBaseParams(params *services.RepositoryBaseParams, isUpdate bool) {
	params.ProjectKey = ""
	params.Environments = nil
	// params.IncludesPattern = "**/*"
	// params.ExcludesPattern = "**/*"
	params.RepoLayoutRef = "simple-default"
	params.Description = strings.ToTitle(params.PackageType) + " Repo for jfrog-client-go local-repository-test"
	if isUpdate {
		params.Notes = "Repo has been updated"
	} else {
		params.Notes = "Repo has been created"
	}
}

func setAdditionalRepositoryBaseParams(params *services.AdditionalRepositoryBaseParams, isUpdate bool) {
	if isUpdate {
		params.BlackedOut = clientutils.Pointer(false)
		params.XrayIndex = clientutils.Pointer(false)
		params.PropertySets = nil
		params.DownloadRedirect = clientutils.Pointer(false)
		params.PriorityResolution = clientutils.Pointer(false)
	} else {
		params.BlackedOut = clientutils.Pointer(false)
		params.XrayIndex = clientutils.Pointer(true)
		params.PropertySets = []string{"artifactory"}
		params.DownloadRedirect = clientutils.Pointer(true)
		params.PriorityResolution = clientutils.Pointer(true)
	}
}

func setCargoRepositoryParams(params *services.CargoRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.CargoAnonymousAccess = clientutils.Pointer(true)
		params.CargoInternalIndex = clientutils.Pointer(true)
	} else {
		params.CargoAnonymousAccess = clientutils.Pointer(false)
		params.CargoInternalIndex = clientutils.Pointer(false)
	}
}

func setDebianRepositoryParams(params *services.DebianRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.DebianTrivialLayout = clientutils.Pointer(true)
		params.OptionalIndexCompressionFormats = []string{"bz2", "lzma"}
	} else {
		params.DebianTrivialLayout = clientutils.Pointer(false)
		params.OptionalIndexCompressionFormats = nil
	}
}

func setDockerRepositoryParams(params *services.DockerRepositoryParams, isUpdate bool) {
	if !isUpdate {
		maxUniqueTags := 18
		params.MaxUniqueTags = &maxUniqueTags
		dockerTagRetention := 10
		params.DockerTagRetention = &dockerTagRetention
		params.DockerApiVersion = "V1"
		params.BlockPushingSchema1 = clientutils.Pointer(false)
	} else {
		maxUniqueTags := 36
		params.MaxUniqueTags = &maxUniqueTags
		dockerTagRetention := 0
		params.DockerTagRetention = &dockerTagRetention
		params.DockerApiVersion = "V2"
		params.BlockPushingSchema1 = clientutils.Pointer(true)
	}
}

func setJavaPackageManagersRepositoryParams(params *services.JavaPackageManagersRepositoryParams, isUpdate bool) {
	if !isUpdate {
		maxUniqueTags := 18
		params.MaxUniqueSnapshots = &maxUniqueTags
		params.HandleReleases = clientutils.Pointer(true)
		params.HandleSnapshots = clientutils.Pointer(true)
		params.SuppressPomConsistencyChecks = clientutils.Pointer(true)
		params.SnapshotVersionBehavior = "non-unique"
		params.ChecksumPolicyType = "server-generated-checksums"
	} else {
		maxUniqueTags := 36
		params.MaxUniqueSnapshots = &maxUniqueTags
		params.HandleReleases = clientutils.Pointer(false)
		params.HandleSnapshots = clientutils.Pointer(false)
		params.SuppressPomConsistencyChecks = clientutils.Pointer(false)
		params.SnapshotVersionBehavior = "unique"
		params.ChecksumPolicyType = "client-checksums"
	}
}

func setNugetRepositoryParams(params *services.NugetRepositoryParams, isUpdate bool) {
	if !isUpdate {
		maxUniqueTags := 24
		params.MaxUniqueSnapshots = &maxUniqueTags
		params.ForceNugetAuthentication = clientutils.Pointer(true)
	} else {
		maxUniqueTags := 18
		params.MaxUniqueSnapshots = &maxUniqueTags
		params.ForceNugetAuthentication = clientutils.Pointer(false)
	}
}

func setRpmRepositoryParams(params *services.RpmRepositoryParams, isUpdate bool) {
	if !isUpdate {
		yumRootDepth := 6
		params.YumRootDepth = &yumRootDepth
		params.CalculateYumMetadata = clientutils.Pointer(true)
		params.EnableFileListsIndexing = clientutils.Pointer(true)
		params.YumGroupFileNames = "filename"
	} else {
		yumRootDepth := 18
		params.YumRootDepth = &yumRootDepth
		params.CalculateYumMetadata = clientutils.Pointer(false)
		params.EnableFileListsIndexing = clientutils.Pointer(false)
		params.YumGroupFileNames = ""
	}
}

func setTerraformRepositoryParams(params *services.TerraformRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.TerraformType = "provider"
	} else {
		params.TerraformType = "module"
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
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
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
	return strings.Contains(value, "Enterprise Plus"), nil
}

func createRepoConfigValidationFunc(repoKey string, expectedConfig interface{}) clientutils.ExecutionHandlerFunc {
	return func() (shouldRetry bool, err error) {
		config, err := getRepoConfig(repoKey)
		if err != nil || config == nil {
			err = errors.New("failed reading repository config for " + repoKey)
			return
		}
		var confMap, expectedConfigMap map[string]interface{}
		if err = json.Unmarshal(config, &confMap); err != nil {
			err = errors.New("failed unmarshalling repository config for " + repoKey)
			return
		}
		tmpJson, err := json.Marshal(expectedConfig)
		if err != nil {
			err = errors.New("failed marshalling expected config for " + repoKey)
			return
		}
		if err = json.Unmarshal(tmpJson, &expectedConfigMap); err != nil {
			err = errors.New("failed unmarshalling expected config for " + repoKey)
			return
		}
		for key, expectedValue := range expectedConfigMap {
			// The password field may be encrypted and won't match the value set, need to handle this during validation
			if key == "password" {
				continue
			}
			// Download Redirect is only supported on Enterprise Plus. Expect false otherwise.
			if key == "downloadRedirect" {
				var eplus bool
				eplus, err = isEnterprisePlus()
				if err != nil {
					return
				}
				if !eplus {
					expectedValue = false
				}
			}
			if !assert.ObjectsAreEqual(confMap[key], expectedValue) {
				err = fmt.Errorf("config validation for '%s' failed. key: '%s'\nexpected: '%s'\nactual: '%s'", repoKey, key, expectedValue, confMap[key])
				shouldRetry = true
				return
			}
		}
		return
	}
}

func validateRepoConfig(t *testing.T, repoKey string, params interface{}) {
	retryExecutor := &clientutils.RetryExecutor{
		MaxRetries: 5,
		// RetriesIntervalMilliSecs in milliseconds
		RetriesIntervalMilliSecs: 10 * 1000,
		ErrorMessage:             "Waiting for Artifactory to evaluate repository operation...",
		ExecutionHandler:         createRepoConfigValidationFunc(repoKey, params),
	}
	err := retryExecutor.Execute()
	assert.NoError(t, err)
}

func deleteRepoOnTestDone(t *testing.T, repoKey string) {
	t.Cleanup(func() {
		if err := testsDeleteRepositoryService.Delete(repoKey); err != nil {
			log.Warn(fmt.Sprintf("Failed to delete repository %s: %+v", repoKey, err))
		}
	})
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

func getAllRepos(t *testing.T, repoType string) *[]services.RepositoryDetails {
	params := services.NewRepositoriesFilterParams()
	params.RepoType = repoType
	params.PackageType = ""
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
		Number:  buildNumber,
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
	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}

	buildName = url.PathEscape(buildName)
	resp, _, err := client.SendDelete(artDetails.GetUrl()+"artifactory-build-info/"+buildName, nil, artHTTPDetails, "")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("failed to delete build " + resp.Status)
	}

	return nil
}

func getIndexedBuilds() ([]string, error) {
	xrayDetails := GetXrayDetails()
	artHTTPDetails := xrayDetails.CreateHttpClientDetails()
	artHTTPDetails.SetContentTypeApplicationJson()
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
	artHTTPDetails.SetContentTypeApplicationJson()
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

func createAccessPingManager() {
	accessDetails := GetAccessDetails()
	client, err := createJfrogHttpClient(&accessDetails)
	failOnHttpClientCreation(err)
	testsAccessPingService = accessServices.NewPingService(client)
	testsAccessPingService.ServiceDetails = accessDetails
}

func getUniqueField(prefix string) string {
	return strings.Join([]string{prefix, strconv.FormatInt(time.Now().Unix(), 10), runtime.GOOS}, "-")
}

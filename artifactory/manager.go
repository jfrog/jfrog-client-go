package artifactory

import (
	"io"

	"github.com/jfrog/jfrog-client-go/auth"

	buildinfo "github.com/jfrog/build-info-go/entities"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	_go "github.com/jfrog/jfrog-client-go/artifactory/services/go"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

type ArtifactoryServicesManagerImp struct {
	client   *jfroghttpclient.JfrogHttpClient
	config   config.Config
	progress ioutils.ProgressMgr
}

func New(config config.Config) (ArtifactoryServicesManager, error) {
	return NewWithProgress(config, nil)
}

func NewWithProgress(config config.Config, progress ioutils.ProgressMgr) (ArtifactoryServicesManager, error) {
	artDetails := config.GetServiceDetails()
	err := artDetails.InitSsh()
	if err != nil {
		return nil, err
	}
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetContext(config.GetContext()).
		SetDialTimeout(config.GetDialTimeout()).
		SetOverallRequestTimeout(config.GetOverallRequestTimeout()).
		SetClientCertPath(artDetails.GetClientCertPath()).
		SetClientCertKeyPath(artDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(artDetails.RunPreRequestFunctions).
		SetContext(config.GetContext()).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		SetHttpClient(config.GetHttpClient()).
		Build()
	if err != nil {
		return nil, err
	}
	if artDetails.GetClient() == nil {
		artDetails.SetClient(client)
	}
	manager, err := NewWithClient(config, client)
	if err != nil {
		return nil, err
	}
	manager.progress = progress
	return manager, err
}

func NewWithClient(config config.Config, client *jfroghttpclient.JfrogHttpClient) (*ArtifactoryServicesManagerImp, error) {
	manager := &ArtifactoryServicesManagerImp{config: config}
	manager.client = client
	return manager, nil
}

func (sm *ArtifactoryServicesManagerImp) CreateLocalRepository() *services.LocalRepositoryService {
	repositoryService := services.NewLocalRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) CreateRemoteRepository() *services.RemoteRepositoryService {
	repositoryService := services.NewRemoteRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) CreateVirtualRepository() *services.VirtualRepositoryService {
	repositoryService := services.NewVirtualRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) CreateFederatedRepository() *services.FederatedRepositoryService {
	repositoryService := services.NewFederatedRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) CreateLocalRepositoryWithParams(params services.LocalRepositoryBaseParams) error {
	repositoryService := services.NewRepositoriesService(sm.client)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService.Create(params, params.Key)
}

func (sm *ArtifactoryServicesManagerImp) CreateRemoteRepositoryWithParams(params services.RemoteRepositoryBaseParams) error {
	repositoryService := services.NewRepositoriesService(sm.client)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService.Create(params, params.Key)
}

func (sm *ArtifactoryServicesManagerImp) CreateVirtualRepositoryWithParams(params services.VirtualRepositoryBaseParams) error {
	repositoryService := services.NewRepositoriesService(sm.client)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService.Create(params, params.Key)
}

func (sm *ArtifactoryServicesManagerImp) CreateFederatedRepositoryWithParams(params services.FederatedRepositoryBaseParams) error {
	repositoryService := services.NewRepositoriesService(sm.client)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService.Create(params, params.Key)
}

func (sm *ArtifactoryServicesManagerImp) CreateRepositoryWithParams(params interface{}, repoName string) error {
	repositoryService := services.NewRepositoriesService(sm.client)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService.Create(params, repoName)
}

func (sm *ArtifactoryServicesManagerImp) UpdateLocalRepository() *services.LocalRepositoryService {
	repositoryService := services.NewLocalRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) UpdateRemoteRepository() *services.RemoteRepositoryService {
	repositoryService := services.NewRemoteRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) UpdateVirtualRepository() *services.VirtualRepositoryService {
	repositoryService := services.NewVirtualRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) UpdateFederatedRepository() *services.FederatedRepositoryService {
	repositoryService := services.NewFederatedRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManagerImp) UpdateRepositoryWithParams(params interface{}, repoName string) error {
	repositoryService := services.NewRepositoriesService(sm.client)
	repositoryService.ArtDetails = sm.config.GetServiceDetails()
	return repositoryService.Update(params, repoName)
}

func (sm *ArtifactoryServicesManagerImp) DeleteRepository(repoKey string) error {
	deleteRepositoryService := services.NewDeleteRepositoryService(sm.client)
	deleteRepositoryService.ArtDetails = sm.config.GetServiceDetails()
	return deleteRepositoryService.Delete(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) GetRepository(repoKey string, repoDetails interface{}) error {
	repositoriesService := services.NewRepositoriesService(sm.client)
	repositoriesService.ArtDetails = sm.config.GetServiceDetails()
	return repositoriesService.Get(repoKey, repoDetails)
}

func (sm *ArtifactoryServicesManagerImp) GetAllRepositories() (*[]services.RepositoryDetails, error) {
	repositoriesService := services.NewRepositoriesService(sm.client)
	repositoriesService.ArtDetails = sm.config.GetServiceDetails()
	return repositoriesService.GetAll()
}

func (sm *ArtifactoryServicesManagerImp) GetAllRepositoriesFiltered(params services.RepositoriesFilterParams) (*[]services.RepositoryDetails, error) {
	repositoriesService := services.NewRepositoriesService(sm.client)
	repositoriesService.ArtDetails = sm.config.GetServiceDetails()
	return repositoriesService.GetWithFilter(params)
}

func (sm *ArtifactoryServicesManagerImp) IsRepoExists(repoKey string) (bool, error) {
	repositoriesService := services.NewRepositoriesService(sm.client)
	repositoriesService.ArtDetails = sm.config.GetServiceDetails()
	return repositoriesService.IsExists(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) CreatePermissionTarget(params services.PermissionTargetParams) error {
	permissionTargetService := services.NewPermissionTargetService(sm.client)
	permissionTargetService.ArtDetails = sm.config.GetServiceDetails()
	return permissionTargetService.Create(params)
}

func (sm *ArtifactoryServicesManagerImp) UpdatePermissionTarget(params services.PermissionTargetParams) error {
	permissionTargetService := services.NewPermissionTargetService(sm.client)
	permissionTargetService.ArtDetails = sm.config.GetServiceDetails()
	return permissionTargetService.Update(params)
}

func (sm *ArtifactoryServicesManagerImp) DeletePermissionTarget(permissionTargetName string) error {
	permissionTargetService := services.NewPermissionTargetService(sm.client)
	permissionTargetService.ArtDetails = sm.config.GetServiceDetails()
	return permissionTargetService.Delete(permissionTargetName)
}

func (sm *ArtifactoryServicesManagerImp) GetPermissionTarget(permissionTargetName string) (*services.PermissionTargetParams, error) {
	permissionTargetService := services.NewPermissionTargetService(sm.client)
	permissionTargetService.ArtDetails = sm.config.GetServiceDetails()
	return permissionTargetService.Get(permissionTargetName)
}

func (sm *ArtifactoryServicesManagerImp) PublishBuildInfo(build *buildinfo.BuildInfo, projectKey string) (*clientutils.Sha256Summary, error) {
	buildInfoService := services.NewBuildInfoService(sm.config.GetServiceDetails(), sm.client)
	buildInfoService.DryRun = sm.config.IsDryRun()
	return buildInfoService.PublishBuildInfo(build, projectKey)
}

func (sm *ArtifactoryServicesManagerImp) DistributeBuild(params services.BuildDistributionParams) error {
	distributionService := services.NewDistributionService(sm.client)
	distributionService.DryRun = sm.config.IsDryRun()
	distributionService.ArtDetails = sm.config.GetServiceDetails()
	return distributionService.BuildDistribute(params)
}

func (sm *ArtifactoryServicesManagerImp) PromoteBuild(params services.PromotionParams) error {
	promotionService := services.NewPromotionService(sm.client)
	promotionService.DryRun = sm.config.IsDryRun()
	promotionService.ArtDetails = sm.config.GetServiceDetails()
	return promotionService.BuildPromote(params)
}

func (sm *ArtifactoryServicesManagerImp) DiscardBuilds(params services.DiscardBuildsParams) error {
	discardService := services.NewDiscardBuildsService(sm.client)
	discardService.ArtDetails = sm.config.GetServiceDetails()
	return discardService.DiscardBuilds(params)
}

func (sm *ArtifactoryServicesManagerImp) XrayScanBuild(params services.XrayScanParams) ([]byte, error) {
	xrayScanService := services.NewXrayScanService(sm.client)
	xrayScanService.ArtDetails = sm.config.GetServiceDetails()
	return xrayScanService.ScanBuild(params)
}

func (sm *ArtifactoryServicesManagerImp) GetPathsToDelete(params services.DeleteParams) (*content.ContentReader, error) {
	deleteService := services.NewDeleteService(sm.config.GetServiceDetails(), sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	return deleteService.GetPathsToDelete(params)
}

func (sm *ArtifactoryServicesManagerImp) DeleteFiles(reader *content.ContentReader) (int, error) {
	deleteService := services.NewDeleteService(sm.config.GetServiceDetails(), sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.Threads = sm.config.GetThreads()
	return deleteService.DeleteFiles(reader)
}

func (sm *ArtifactoryServicesManagerImp) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	readFileService := services.NewReadFileService(sm.config.GetServiceDetails(), sm.client)
	readFileService.DryRun = sm.config.IsDryRun()
	return readFileService.ReadRemoteFile(readPath)
}

func (sm *ArtifactoryServicesManagerImp) initDownloadService() *services.DownloadService {
	downloadService := services.NewDownloadService(sm.config.GetServiceDetails(), sm.client)
	downloadService.DryRun = sm.config.IsDryRun()
	downloadService.Threads = sm.config.GetThreads()
	downloadService.Progress = sm.progress
	return downloadService
}

func (sm *ArtifactoryServicesManagerImp) DownloadFiles(params ...services.DownloadParams) (totalDownloaded, totalFailed int, err error) {
	downloadService := sm.initDownloadService()
	summary, e := downloadService.DownloadFiles(params...)
	if e != nil {
		return 0, 0, e
	}
	return summary.TotalSucceeded, summary.TotalFailed, nil
}

func (sm *ArtifactoryServicesManagerImp) DownloadFilesWithSummary(params ...services.DownloadParams) (operationSummary *utils.OperationSummary, err error) {
	downloadService := sm.initDownloadService()
	downloadService.SetSaveSummary(true)
	return downloadService.DownloadFiles(params...)
}

func (sm *ArtifactoryServicesManagerImp) GetUnreferencedGitLfsFiles(params services.GitLfsCleanParams) (*content.ContentReader, error) {
	gitLfsCleanService := services.NewGitLfsCleanService(sm.config.GetServiceDetails(), sm.client)
	gitLfsCleanService.DryRun = sm.config.IsDryRun()
	return gitLfsCleanService.GetUnreferencedGitLfsFiles(params)
}

func (sm *ArtifactoryServicesManagerImp) SearchFiles(params services.SearchParams) (*content.ContentReader, error) {
	searchService := services.NewSearchService(sm.config.GetServiceDetails(), sm.client)
	return searchService.Search(params)
}

func (sm *ArtifactoryServicesManagerImp) Aql(aql string) (io.ReadCloser, error) {
	aqlService := services.NewAqlService(sm.config.GetServiceDetails(), sm.client)
	return aqlService.ExecAql(aql)
}

func (sm *ArtifactoryServicesManagerImp) SetProps(params services.PropsParams) (int, error) {
	setPropsService := services.NewPropsService(sm.client)
	setPropsService.ArtDetails = sm.config.GetServiceDetails()
	setPropsService.Threads = sm.config.GetThreads()
	return setPropsService.SetProps(params)
}

func (sm *ArtifactoryServicesManagerImp) DeleteProps(params services.PropsParams) (int, error) {
	setPropsService := services.NewPropsService(sm.client)
	setPropsService.ArtDetails = sm.config.GetServiceDetails()
	setPropsService.Threads = sm.config.GetThreads()
	return setPropsService.DeleteProps(params)
}

func (sm *ArtifactoryServicesManagerImp) GetItemProps(relativePath string) (*utils.ItemProperties, error) {
	setPropsService := services.NewPropsService(sm.client)
	setPropsService.ArtDetails = sm.config.GetServiceDetails()
	return setPropsService.GetItemProperties(relativePath)
}

func (sm *ArtifactoryServicesManagerImp) initUploadService() *services.UploadService {
	uploadService := services.NewUploadService(sm.client)
	uploadService.Threads = sm.config.GetThreads()
	uploadService.ArtDetails = sm.config.GetServiceDetails()
	uploadService.DryRun = sm.config.IsDryRun()
	uploadService.Progress = sm.progress
	httpClientDetails := uploadService.ArtDetails.CreateHttpClientDetails()
	uploadService.MultipartUpload = utils.NewMultipartUpload(sm.client, &httpClientDetails, uploadService.ArtDetails.GetUrl())
	return uploadService
}

func (sm *ArtifactoryServicesManagerImp) UploadFiles(params ...services.UploadParams) (totalUploaded, totalFailed int, err error) {
	uploadService := sm.initUploadService()
	summary, e := uploadService.UploadFiles(params...)
	if summary == nil {
		return 0, 0, e
	}
	return summary.TotalSucceeded, summary.TotalFailed, e
}

func (sm *ArtifactoryServicesManagerImp) UploadFilesWithSummary(params ...services.UploadParams) (operationSummary *utils.OperationSummary, err error) {
	uploadService := sm.initUploadService()
	uploadService.SetSaveSummary(true)
	return uploadService.UploadFiles(params...)
}

func (sm *ArtifactoryServicesManagerImp) Copy(params ...services.MoveCopyParams) (successCount, failedCount int, err error) {
	copyService := services.NewMoveCopyService(sm.config.GetServiceDetails(), sm.client, services.COPY)
	copyService.DryRun = sm.config.IsDryRun()
	copyService.Threads = sm.config.GetThreads()
	return copyService.MoveCopyServiceMoveFilesWrapper(params...)
}

func (sm *ArtifactoryServicesManagerImp) Move(params ...services.MoveCopyParams) (successCount, failedCount int, err error) {
	moveService := services.NewMoveCopyService(sm.config.GetServiceDetails(), sm.client, services.MOVE)
	moveService.DryRun = sm.config.IsDryRun()
	moveService.Threads = sm.config.GetThreads()
	return moveService.MoveCopyServiceMoveFilesWrapper(params...)
}

func (sm *ArtifactoryServicesManagerImp) PublishGoProject(params _go.GoParams) (*utils.OperationSummary, error) {
	goService := _go.NewGoService(sm.client)
	goService.ArtDetails = sm.config.GetServiceDetails()
	return goService.PublishPackage(params)
}

func (sm *ArtifactoryServicesManagerImp) Ping() ([]byte, error) {
	pingService := services.NewPingService(sm.config.GetServiceDetails(), sm.client)
	return pingService.Ping()
}

func (sm *ArtifactoryServicesManagerImp) GetConfig() config.Config {
	return sm.config
}

func (sm *ArtifactoryServicesManagerImp) GetBuildInfo(params services.BuildInfoParams) (*buildinfo.PublishedBuildInfo, bool, error) {
	buildInfoService := services.NewBuildInfoService(sm.config.GetServiceDetails(), sm.client)
	return buildInfoService.GetBuildInfo(params)
}

func (sm *ArtifactoryServicesManagerImp) CreateAPIKey() (string, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.CreateAPIKey()
}

func (sm *ArtifactoryServicesManagerImp) RegenerateAPIKey() (string, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.RegenerateAPIKey()
}

func (sm *ArtifactoryServicesManagerImp) GetAPIKey() (string, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.GetAPIKey()
}

func (sm *ArtifactoryServicesManagerImp) CreateToken(params services.CreateTokenParams) (auth.CreateTokenResponseData, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.CreateToken(params)
}

func (sm *ArtifactoryServicesManagerImp) GetTokens() (services.GetTokensResponseData, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.GetTokens()
}

func (sm *ArtifactoryServicesManagerImp) GetUserTokens(username string) ([]string, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.GetUserTokens(username)
}

func (sm *ArtifactoryServicesManagerImp) RefreshToken(params services.ArtifactoryRefreshTokenParams) (auth.CreateTokenResponseData, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.RefreshToken(params)
}

func (sm *ArtifactoryServicesManagerImp) RevokeToken(params services.RevokeTokenParams) (string, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetServiceDetails()
	return securityService.RevokeToken(params)
}

func (sm *ArtifactoryServicesManagerImp) CreateReplication(params services.CreateReplicationParams) error {
	replicationService := services.NewCreateReplicationService(sm.client)
	replicationService.ArtDetails = sm.config.GetServiceDetails()
	return replicationService.CreateReplication(params)
}

func (sm *ArtifactoryServicesManagerImp) UpdateReplication(params services.UpdateReplicationParams) error {
	replicationService := services.NewUpdateReplicationService(sm.client)
	replicationService.ArtDetails = sm.config.GetServiceDetails()
	return replicationService.UpdateReplication(params)
}

func (sm *ArtifactoryServicesManagerImp) DeleteReplication(repoKey string) error {
	deleteReplicationService := services.NewDeleteReplicationService(sm.client)
	deleteReplicationService.ArtDetails = sm.config.GetServiceDetails()
	return deleteReplicationService.DeleteReplication(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) GetReplication(repoKey string) ([]utils.ReplicationParams, error) {
	getPushReplicationService := services.NewGetReplicationService(sm.client)
	getPushReplicationService.ArtDetails = sm.config.GetServiceDetails()
	return getPushReplicationService.GetReplication(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) ConvertLocalToFederatedRepository(repoKey string) error {
	getFederationService := services.NewFederationService(sm.client)
	getFederationService.ArtDetails = sm.config.GetServiceDetails()
	return getFederationService.ConvertLocalToFederated(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) TriggerFederatedRepositoryFullSyncAll(repoKey string) error {
	getFederationService := services.NewFederationService(sm.client)
	getFederationService.ArtDetails = sm.config.GetServiceDetails()
	return getFederationService.TriggerFederatedFullSyncAll(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) TriggerFederatedRepositoryFullSyncMirror(repoKey string, mirrorUrl string) error {
	getFederationService := services.NewFederationService(sm.client)
	getFederationService.ArtDetails = sm.config.GetServiceDetails()
	return getFederationService.TriggerFederatedFullSyncMirror(repoKey, mirrorUrl)
}

func (sm *ArtifactoryServicesManagerImp) GetVersion() (string, error) {
	systemService := services.NewSystemService(sm.config.GetServiceDetails(), sm.client)
	return systemService.GetVersion()
}

func (sm *ArtifactoryServicesManagerImp) GetServiceId() (string, error) {
	systemService := services.NewSystemService(sm.config.GetServiceDetails(), sm.client)
	return systemService.GetServiceId()
}

func (sm *ArtifactoryServicesManagerImp) GetRunningNodes() ([]string, error) {
	systemService := services.NewSystemService(sm.config.GetServiceDetails(), sm.client)
	return systemService.GetRunningNodes()
}

func (sm *ArtifactoryServicesManagerImp) GetConfigDescriptor() (string, error) {
	systemService := services.NewSystemService(sm.config.GetServiceDetails(), sm.client)
	return systemService.GetConfigDescriptor()
}

func (sm *ArtifactoryServicesManagerImp) ActivateKeyEncryption() error {
	systemService := services.NewSystemService(sm.config.GetServiceDetails(), sm.client)
	return systemService.ActivateKeyEncryption()
}

func (sm *ArtifactoryServicesManagerImp) DeactivateKeyEncryption() (bool, error) {
	systemService := services.NewSystemService(sm.config.GetServiceDetails(), sm.client)
	return systemService.DeactivateKeyEncryption()
}

func (sm *ArtifactoryServicesManagerImp) GetGroup(params services.GroupParams) (*services.Group, error) {
	groupService := services.NewGroupService(sm.client)
	groupService.ArtDetails = sm.config.GetServiceDetails()
	return groupService.GetGroup(params)
}

func (sm *ArtifactoryServicesManagerImp) GetAllGroups() (*[]string, error) {
	groupService := services.NewGroupService(sm.client)
	groupService.ArtDetails = sm.config.GetServiceDetails()
	return groupService.GetAllGroups()
}

func (sm *ArtifactoryServicesManagerImp) CreateGroup(params services.GroupParams) error {
	groupService := services.NewGroupService(sm.client)
	groupService.ArtDetails = sm.config.GetServiceDetails()
	return groupService.CreateGroup(params)
}

func (sm *ArtifactoryServicesManagerImp) UpdateGroup(params services.GroupParams) error {
	groupService := services.NewGroupService(sm.client)
	groupService.ArtDetails = sm.config.GetServiceDetails()
	return groupService.UpdateGroup(params)
}

func (sm *ArtifactoryServicesManagerImp) DeleteGroup(name string) error {
	groupService := services.NewGroupService(sm.client)
	groupService.ArtDetails = sm.config.GetServiceDetails()
	return groupService.DeleteGroup(name)
}

func (sm *ArtifactoryServicesManagerImp) GetUser(params services.UserParams) (*services.User, error) {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.GetUser(params)
}

func (sm *ArtifactoryServicesManagerImp) GetAllUsers() ([]*services.User, error) {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.GetAllUsers()
}

func (sm *ArtifactoryServicesManagerImp) CreateUser(params services.UserParams) error {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.CreateUser(params)
}

func (sm *ArtifactoryServicesManagerImp) UpdateUser(params services.UserParams) error {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.UpdateUser(params)
}

func (sm *ArtifactoryServicesManagerImp) DeleteUser(name string) error {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.DeleteUser(name)
}

func (sm *ArtifactoryServicesManagerImp) GetLockedUsers() ([]string, error) {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.GetLockedUsers()
}

func (sm *ArtifactoryServicesManagerImp) UnlockUser(name string) error {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.UnlockUser(name)
}

func (sm *ArtifactoryServicesManagerImp) PromoteDocker(params services.DockerPromoteParams) error {
	systemService := services.NewDockerPromoteService(sm.config.GetServiceDetails(), sm.client)
	return systemService.PromoteDocker(params)
}

func (sm *ArtifactoryServicesManagerImp) Export(params services.ExportParams) error {
	exportService := services.NewExportService(sm.config.GetServiceDetails(), sm.client)
	return exportService.Export(params)
}

func (sm *ArtifactoryServicesManagerImp) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *ArtifactoryServicesManagerImp) FolderInfo(relativePath string) (*utils.FolderInfo, error) {
	storageService := services.NewStorageService(sm.config.GetServiceDetails(), sm.client)
	return storageService.FolderInfo(relativePath)
}

func (sm *ArtifactoryServicesManagerImp) FileInfo(relativePath string) (*utils.FileInfo, error) {
	storageService := services.NewStorageService(sm.config.GetServiceDetails(), sm.client)
	return storageService.FileInfo(relativePath)
}

func (sm *ArtifactoryServicesManagerImp) FileList(relativePath string, optionalParams utils.FileListParams) (*utils.FileListResponse, error) {
	storageService := services.NewStorageService(sm.config.GetServiceDetails(), sm.client)
	return storageService.FileList(relativePath, optionalParams)
}

func (sm *ArtifactoryServicesManagerImp) GetStorageInfo() (*utils.StorageInfo, error) {
	storageService := services.NewStorageService(sm.config.GetServiceDetails(), sm.client)
	return storageService.StorageInfo()
}

func (sm *ArtifactoryServicesManagerImp) CalculateStorageInfo() error {
	storageService := services.NewStorageService(sm.config.GetServiceDetails(), sm.client)
	return storageService.StorageInfoRefresh()
}

func (sm *ArtifactoryServicesManagerImp) ReleaseBundleImport(filePath string) error {
	releaseService := services.NewReleaseService(sm.config.GetServiceDetails(), sm.client)
	return releaseService.ImportReleaseBundle(filePath)
}

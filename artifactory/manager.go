package artifactory

import (
	"io"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	_go "github.com/jfrog/jfrog-client-go/artifactory/services/go"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

type ArtifactoryServicesManagerImp struct {
	client   *jfroghttpclient.JfrogHttpClient
	config   config.Config
	progress ioutils.ProgressMgr
}

func New(artDetails *auth.ServiceDetails, config config.Config) (ArtifactoryServicesManager, error) {
	return NewWithProgress(artDetails, config, nil)
}

func NewWithProgress(artDetails *auth.ServiceDetails, config config.Config, progress ioutils.ProgressMgr) (ArtifactoryServicesManager, error) {
	err := (*artDetails).InitSsh()
	if err != nil {
		return nil, err
	}
	manager := &ArtifactoryServicesManagerImp{config: config, progress: progress}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetServiceDetails(artDetails).
		SetContext(config.GetContext()).
		Build()
	if err != nil {
		return nil, err
	}
	return manager, err
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

func (sm *ArtifactoryServicesManagerImp) DeleteRepository(repoKey string) error {
	deleteRepositoryService := services.NewDeleteRepositoryService(sm.client)
	deleteRepositoryService.ArtDetails = sm.config.GetServiceDetails()
	return deleteRepositoryService.Delete(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) GetRepository(repoKey string) (*services.RepositoryDetails, error) {
	getRepositoryService := services.NewGetRepositoryService(sm.client)
	getRepositoryService.ArtDetails = sm.config.GetServiceDetails()
	return getRepositoryService.Get(repoKey)
}

func (sm *ArtifactoryServicesManagerImp) GetAllRepositories() (*[]services.RepositoryDetails, error) {
	getRepositoryService := services.NewGetRepositoryService(sm.client)
	getRepositoryService.ArtDetails = sm.config.GetServiceDetails()
	return getRepositoryService.GetAll()
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

func (sm *ArtifactoryServicesManagerImp) PublishBuildInfo(build *buildinfo.BuildInfo, projectKey string) error {
	buildInfoService := services.NewBuildInfoService(sm.client)
	buildInfoService.DryRun = sm.config.IsDryRun()
	buildInfoService.ArtDetails = sm.config.GetServiceDetails()
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
	deleteService := services.NewDeleteService(sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.ArtDetails = sm.config.GetServiceDetails()
	return deleteService.GetPathsToDelete(params)
}

func (sm *ArtifactoryServicesManagerImp) DeleteFiles(reader *content.ContentReader) (int, error) {
	deleteService := services.NewDeleteService(sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.ArtDetails = sm.config.GetServiceDetails()
	deleteService.Threads = sm.config.GetThreads()
	return deleteService.DeleteFiles(reader)
}

func (sm *ArtifactoryServicesManagerImp) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	readFileService := services.NewReadFileService(sm.client)
	readFileService.DryRun = sm.config.IsDryRun()
	readFileService.ArtDetails = sm.config.GetServiceDetails()
	return readFileService.ReadRemoteFile(readPath)
}

func (sm *ArtifactoryServicesManagerImp) initDownloadService() *services.DownloadService {
	downloadService := services.NewDownloadService(sm.client)
	downloadService.DryRun = sm.config.IsDryRun()
	downloadService.ArtDetails = sm.config.GetServiceDetails()
	downloadService.Threads = sm.config.GetThreads()
	downloadService.Progress = sm.progress
	return downloadService
}

func (sm *ArtifactoryServicesManagerImp) DownloadFiles(params ...services.DownloadParams) (totalDownloaded, totalExpected int, err error) {
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
	gitLfsCleanService := services.NewGitLfsCleanService(sm.client)
	gitLfsCleanService.DryRun = sm.config.IsDryRun()
	gitLfsCleanService.ArtDetails = sm.config.GetServiceDetails()
	return gitLfsCleanService.GetUnreferencedGitLfsFiles(params)
}

func (sm *ArtifactoryServicesManagerImp) SearchFiles(params services.SearchParams) (*content.ContentReader, error) {
	searchService := services.NewSearchService(sm.client)
	searchService.ArtDetails = sm.config.GetServiceDetails()
	return searchService.Search(params)
}

func (sm *ArtifactoryServicesManagerImp) Aql(aql string) (io.ReadCloser, error) {
	aqlService := services.NewAqlService(sm.client)
	aqlService.ArtDetails = sm.config.GetServiceDetails()
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
func (sm *ArtifactoryServicesManagerImp) initUploadService() *services.UploadService {
	uploadService := services.NewUploadService(sm.client)
	uploadService.Threads = sm.config.GetThreads()
	uploadService.ArtDetails = sm.config.GetServiceDetails()
	uploadService.DryRun = sm.config.IsDryRun()
	uploadService.Progress = sm.progress
	return uploadService
}

func (sm *ArtifactoryServicesManagerImp) UploadFiles(params ...services.UploadParams) (totalUploaded, totalFailed int, err error) {
	uploadService := sm.initUploadService()
	summary, e := uploadService.UploadFiles(params...)
	if e != nil {
		return 0, 0, e
	}
	return summary.TotalSucceeded, summary.TotalFailed, nil
}

func (sm *ArtifactoryServicesManagerImp) UploadFilesWithSummary(params ...services.UploadParams) (operationSummary *utils.OperationSummary, err error) {
	uploadService := sm.initUploadService()
	uploadService.SetSaveSummary(true)
	return uploadService.UploadFiles(params...)
}

func (sm *ArtifactoryServicesManagerImp) Copy(params ...services.MoveCopyParams) (successCount, failedCount int, err error) {
	copyService := services.NewMoveCopyService(sm.client, services.COPY)
	copyService.DryRun = sm.config.IsDryRun()
	copyService.ArtDetails = sm.config.GetServiceDetails()
	copyService.Threads = sm.config.GetThreads()
	return copyService.MoveCopyServiceMoveFilesWrapper(params...)
}

func (sm *ArtifactoryServicesManagerImp) Move(params ...services.MoveCopyParams) (successCount, failedCount int, err error) {
	moveService := services.NewMoveCopyService(sm.client, services.MOVE)
	moveService.DryRun = sm.config.IsDryRun()
	moveService.ArtDetails = sm.config.GetServiceDetails()
	moveService.Threads = sm.config.GetThreads()
	return moveService.MoveCopyServiceMoveFilesWrapper(params...)
}

func (sm *ArtifactoryServicesManagerImp) PublishGoProject(params _go.GoParams) error {
	goService := _go.NewGoService(sm.client)
	goService.ArtDetails = sm.config.GetServiceDetails()
	return goService.PublishPackage(params)
}

func (sm *ArtifactoryServicesManagerImp) Ping() ([]byte, error) {
	pingService := services.NewPingService(sm.client)
	pingService.ArtDetails = sm.config.GetServiceDetails()
	return pingService.Ping()
}

func (sm *ArtifactoryServicesManagerImp) GetConfig() config.Config {
	return sm.config
}

func (sm *ArtifactoryServicesManagerImp) GetBuildInfo(params services.BuildInfoParams) (*buildinfo.PublishedBuildInfo, bool, error) {
	buildInfoService := services.NewBuildInfoService(sm.client)
	buildInfoService.ArtDetails = sm.config.GetServiceDetails()
	return buildInfoService.GetBuildInfo(params)
}

func (sm *ArtifactoryServicesManagerImp) CreateToken(params services.CreateTokenParams) (services.CreateTokenResponseData, error) {
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

func (sm *ArtifactoryServicesManagerImp) RefreshToken(params services.RefreshTokenParams) (services.CreateTokenResponseData, error) {
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

func (sm *ArtifactoryServicesManagerImp) GetVersion() (string, error) {
	systemService := services.NewSystemService(sm.client)
	systemService.ArtDetails = sm.config.GetServiceDetails()
	return systemService.GetVersion()
}

func (sm *ArtifactoryServicesManagerImp) GetServiceId() (string, error) {
	systemService := services.NewSystemService(sm.client)
	systemService.ArtDetails = sm.config.GetServiceDetails()
	return systemService.GetServiceId()
}

func (sm *ArtifactoryServicesManagerImp) GetGroup(params services.GroupParams) (*services.Group, error) {
	groupService := services.NewGroupService(sm.client)
	groupService.ArtDetails = sm.config.GetServiceDetails()
	return groupService.GetGroup(params)
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

func (sm *ArtifactoryServicesManagerImp) PromoteDocker(params services.DockerPromoteParams) error {
	systemService := services.NewDockerPromoteService(sm.client)
	systemService.ArtDetails = sm.config.GetServiceDetails()
	return systemService.PromoteDocker(params)
}

func (sm *ArtifactoryServicesManagerImp) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

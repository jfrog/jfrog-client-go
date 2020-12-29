package artifactory

import (
	"io"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	_go "github.com/jfrog/jfrog-client-go/artifactory/services/go"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

type ArtifactoryServicesManagerImp struct {
	client   *rthttpclient.ArtifactoryHttpClient
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
	manager.client, err = rthttpclient.ArtifactoryClientBuilder().
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

func (sm *ArtifactoryServicesManagerImp) PublishBuildInfo(build *buildinfo.BuildInfo, project string) error {
	buildInfoService := services.NewBuildInfoService(sm.client)
	buildInfoService.DryRun = sm.config.IsDryRun()
	buildInfoService.Project = project
	buildInfoService.ArtDetails = sm.config.GetServiceDetails()
	return buildInfoService.PublishBuildInfo(build)
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
	return downloadService.DownloadFiles(params...)
}

func (sm *ArtifactoryServicesManagerImp) DownloadFilesWithResultReader(params ...services.DownloadParams) (resultReader *content.ContentReader, totalDownloaded, totalExpected int, err error) {
	downloadService := sm.initDownloadService()
	rw, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return
	}
	defer rw.Close()
	downloadService.ResultWriter = rw
	totalDownloaded, totalExpected, err = downloadService.DownloadFiles(params...)
	if err != nil {
		return
	}
	resultReader = content.NewContentReader(downloadService.ResultWriter.GetFilePath(), content.DefaultKey)
	return
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
	return uploadService.UploadFiles(params...)
}

func (sm *ArtifactoryServicesManagerImp) UploadFilesWithResultReader(params ...services.UploadParams) (resultReader *content.ContentReader, totalUploaded, totalFailed int, err error) {
	uploadService := sm.initUploadService()
	resultWriter, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return
	}
	defer resultWriter.Close()
	uploadService.ResultWriter = resultWriter
	totalUploaded, totalFailed, err = uploadService.UploadFiles(params...)
	if err != nil {
		return
	}
	resultReader = content.NewContentReader(uploadService.ResultWriter.GetFilePath(), content.DefaultKey)
	return
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

func (sm *ArtifactoryServicesManagerImp) GetUser(name string) (*services.User, error) {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.GetUser(name)
}

func (sm *ArtifactoryServicesManagerImp) CreateUser(user services.User) error {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.CreateOrUpdateUser(user)
}

func (sm *ArtifactoryServicesManagerImp) DeleteUser(name string) error {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.DeleteUser(name)
}

func (sm *ArtifactoryServicesManagerImp) UserExists(name string) (bool, error) {
	userService := services.NewUserService(sm.client)
	userService.ArtDetails = sm.config.GetServiceDetails()
	return userService.UserExists(name)
}

func (sm *ArtifactoryServicesManagerImp) PromoteDocker(params services.DockerPromoteParams) error {
	systemService := services.NewDockerPromoteService(sm.client)
	systemService.ArtDetails = sm.config.GetServiceDetails()
	return systemService.PromoteDocker(params)
}

func (sm *ArtifactoryServicesManagerImp) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}

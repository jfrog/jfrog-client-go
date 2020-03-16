package artifactory

import (
	"io"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/go"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
)

type ArtifactoryServicesManager struct {
	client   *rthttpclient.ArtifactoryHttpClient
	config   config.Config
	progress ioutils.Progress
}

func New(artDetails *auth.CommonDetails, config config.Config) (*ArtifactoryServicesManager, error) {
	return NewWithProgress(artDetails, config, nil)
}

func NewWithProgress(artDetails *auth.CommonDetails, config config.Config, progress ioutils.Progress) (*ArtifactoryServicesManager, error) {
	var err error
	manager := &ArtifactoryServicesManager{config: config, progress: progress}
	manager.client, err = rthttpclient.ArtifactoryClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetCommonDetails(artDetails).
		Build()
	if err != nil {
		return nil, err
	}
	return manager, err
}

func (sm *ArtifactoryServicesManager) CreateLocalRepository() *services.LocalRepositoryService {
	repositoryService := services.NewLocalRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetCommonDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManager) CreateRemoteRepository() *services.RemoteRepositoryService {
	repositoryService := services.NewRemoteRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetCommonDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManager) CreateVirtualRepository() *services.VirtualRepositoryService {
	repositoryService := services.NewVirtualRepositoryService(sm.client, false)
	repositoryService.ArtDetails = sm.config.GetCommonDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManager) UpdateLocalRepository() *services.LocalRepositoryService {
	repositoryService := services.NewLocalRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetCommonDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManager) UpdateRemoteRepository() *services.RemoteRepositoryService {
	repositoryService := services.NewRemoteRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetCommonDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManager) UpdateVirtualRepository() *services.VirtualRepositoryService {
	repositoryService := services.NewVirtualRepositoryService(sm.client, true)
	repositoryService.ArtDetails = sm.config.GetCommonDetails()
	return repositoryService
}

func (sm *ArtifactoryServicesManager) DeleteRepository(repoKey string) error {
	deleteRepositoryService := services.NewDeleteRepositoryService(sm.client)
	deleteRepositoryService.ArtDetails = sm.config.GetCommonDetails()
	return deleteRepositoryService.Delete(repoKey)
}

func (sm *ArtifactoryServicesManager) PublishBuildInfo(build *buildinfo.BuildInfo) error {
	buildInfoService := services.NewBuildInfoService(sm.client)
	buildInfoService.DryRun = sm.config.IsDryRun()
	buildInfoService.ArtDetails = sm.config.GetCommonDetails()
	return buildInfoService.PublishBuildInfo(build)
}

func (sm *ArtifactoryServicesManager) DistributeBuild(params services.BuildDistributionParams) error {
	distributionService := services.NewDistributionService(sm.client)
	distributionService.DryRun = sm.config.IsDryRun()
	distributionService.ArtDetails = sm.config.GetCommonDetails()
	return distributionService.BuildDistribute(params)
}

func (sm *ArtifactoryServicesManager) PromoteBuild(params services.PromotionParams) error {
	promotionService := services.NewPromotionService(sm.client)
	promotionService.DryRun = sm.config.IsDryRun()
	promotionService.ArtDetails = sm.config.GetCommonDetails()
	return promotionService.BuildPromote(params)
}

func (sm *ArtifactoryServicesManager) DiscardBuilds(params services.DiscardBuildsParams) error {
	discardService := services.NewDiscardBuildsService(sm.client)
	discardService.ArtDetails = sm.config.GetCommonDetails()
	return discardService.DiscardBuilds(params)
}

func (sm *ArtifactoryServicesManager) XrayScanBuild(params services.XrayScanParams) ([]byte, error) {
	xrayScanService := services.NewXrayScanService(sm.client)
	xrayScanService.ArtDetails = sm.config.GetCommonDetails()
	return xrayScanService.ScanBuild(params)
}

func (sm *ArtifactoryServicesManager) GetPathsToDelete(params services.DeleteParams) ([]utils.ResultItem, error) {
	deleteService := services.NewDeleteService(sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.ArtDetails = sm.config.GetCommonDetails()
	return deleteService.GetPathsToDelete(params)
}

func (sm *ArtifactoryServicesManager) DeleteFiles(resultItems []utils.ResultItem) (int, error) {
	deleteService := services.NewDeleteService(sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.ArtDetails = sm.config.GetCommonDetails()
	deleteService.Threads = sm.config.GetThreads()
	return deleteService.DeleteFiles(resultItems)
}

func (sm *ArtifactoryServicesManager) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	readFileService := services.NewReadFileService(sm.client)
	readFileService.DryRun = sm.config.IsDryRun()
	readFileService.ArtDetails = sm.config.GetCommonDetails()
	return readFileService.ReadRemoteFile(readPath)
}

func (sm *ArtifactoryServicesManager) DownloadFiles(params ...services.DownloadParams) ([]utils.FileInfo, int, error) {
	downloadService := services.NewDownloadService(sm.client)
	downloadService.DryRun = sm.config.IsDryRun()
	downloadService.ArtDetails = sm.config.GetCommonDetails()
	downloadService.Threads = sm.config.GetThreads()
	downloadService.Progress = sm.progress
	return downloadService.DownloadFiles(params...)
}

func (sm *ArtifactoryServicesManager) GetUnreferencedGitLfsFiles(params services.GitLfsCleanParams) ([]utils.ResultItem, error) {
	gitLfsCleanService := services.NewGitLfsCleanService(sm.client)
	gitLfsCleanService.DryRun = sm.config.IsDryRun()
	gitLfsCleanService.ArtDetails = sm.config.GetCommonDetails()
	return gitLfsCleanService.GetUnreferencedGitLfsFiles(params)
}

func (sm *ArtifactoryServicesManager) SearchFiles(params services.SearchParams) ([]utils.ResultItem, error) {
	searchService := services.NewSearchService(sm.client)
	searchService.ArtDetails = sm.config.GetCommonDetails()
	return searchService.Search(params)
}

func (sm *ArtifactoryServicesManager) Aql(aql string) ([]byte, error) {
	aqlService := services.NewAqlService(sm.client)
	aqlService.ArtDetails = sm.config.GetCommonDetails()
	return aqlService.ExecAql(aql)
}

func (sm *ArtifactoryServicesManager) SetProps(params services.PropsParams) (int, error) {
	setPropsService := services.NewPropsService(sm.client)
	setPropsService.ArtDetails = sm.config.GetCommonDetails()
	setPropsService.Threads = sm.config.GetThreads()
	return setPropsService.SetProps(params)
}

func (sm *ArtifactoryServicesManager) DeleteProps(params services.PropsParams) (int, error) {
	setPropsService := services.NewPropsService(sm.client)
	setPropsService.ArtDetails = sm.config.GetCommonDetails()
	setPropsService.Threads = sm.config.GetThreads()
	return setPropsService.DeleteProps(params)
}

func (sm *ArtifactoryServicesManager) UploadFiles(params ...services.UploadParams) (artifactsFileInfo []utils.FileInfo, totalUploaded, totalFailed int, err error) {
	uploadService := services.NewUploadService(sm.client)
	uploadService.Threads = sm.config.GetThreads()
	uploadService.ArtDetails = sm.config.GetCommonDetails()
	uploadService.DryRun = sm.config.IsDryRun()
	uploadService.Progress = sm.progress
	return uploadService.UploadFiles(params...)
}

func (sm *ArtifactoryServicesManager) Copy(params services.MoveCopyParams) (successCount, failedCount int, err error) {
	copyService := services.NewMoveCopyService(sm.client, services.COPY)
	copyService.DryRun = sm.config.IsDryRun()
	copyService.ArtDetails = sm.config.GetCommonDetails()
	return copyService.MoveCopyServiceMoveFilesWrapper(params)
}

func (sm *ArtifactoryServicesManager) Move(params services.MoveCopyParams) (successCount, failedCount int, err error) {
	moveService := services.NewMoveCopyService(sm.client, services.MOVE)
	moveService.DryRun = sm.config.IsDryRun()
	moveService.ArtDetails = sm.config.GetCommonDetails()
	return moveService.MoveCopyServiceMoveFilesWrapper(params)
}

func (sm *ArtifactoryServicesManager) PublishGoProject(params _go.GoParams) error {
	goService := _go.NewGoService(sm.client)
	goService.ArtDetails = sm.config.GetCommonDetails()
	return goService.PublishPackage(params)
}

func (sm *ArtifactoryServicesManager) Ping() ([]byte, error) {
	pingService := services.NewPingService(sm.client)
	pingService.ArtDetails = sm.config.GetCommonDetails()
	return pingService.Ping()
}

func (sm *ArtifactoryServicesManager) GetConfig() config.Config {
	return sm.config
}

func (sm *ArtifactoryServicesManager) GetBuildInfo(params services.BuildInfoParams) (*buildinfo.BuildInfo, error) {
	buildInfoService := services.NewBuildInfoService(sm.client)
	buildInfoService.ArtDetails = sm.config.GetCommonDetails()
	return buildInfoService.GetBuildInfo(params)
}

func (sm *ArtifactoryServicesManager) CreateToken(params services.CreateTokenParams) (services.CreateTokenResponseData, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetCommonDetails()
	return securityService.CreateToken(params)
}

func (sm *ArtifactoryServicesManager) GetTokens() (services.GetTokensResponseData, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetCommonDetails()
	return securityService.GetTokens()
}

func (sm *ArtifactoryServicesManager) RefreshToken(params services.RefreshTokenParams) (services.CreateTokenResponseData, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetCommonDetails()
	return securityService.RefreshToken(params)
}

func (sm *ArtifactoryServicesManager) RevokeToken(params services.RevokeTokenParams) (string, error) {
	securityService := services.NewSecurityService(sm.client)
	securityService.ArtDetails = sm.config.GetCommonDetails()
	return securityService.RevokeToken(params)
}

func (sm *ArtifactoryServicesManager) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}

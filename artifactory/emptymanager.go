package artifactory

import (
	"io"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	_go "github.com/jfrog/jfrog-client-go/artifactory/services/go"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

type ArtifactoryServicesManager interface {
	CreateLocalRepository() *services.LocalRepositoryService
	CreateRemoteRepository() *services.RemoteRepositoryService
	CreateVirtualRepository() *services.VirtualRepositoryService
	UpdateLocalRepository() *services.LocalRepositoryService
	UpdateRemoteRepository() *services.RemoteRepositoryService
	UpdateVirtualRepository() *services.VirtualRepositoryService
	DeleteRepository(repoKey string) error
	GetRepository(repoKey string) (*services.RepositoryDetails, error)
	CreatePermissionTarget(params services.PermissionTargetParams) error
	UpdatePermissionTarget(params services.PermissionTargetParams) error
	DeletePermissionTarget(permissionTargetName string) error
	PublishBuildInfo(build *buildinfo.BuildInfo) error
	DistributeBuild(params services.BuildDistributionParams) error
	PromoteBuild(params services.PromotionParams) error
	DiscardBuilds(params services.DiscardBuildsParams) error
	XrayScanBuild(params services.XrayScanParams) ([]byte, error)
	GetPathsToDelete(params services.DeleteParams) (*content.ContentReader, error)
	DeleteFiles(reader *content.ContentReader) (int, error)
	ReadRemoteFile(readPath string) (io.ReadCloser, error)
	initDownloadService() *services.DownloadService
	DownloadFiles(params ...services.DownloadParams) (totalDownloaded, totalExpected int, err error)
	DownloadFilesWithResultReader(params ...services.DownloadParams) (resultReader *content.ContentReader, totalDownloaded, totalExpected int, err error)
	GetUnreferencedGitLfsFiles(params services.GitLfsCleanParams) (*content.ContentReader, error)
	SearchFiles(params services.SearchParams) (*content.ContentReader, error)
	Aql(aql string) (io.ReadCloser, error)
	SetProps(params services.PropsParams) (int, error)
	DeleteProps(params services.PropsParams) (int, error)
	UploadFiles(params ...services.UploadParams) (artifactsFileInfo []utils.FileInfo, totalUploaded, totalFailed int, err error)
	Copy(params services.MoveCopyParams) (successCount, failedCount int, err error)
	Move(params services.MoveCopyParams) (successCount, failedCount int, err error)
	PublishGoProject(params _go.GoParams) error
	Ping() ([]byte, error)
	GetConfig() config.Config
	GetBuildInfo(params services.BuildInfoParams) (*buildinfo.BuildInfo, error)
	CreateToken(params services.CreateTokenParams) (services.CreateTokenResponseData, error)
	GetTokens() (services.GetTokensResponseData, error)
	RefreshToken(params services.RefreshTokenParams) (services.CreateTokenResponseData, error)
	RevokeToken(params services.RevokeTokenParams) (string, error)
	CreateReplication(params services.CreateReplicationParams) error
	UpdateReplication(params services.UpdateReplicationParams) error
	DeleteReplication(repoKey string) error
	GetReplication(repoKey string) ([]utils.ReplicationParams, error)
	GetVersion() (string, error)
	GetServiceId() (string, error)
	PromoteDocker(params services.DockerPromoteParams) error
	Client() *rthttpclient.ArtifactoryHttpClient
}

type EmptyArtifactoryServicesManager struct {
}

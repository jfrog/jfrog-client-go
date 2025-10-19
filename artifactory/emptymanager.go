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
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

type ArtifactoryServicesManager interface {
	CreateUpdateRepositoriesInBatch(body []byte, isUpdate bool) error
	CreateLocalRepository() *services.LocalRepositoryService
	CreateLocalRepositoryWithParams(params services.LocalRepositoryBaseParams) error
	CreateRemoteRepository() *services.RemoteRepositoryService
	CreateRemoteRepositoryWithParams(params services.RemoteRepositoryBaseParams) error
	CreateVirtualRepository() *services.VirtualRepositoryService
	CreateVirtualRepositoryWithParams(params services.VirtualRepositoryBaseParams) error
	CreateFederatedRepository() *services.FederatedRepositoryService
	CreateFederatedRepositoryWithParams(params services.FederatedRepositoryBaseParams) error
	CreateRepositoryWithParams(params interface{}, repoName string) error
	UpdateLocalRepository() *services.LocalRepositoryService
	UpdateRemoteRepository() *services.RemoteRepositoryService
	UpdateVirtualRepository() *services.VirtualRepositoryService
	UpdateFederatedRepository() *services.FederatedRepositoryService
	UpdateRepositoryWithParams(params interface{}, repoName string) error
	DeleteRepository(repoKey string) error
	GetRepository(repoKey string, repoDetails interface{}) error
	GetAllRepositories() (*[]services.RepositoryDetails, error)
	GetAllRepositoriesFiltered(params services.RepositoriesFilterParams) (*[]services.RepositoryDetails, error)
	IsRepoExists(repoKey string) (bool, error)
	CreatePermissionTarget(params services.PermissionTargetParams) error
	UpdatePermissionTarget(params services.PermissionTargetParams) error
	DeletePermissionTarget(permissionTargetName string) error
	GetPermissionTarget(permissionTargetName string) (*services.PermissionTargetParams, error)
	GetAllPermissionTargets() (*[]services.PermissionTargetParams, error)
	PublishBuildInfo(build *buildinfo.BuildInfo, projectKey string) (*clientutils.Sha256Summary, error)
	DeleteBuildInfo(build *buildinfo.BuildInfo, projectKey string, buildNumberFrequency int) error
	DistributeBuild(params services.BuildDistributionParams) error
	PromoteBuild(params services.PromotionParams) error
	DiscardBuilds(params services.DiscardBuildsParams) error
	XrayScanBuild(params services.XrayScanParams) ([]byte, error)
	GetPathsToDelete(params services.DeleteParams) (*content.ContentReader, error)
	DeleteFiles(reader *content.ContentReader) (int, error)
	ReadRemoteFile(readPath string) (io.ReadCloser, error)
	DownloadFiles(params ...services.DownloadParams) (totalDownloaded, totalFailed int, err error)
	DownloadFilesWithSummary(params ...services.DownloadParams) (operationSummary *utils.OperationSummary, err error)
	GetUnreferencedGitLfsFiles(params services.GitLfsCleanParams) (*content.ContentReader, error)
	SearchFiles(params services.SearchParams) (*content.ContentReader, error)
	Aql(aql string) (io.ReadCloser, error)
	SetProps(params services.PropsParams) (int, error)
	DeleteProps(params services.PropsParams) (int, error)
	GetItemProps(relativePath string) (*utils.ItemProperties, error)
	UploadFiles(uploadServiceOptions UploadServiceOptions, params ...services.UploadParams) (totalUploaded, totalFailed int, err error)
	UploadFilesWithSummary(uploadServiceOptions UploadServiceOptions, params ...services.UploadParams) (operationSummary *utils.OperationSummary, err error)
	Copy(params ...services.MoveCopyParams) (successCount, failedCount int, err error)
	Move(params ...services.MoveCopyParams) (successCount, failedCount int, err error)
	PublishGoProject(params _go.GoParams) (*utils.OperationSummary, error)
	Ping() ([]byte, error)
	GetConfig() config.Config
	GetBuildInfo(params services.BuildInfoParams) (*buildinfo.PublishedBuildInfo, bool, error)
	GetBuildRuns(params services.BuildInfoParams) (*buildinfo.BuildRuns, bool, error)
	CreateAPIKey() (string, error)
	RegenerateAPIKey() (string, error)
	GetAPIKey() (string, error)
	CreateToken(params services.CreateTokenParams) (auth.CreateTokenResponseData, error)
	GetTokens() (services.GetTokensResponseData, error)
	GetUserTokens(username string) ([]string, error)
	RefreshToken(params services.ArtifactoryRefreshTokenParams) (auth.CreateTokenResponseData, error)
	RevokeToken(params services.RevokeTokenParams) (string, error)
	CreateReplication(params services.CreateReplicationParams) error
	UpdateReplication(params services.UpdateReplicationParams) error
	DeleteReplication(repoKey string) error
	GetReplication(repoKey string) ([]utils.ReplicationParams, error)
	GetVersion() (string, error)
	GetRunningNodes() ([]string, error)
	GetServiceId() (string, error)
	GetConfigDescriptor() (string, error)
	ActivateKeyEncryption() error
	DeactivateKeyEncryption() (bool, error)
	PromoteDocker(params services.DockerPromoteParams) error
	Client() *jfroghttpclient.JfrogHttpClient
	GetGroup(params services.GroupParams) (*services.Group, error)
	GetAllGroups() (*[]string, error)
	CreateGroup(params services.GroupParams) error
	UpdateGroup(params services.GroupParams) error
	DeleteGroup(name string) error
	GetUser(params services.UserParams) (*services.User, error)
	GetAllUsers() ([]*services.User, error)
	CreateUser(params services.UserParams) error
	UpdateUser(params services.UserParams) error
	DeleteUser(name string) error
	GetLockedUsers() ([]string, error)
	UnlockUser(name string) error
	ConvertLocalToFederatedRepository(repoKey string) error
	TriggerFederatedRepositoryFullSyncAll(repoKey string) error
	TriggerFederatedRepositoryFullSyncMirror(repoKey string, mirrorUrl string) error
	Export(params services.ExportParams) error
	FolderInfo(relativePath string) (*utils.FolderInfo, error)
	FileInfo(relativePath string) (*utils.FileInfo, error)
	FileList(relativePath string, optionalParams utils.FileListParams) (*utils.FileListResponse, error)
	GetStorageInfo() (*utils.StorageInfo, error)
	CalculateStorageInfo() error
	ImportReleaseBundle(string) error
	GetPackageLeadFile(leadFileParams services.LeadFileParams) ([]byte, error)
	UploadTrustedKey(params services.TrustedKeyParams) (*services.TrustedKeyResponse, error)
}

// By using this struct, you have the option of overriding only some of the ArtifactoryServicesManager
// interface's methods, but still implement this interface.
// This comes in very handy for tests.
type EmptyArtifactoryServicesManager struct {
}

func (esm *EmptyArtifactoryServicesManager) CreateUpdateRepositoriesInBatch(_ []byte, _ bool) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateLocalRepository() *services.LocalRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateLocalRepositoryWithParams(services.LocalRepositoryBaseParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateRemoteRepository() *services.RemoteRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateRemoteRepositoryWithParams(services.RemoteRepositoryBaseParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateVirtualRepository() *services.VirtualRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateVirtualRepositoryWithParams(services.VirtualRepositoryBaseParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateFederatedRepository() *services.FederatedRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateFederatedRepositoryWithParams(services.FederatedRepositoryBaseParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateRepositoryWithParams(interface{}, string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateLocalRepository() *services.LocalRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateRemoteRepository() *services.RemoteRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateVirtualRepository() *services.VirtualRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateFederatedRepository() *services.FederatedRepositoryService {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateRepositoryWithParams(interface{}, string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteRepository(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetRepository(string, interface{}) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) IsRepoExists(string) (bool, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreatePermissionTarget(services.PermissionTargetParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdatePermissionTarget(services.PermissionTargetParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeletePermissionTarget(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetPermissionTarget(string) (*services.PermissionTargetParams, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetAllPermissionTargets() (*[]services.PermissionTargetParams, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) PublishBuildInfo(*buildinfo.BuildInfo, string) (*clientutils.Sha256Summary, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DistributeBuild(services.BuildDistributionParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) PromoteBuild(services.PromotionParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DiscardBuilds(services.DiscardBuildsParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) XrayScanBuild(services.XrayScanParams) ([]byte, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetPathsToDelete(services.DeleteParams) (*content.ContentReader, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteFiles(*content.ContentReader) (int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) ReadRemoteFile(string) (io.ReadCloser, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DownloadFiles(...services.DownloadParams) (int, int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DownloadFilesWithSummary(...services.DownloadParams) (*utils.OperationSummary, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetUnreferencedGitLfsFiles(services.GitLfsCleanParams) (*content.ContentReader, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) SearchFiles(services.SearchParams) (*content.ContentReader, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) Aql(string) (io.ReadCloser, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) SetProps(services.PropsParams) (int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteProps(services.PropsParams) (int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetItemProps(string) (*utils.ItemProperties, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UploadFiles(_ UploadServiceOptions, _ ...services.UploadParams) (int, int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UploadFilesWithSummary(_ UploadServiceOptions, _ ...services.UploadParams) (*utils.OperationSummary, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) Copy(...services.MoveCopyParams) (int, int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) Move(...services.MoveCopyParams) (int, int, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) PublishGoProject(_go.GoParams) (*utils.OperationSummary, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) Ping() ([]byte, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetConfig() config.Config {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetBuildInfo(services.BuildInfoParams) (*buildinfo.PublishedBuildInfo, bool, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetBuildRuns(services.BuildInfoParams) (*buildinfo.BuildRuns, bool, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateAPIKey() (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) RegenerateAPIKey() (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetAPIKey() (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateToken(services.CreateTokenParams) (auth.CreateTokenResponseData, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetTokens() (services.GetTokensResponseData, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetUserTokens(string) ([]string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) RefreshToken(services.ArtifactoryRefreshTokenParams) (auth.CreateTokenResponseData, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) RevokeToken(services.RevokeTokenParams) (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateReplication(services.CreateReplicationParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateReplication(services.UpdateReplicationParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteReplication(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetReplication(string) ([]utils.ReplicationParams, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetRunningNodes() ([]string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetServiceId() (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetConfigDescriptor() (string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) ActivateKeyEncryption() error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeactivateKeyEncryption() (bool, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) PromoteDocker(services.DockerPromoteParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetAllRepositories() (*[]services.RepositoryDetails, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetAllRepositoriesFiltered(services.RepositoriesFilterParams) (*[]services.RepositoryDetails, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetUser(services.UserParams) (*services.User, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetAllUsers() ([]*services.User, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateUser(services.UserParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateUser(services.UserParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteUser(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetLockedUsers() ([]string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UnlockUser(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetGroup(services.GroupParams) (*services.Group, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetAllGroups() (*[]string, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CreateGroup(services.GroupParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UpdateGroup(services.GroupParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteGroup(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) ConvertLocalToFederatedRepository(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) TriggerFederatedRepositoryFullSyncAll(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) TriggerFederatedRepositoryFullSyncMirror(string, string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) Export(services.ExportParams) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) FolderInfo(string) (*utils.FolderInfo, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) FileInfo(string) (*utils.FileInfo, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) FileList(string, utils.FileListParams) (*utils.FileListResponse, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetStorageInfo() (*utils.StorageInfo, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) CalculateStorageInfo() error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) ImportReleaseBundle(string) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetPackageLeadFile(services.LeadFileParams) ([]byte, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) UploadTrustedKey(services.TrustedKeyParams) (*services.TrustedKeyResponse, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) DeleteBuildInfo(*buildinfo.BuildInfo, string, int) error {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetRepositoriesStats(string) ([]byte, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetJPDsStats(string) ([]byte, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetReleaseBundlesStats(string) ([]byte, error) {
	panic("Failed: Method is not implemented")
}

func (esm *EmptyArtifactoryServicesManager) GetTokenDetails(string, string) ([]byte, error) {
	panic("Failed: Method is not implemented")
}

// Compile time check of interface implementation.
// Since EmptyArtifactoryServicesManager can be used by tests external to this project, we want this project's tests to fail,
// if EmptyArtifactoryServicesManager stops implementing the ArtifactoryServicesManager interface.
var _ ArtifactoryServicesManager = (*EmptyArtifactoryServicesManager)(nil)

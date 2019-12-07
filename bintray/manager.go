package bintray

import (
	"github.com/jfrog/jfrog-client-go/bintray/services"
	"github.com/jfrog/jfrog-client-go/bintray/services/accesskeys"
	"github.com/jfrog/jfrog-client-go/bintray/services/entitlements"
	"github.com/jfrog/jfrog-client-go/bintray/services/gpg"
	"github.com/jfrog/jfrog-client-go/bintray/services/logs"
	"github.com/jfrog/jfrog-client-go/bintray/services/packages"
	"github.com/jfrog/jfrog-client-go/bintray/services/repositories"
	"github.com/jfrog/jfrog-client-go/bintray/services/url"
	"github.com/jfrog/jfrog-client-go/bintray/services/utils"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"github.com/jfrog/jfrog-client-go/httpclient"
)

type ServicesManager struct {
	client *httpclient.HttpClient
	config Config
}

func New(config Config) (*ServicesManager, error) {
	var err error
	manager := &ServicesManager{config: config}
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil, err
	}
	manager.client = client
	return manager, err
}

func (sm *ServicesManager) newDownloadService() *services.DownloadService {
	downloadService := services.NewDownloadService(sm.client)
	downloadService.BintrayDetails = sm.config.GetBintrayDetails()
	downloadService.Threads = sm.config.GetThreads()
	return downloadService
}

func (sm *ServicesManager) DownloadFile(params *services.DownloadFileParams) (totalUploaded, totalFailed int, err error) {
	downloadService := sm.newDownloadService()
	return downloadService.DownloadFile(params)
}

func (sm *ServicesManager) DownloadVersion(params *services.DownloadVersionParams) (totalDownloaded, totalFailed int, err error) {
	downloadService := sm.newDownloadService()
	return downloadService.DownloadVersion(params)
}

func (sm *ServicesManager) UploadFiles(params *services.UploadParams) (totalUploaded, totalFailed int, err error) {
	uploadService := services.NewUploadService(sm.client)
	uploadService.BintrayDetails = sm.config.GetBintrayDetails()
	uploadService.DryRun = sm.config.IsDryRun()
	uploadService.Threads = sm.config.GetThreads()
	return uploadService.Upload(params)
}

func (sm *ServicesManager) newVersionService() *versions.VersionService {
	versionService := versions.NewService(sm.client)
	versionService.BintrayDetails = sm.config.GetBintrayDetails()
	return versionService
}

func (sm *ServicesManager) CreateVersion(params *versions.Params) error {
	return sm.newVersionService().Create(params)
}

func (sm *ServicesManager) UpdateVersion(params *versions.Params) error {
	return sm.newVersionService().Update(params)
}

func (sm *ServicesManager) PublishVersion(path *versions.Path) error {
	return sm.newVersionService().Publish(path)
}

func (sm *ServicesManager) DeleteVersion(path *versions.Path) error {
	return sm.newVersionService().Delete(path)
}

func (sm *ServicesManager) ShowVersion(path *versions.Path) error {
	return sm.newVersionService().Show(path)
}

func (sm *ServicesManager) IsVersionExists(path *versions.Path) (bool, error) {
	return sm.newVersionService().IsVersionExists(path)
}

func (sm *ServicesManager) CalcMetadata(path *versions.Path) (bool, error) {
	return sm.newVersionService().CalcMetadata(path)
}

func (sm *ServicesManager) newPackageService() *packages.PackageService {
	packageService := packages.NewService(sm.client)
	packageService.BintrayDetails = sm.config.GetBintrayDetails()
	return packageService
}

func (sm *ServicesManager) CreatePackage(params *packages.Params) error {
	return sm.newPackageService().Create(params)
}

func (sm *ServicesManager) UpdatePackage(params *packages.Params) error {
	return sm.newPackageService().Update(params)
}

func (sm *ServicesManager) DeletePackage(path *packages.Path) error {
	return sm.newPackageService().Delete(path)
}

func (sm *ServicesManager) ShowPackage(path *packages.Path) error {
	return sm.newPackageService().Show(path)
}

func (sm *ServicesManager) IsPackageExists(path *packages.Path) (bool, error) {
	return sm.newPackageService().IsPackageExists(path)
}

func (sm *ServicesManager) IsRepoExists(path *repositories.Path) (bool, error) {
	repositoryService := repositories.NewService(sm.client)
	repositoryService.BintrayDetails = sm.config.GetBintrayDetails()
	return repositoryService.IsRepoExists(path)
}

func (sm *ServicesManager) CreateReposIfNeeded(path *repositories.Path, params *repositories.Config, repoConfig string) (bool, error) {
	repositoryService := repositories.NewService(sm.client)
	repositoryService.BintrayDetails = sm.config.GetBintrayDetails()
	return repositoryService.CreateReposIfNeeded(path, params, repoConfig)
}

func (sm *ServicesManager) ExecCreateRepoRest(path *repositories.Path, params *repositories.Config, repoConfig string) (bool, error) {
	repositoryService := repositories.NewService(sm.client)
	repositoryService.BintrayDetails = sm.config.GetBintrayDetails()
	return repositoryService.ExecCreateRepoRest(path, params, repoConfig)
}

func (sm *ServicesManager) ExecDeleteRepoRest(path *repositories.Path) error {
	repositoryService := repositories.NewService(sm.client)
	repositoryService.BintrayDetails = sm.config.GetBintrayDetails()
	return repositoryService.ExecDeleteRepoRest(path)
}

func (sm *ServicesManager) newAccessKeysService() *accesskeys.AccessKeysService {
	accessKeysService := accesskeys.NewService(sm.client)
	accessKeysService.BintrayDetails = sm.config.GetBintrayDetails()
	return accessKeysService
}

func (sm *ServicesManager) CreateAccessKey(params *accesskeys.Params) error {
	return sm.newAccessKeysService().Create(params)
}

func (sm *ServicesManager) UpdateAccessKey(params *accesskeys.Params) error {
	return sm.newAccessKeysService().Update(params)
}

func (sm *ServicesManager) ShowAllAccessKeys(org string) error {
	return sm.newAccessKeysService().ShowAll(org)
}

func (sm *ServicesManager) ShowAccessKey(org, id string) error {
	return sm.newAccessKeysService().Show(org, id)
}

func (sm *ServicesManager) DeleteAccessKey(org, id string) error {
	return sm.newAccessKeysService().Delete(org, id)
}

func (sm *ServicesManager) newGpgService() *gpg.GpgService {
	gpgService := gpg.NewService(sm.client)
	gpgService.BintrayDetails = sm.config.GetBintrayDetails()
	return gpgService
}

func (sm *ServicesManager) GpgSignFile(pathDetails *utils.PathDetails, passphrase string) error {
	return sm.newGpgService().SignFile(pathDetails, passphrase)
}

func (sm *ServicesManager) GpgSignVersion(versionPath *versions.Path, passphrase string) error {
	return sm.newGpgService().SignVersion(versionPath, passphrase)
}

func (sm *ServicesManager) newLogsService() *logs.LogsService {
	logsService := logs.NewService(sm.client)
	logsService.BintrayDetails = sm.config.GetBintrayDetails()
	return logsService
}

func (sm *ServicesManager) LogsList(versionPath *versions.Path) error {
	return sm.newLogsService().List(versionPath)
}

func (sm *ServicesManager) DownloadLog(versionPath *versions.Path, logName string) error {
	return sm.newLogsService().Download(versionPath, logName)
}

func (sm *ServicesManager) SignUrl(params *url.Params) error {
	signUrlService := url.NewService(sm.client)
	signUrlService.BintrayDetails = sm.config.GetBintrayDetails()
	return signUrlService.SignVersion(params)
}

func (sm *ServicesManager) newEntitlementService() *entitlements.EntitlementsService {
	entitlementsService := entitlements.NewService(sm.client)
	entitlementsService.BintrayDetails = sm.config.GetBintrayDetails()
	return entitlementsService
}

func (sm *ServicesManager) ShowAllEntitlements(versionPath *versions.Path) error {
	return sm.newEntitlementService().ShowAll(versionPath)
}

func (sm *ServicesManager) ShowEntitlement(id string, path *versions.Path) error {
	return sm.newEntitlementService().Show(id, path)
}

func (sm *ServicesManager) CreateEntitlement(params *entitlements.Params) error {
	return sm.newEntitlementService().Create(params)
}

func (sm *ServicesManager) UpdateEntitlement(params *entitlements.Params) error {
	return sm.newEntitlementService().Update(params)
}

func (sm *ServicesManager) DeleteEntitlement(id string, path *versions.Path) error {
	return sm.newEntitlementService().Delete(id, path)
}

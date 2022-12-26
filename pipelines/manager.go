package pipelines

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/pipelines/services"
)

type PipelinesServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func New(config config.Config) (*PipelinesServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &PipelinesServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetContext(config.GetContext()).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()
	return manager, err
}

func (sm *PipelinesServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *PipelinesServicesManager) GetSystemInfo() (*services.PipelinesSystemInfo, error) {
	systemService := services.NewSystemService(sm.client)
	systemService.ServiceDetails = sm.config.GetServiceDetails()
	return systemService.GetSystemInfo()
}

func (sm *PipelinesServicesManager) CreateGithubIntegration(integrationName, token string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateGithubIntegration(integrationName, token)
}

func (sm *PipelinesServicesManager) CreateGithubEnterpriseIntegration(integrationName, url, token string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateGithubEnterpriseIntegration(integrationName, url, token)
}

func (sm *PipelinesServicesManager) CreateBitbucketIntegration(integrationName, username, token string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateBitbucketIntegration(integrationName, username, token)
}

func (sm *PipelinesServicesManager) CreateBitbucketServerIntegration(integrationName, url, username, passwordOrToken string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateBitbucketServerIntegration(integrationName, url, username, passwordOrToken)
}

func (sm *PipelinesServicesManager) CreateGitlabIntegration(integrationName, url, token string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateGitlabIntegration(integrationName, url, token)
}

func (sm *PipelinesServicesManager) CreateArtifactoryIntegration(integrationName, url, user, apikey string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateArtifactoryIntegration(integrationName, url, user, apikey)
}

func (sm *PipelinesServicesManager) GetIntegrationById(integrationId int) (*services.Integration, error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.GetIntegrationById(integrationId)
}

func (sm *PipelinesServicesManager) GetIntegrationByName(integrationName string) (*services.Integration, error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.GetIntegrationByName(integrationName)
}

func (sm *PipelinesServicesManager) GetAllIntegrations() ([]services.Integration, error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.GetAllIntegrations()
}

func (sm *PipelinesServicesManager) DeleteIntegration(integrationId int) error {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.DeleteIntegration(integrationId)
}

func (sm *PipelinesServicesManager) AddPipelineSource(projectIntegrationId int, repositoryFullName, branch, fileFilter string) (id int, err error) {
	sourcesService := services.NewSourcesService(sm.client)
	sourcesService.ServiceDetails = sm.config.GetServiceDetails()
	return sourcesService.AddSource(projectIntegrationId, repositoryFullName, branch, fileFilter)
}

func (sm *PipelinesServicesManager) GetPipelineRunStatusByBranch(branch, pipeline string, isMultiBranch bool) (*services.PipelineRunStatusResponse, error) {
	runService := services.NewRunService(sm.client)
	runService.ServiceDetails = sm.config.GetServiceDetails()
	return runService.GetRunStatus(branch, pipeline, isMultiBranch)
}

func (sm *PipelinesServicesManager) TriggerPipelineRun(branch, pipeline string, isMultiBranch bool) error {
	runService := services.NewRunService(sm.client)
	runService.ServiceDetails = sm.config.GetServiceDetails()
	return runService.TriggerPipelineRun(branch, pipeline, isMultiBranch)
}

func (sm *PipelinesServicesManager) SyncPipelineResource(branch, repoFullName string) error {
	syncService := services.NewSyncService(sm.client)
	syncService.ServiceDetails = sm.config.GetServiceDetails()
	return syncService.SyncPipelineSource(branch, repoFullName)
}

func (sm *PipelinesServicesManager) GetSyncStatusForPipelineResource(repo, branch string) ([]services.PipelineSyncStatus, error) {
	syncStatusService := services.NewSyncStatusService(sm.client)
	syncStatusService.ServiceDetails = sm.config.GetServiceDetails()
	pipSyncStatus, err := syncStatusService.GetSyncPipelineResourceStatus(repo, branch)
	return pipSyncStatus, err
}

func (sm *PipelinesServicesManager) CancelRun(runID int) error {
	runService := services.NewRunService(sm.client)
	runService.ServiceDetails = sm.config.GetServiceDetails()
	return runService.CancelRun(runID)
}

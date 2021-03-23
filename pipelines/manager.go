package pipelines

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/pipelines/services"
)

type PipelinesServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func New(details *auth.ServiceDetails, config config.Config) (*PipelinesServicesManager, error) {
	err := (*details).InitSsh()
	if err != nil {
		return nil, err
	}
	manager := &PipelinesServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetServiceDetails(details).
		SetContext(config.GetContext()).
		Build()
	if err != nil {
		return nil, err
	}
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

func (sm *PipelinesServicesManager) CreateBitbucketIntegration(integrationName, username, token string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateBitbucketIntegration(integrationName, username, token)
}

func (sm *PipelinesServicesManager) CreateGitlabIntegration(integrationName, username, token string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateGitlabIntegration(integrationName, username, token)
}

func (sm *PipelinesServicesManager) CreateArtifactoryIntegration(integrationName, url, user, apikey string) (id int, err error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.CreateArtifactoryIntegration(integrationName, url, user, apikey)
}

func (sm *PipelinesServicesManager) GetIntegration(integrationId int) (*services.Integration, error) {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.GetIntegration(integrationId)
}

func (sm *PipelinesServicesManager) DeleteIntegration(integrationId int) error {
	integrationsService := services.NewIntegrationsService(sm.client)
	integrationsService.ServiceDetails = sm.config.GetServiceDetails()
	return integrationsService.DeleteIntegration(integrationId)
}

func (sm *PipelinesServicesManager) AddPipelineSource(projectIntegrationId int, repositoryFullName, branch, fileFilter string) (id int, err error) {
	sourcesService := services.NewSourcesService(sm.client)
	sourcesService.ServiceDetails = sm.config.GetServiceDetails()
	return sourcesService.AddPipelineSource(projectIntegrationId, repositoryFullName, branch, fileFilter)
}

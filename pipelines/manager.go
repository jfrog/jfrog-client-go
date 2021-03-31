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

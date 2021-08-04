package access

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

type AccessServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func New(config config.Config) (*AccessServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &AccessServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetContext(config.GetContext()).
		SetRetries(config.GetHttpRetries()).
		Build()

	return manager, err
}

func (sm *AccessServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *AccessServicesManager) CreateProject(params services.ProjectParams) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.CreateProject(params)
}

func (sm *AccessServicesManager) UpdateProject(params services.ProjectParams) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.UpdateProject(params)
}

func (sm *AccessServicesManager) DeleteProject(projectKey string) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.DeleteProject(projectKey)
}

func (sm *AccessServicesManager) AssignRepoToProject(repoName, projectKey string, isForce bool) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.AssignRepoToProject(repoName, projectKey, isForce)
}

func (sm *AccessServicesManager) UnassignRepoFromProject(repoName string) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.UnassignRepoFromProject(repoName)
}

package access

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/auth"
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
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()

	return manager, err
}

func (sm *AccessServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *AccessServicesManager) Ping() ([]byte, error) {
	projectService := services.NewPingService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.Ping()
}

func (sm *AccessServicesManager) CreateProject(params services.ProjectParams) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.Create(params)
}

func (sm *AccessServicesManager) GetProject(projectKey string) (*services.Project, error) {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.Get(projectKey)
}

func (sm *AccessServicesManager) GetAllProjects() ([]services.Project, error) {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.GetAll()
}

func (sm *AccessServicesManager) UpdateProject(params services.ProjectParams) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.Update(params)
}

func (sm *AccessServicesManager) DeleteProject(projectKey string) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.Delete(projectKey)
}

func (sm *AccessServicesManager) AssignRepoToProject(repoName, projectKey string, isForce bool) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.AssignRepo(repoName, projectKey, isForce)
}

func (sm *AccessServicesManager) UnassignRepoFromProject(repoName string) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.UnassignRepo(repoName)
}

func (sm *AccessServicesManager) GetProjectsGroups(projectKey string) (*[]services.ProjectGroup, error) {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.GetGroups(projectKey)
}

func (sm *AccessServicesManager) GetProjectsGroup(projectKey string, groupName string) (*services.ProjectGroup, error) {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.GetGroup(projectKey, groupName)
}

func (sm *AccessServicesManager) UpdateGroupInProject(projectKey string, groupName string, group services.ProjectGroup) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.UpdateGroup(projectKey, groupName, group)
}

func (sm *AccessServicesManager) DeleteExistingProjectGroup(projectKey string, groupName string) error {
	projectService := services.NewProjectService(sm.client)
	projectService.ServiceDetails = sm.config.GetServiceDetails()
	return projectService.DeleteExistingGroup(projectKey, groupName)
}

func (sm *AccessServicesManager) CreateAccessToken(params services.CreateTokenParams) (auth.CreateTokenResponseData, error) {
	tokenService := services.NewTokenService(sm.client)
	tokenService.ServiceDetails = sm.config.GetServiceDetails()
	return tokenService.CreateAccessToken(params)
}

func (sm *AccessServicesManager) RefreshAccessToken(params services.CreateTokenParams) (auth.CreateTokenResponseData, error) {
	tokenService := services.NewTokenService(sm.client)
	tokenService.ServiceDetails = sm.config.GetServiceDetails()
	return tokenService.RefreshAccessToken(params)
}

func (sm *AccessServicesManager) InviteUser(email, source string) error {
	inviteService := services.NewInviteService(sm.client)
	inviteService.ServiceDetails = sm.config.GetServiceDetails()
	return inviteService.InviteUser(email, source)
}

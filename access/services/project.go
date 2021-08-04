package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const projectsApi = "v1/projects"

type ProjectParams struct {
	ProjectDetails Project
}

func NewProjectParams() ProjectParams {
	return ProjectParams{}
}

type Project struct {
	DisplayName       string           `json:"display_name,omitempty"`
	Description       string           `json:"description,omitempty"`
	AdminPrivileges   *AdminPrivileges `json:"admin_privileges,omitempty"`
	SoftLimit         bool             `json:"soft_limit,omitempty"`
	StorageQuotaBytes float64          `json:"storage_quota_bytes,omitempty"`
	ProjectKey        string           `json:"project_key,omitempty"`
}

type AdminPrivileges struct {
	ManageMembers   bool `json:"manage_members,omitempty"`
	ManageResources bool `json:"manage_resources,omitempty"`
	IndexResources  bool `json:"index_resources,omitempty"`
}

type ProjectService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

func NewProjectService(client *jfroghttpclient.JfrogHttpClient) *ProjectService {
	return &ProjectService{client: client}
}

func (us *ProjectService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	us.ServiceDetails = rt
}

func (us *ProjectService) getProjectsBaseUrl() string {
	return fmt.Sprintf("%s%s", us.ServiceDetails.GetUrl(), projectsApi)
}

func (us *ProjectService) GetProject(projectKey string) (u *Project, err error) {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", us.getProjectsBaseUrl(), projectKey)
	resp, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// The case the requested user is not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return &project, nil
}

func (us *ProjectService) CreateProject(params ProjectParams) error {
	project, err := us.GetProject(params.ProjectDetails.ProjectKey)
	if err != nil {
		return err
	}
	if project != nil {
		return errorutils.CheckError(fmt.Errorf("project '%s' already exists", project.ProjectKey))
	}
	content, httpDetails, err := us.createOrUpdateProjectRequest(params.ProjectDetails)
	if err != nil {
		return err
	}
	resp, body, err := us.client.SendPost(us.getProjectsBaseUrl(), content, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

func (us *ProjectService) UpdateProject(params ProjectParams) error {
	content, httpDetails, err := us.createOrUpdateProjectRequest(params.ProjectDetails)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", us.getProjectsBaseUrl(), params.ProjectDetails.ProjectKey)
	resp, body, err := us.client.SendPut(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

func (us *ProjectService) createOrUpdateProjectRequest(project Project) (requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = us.ServiceDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(project)
	if errorutils.CheckError(err) != nil {
		return
	}
	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	return
}

func (us *ProjectService) DeleteProject(projectKey string) error {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", us.getProjectsBaseUrl(), projectKey)
	resp, body, err := us.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return errorutils.CheckError(fmt.Errorf("no response provided (including status code)"))
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusNoContent); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return err
}

func (us *ProjectService) AssignRepoToProject(repoName, projectKey string, isForce bool) error {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/_/attach/repositories/%s/%s?force=%b", us.getProjectsBaseUrl(), repoName, projectKey, isForce)
	resp, body, err := us.client.SendPut(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusNoContent); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}
func (us *ProjectService) UnassignRepoFromProject(repoName string) error {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/_/attach/repositories/%s", us.getProjectsBaseUrl(), repoName)
	resp, body, err := us.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusNoContent); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

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

const projectsApi = "api/v1/projects"

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
	SoftLimit         *bool            `json:"soft_limit,omitempty"`
	StorageQuotaBytes float64          `json:"storage_quota_bytes,omitempty"`
	ProjectKey        string           `json:"project_key,omitempty"`
}

type AdminPrivileges struct {
	ManageMembers   *bool `json:"manage_members,omitempty"`
	ManageResources *bool `json:"manage_resources,omitempty"`
	IndexResources  *bool `json:"index_resources,omitempty"`
}

type ProjectService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

func NewProjectService(client *jfroghttpclient.JfrogHttpClient) *ProjectService {
	return &ProjectService{client: client}
}

func (ps *ProjectService) getProjectsBaseUrl() string {
	return fmt.Sprintf("%s%s", ps.ServiceDetails.GetUrl(), projectsApi)
}

func (ps *ProjectService) Get(projectKey string) (u *Project, err error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", ps.getProjectsBaseUrl(), projectKey)
	resp, body, _, err := ps.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// In case the requested project is not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	var project Project
	err = json.Unmarshal(body, &project)
	return &project, errorutils.CheckError(err)
}

func (ps *ProjectService) Create(params ProjectParams) error {
	project, err := ps.Get(params.ProjectDetails.ProjectKey)
	if err != nil {
		return err
	}
	if project != nil {
		return errorutils.CheckErrorf("project '%s' already exists", project.ProjectKey)
	}
	content, httpDetails, err := ps.createOrUpdateRequest(params.ProjectDetails)
	if err != nil {
		return err
	}
	resp, body, err := ps.client.SendPost(ps.getProjectsBaseUrl(), content, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

func (ps *ProjectService) Update(params ProjectParams) error {
	content, httpDetails, err := ps.createOrUpdateRequest(params.ProjectDetails)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", ps.getProjectsBaseUrl(), params.ProjectDetails.ProjectKey)
	resp, body, err := ps.client.SendPut(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

func (ps *ProjectService) createOrUpdateRequest(project Project) (requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = ps.ServiceDetails.CreateHttpClientDetails()
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

func (ps *ProjectService) Delete(projectKey string) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", ps.getProjectsBaseUrl(), projectKey)
	resp, body, err := ps.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return errorutils.CheckErrorf("no response provided (including status code)")
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusNoContent); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return err
}

func (ps *ProjectService) AssignRepo(repoName, projectKey string, isForce bool) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/_/attach/repositories/%s/%s?force=%t", ps.getProjectsBaseUrl(), repoName, projectKey, isForce)
	resp, body, err := ps.client.SendPut(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusNoContent); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

func (ps *ProjectService) UnassignRepo(repoName string) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/_/attach/repositories/%s", ps.getProjectsBaseUrl(), repoName)
	resp, body, err := ps.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusNoContent); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

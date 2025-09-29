package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
)

const projectsApi = "api/v1/projects"

type ProjectParams struct {
	ProjectDetails Project
}

func NewProjectParams() ProjectParams {
	return ProjectParams{}
}

type Project struct {
	DisplayName       string           `json:"display_name,omitempty" log:"Name"`
	Description       string           `json:"description,omitempty"`
	AdminPrivileges   *AdminPrivileges `json:"admin_privileges,omitempty"`
	SoftLimit         *bool            `json:"soft_limit,omitempty"`
	StorageQuotaBytes float64          `json:"storage_quota_bytes,omitempty"`
	ProjectKey        string           `json:"project_key,omitempty" log:"ProjectKey"`
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

type ProjectGroup struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type ProjectGroups struct {
	Members []ProjectGroup `json:"members"`
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
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var project Project
	err = json.Unmarshal(body, &project)
	return &project, errorutils.CheckError(err)
}

func (ps *ProjectService) GetAll() ([]Project, error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := ps.getProjectsBaseUrl()
	resp, body, _, err := ps.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var projects []Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, errorutils.CheckErrorf("failed extracting projects list from payload: %s", err.Error())
	}
	return projects, nil
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
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated)
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
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated)
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
	if err != nil {
		return err
	}
	if resp == nil {
		return errorutils.CheckErrorf("no response provided (including status code)")
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}

func (ps *ProjectService) AssignRepo(repoName, projectKey string, isForce bool) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/_/attach/repositories/%s/%s?force=%t", ps.getProjectsBaseUrl(), repoName, projectKey, isForce)
	resp, body, err := ps.client.SendPut(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusNoContent)
}

func (ps *ProjectService) UnassignRepo(repoName string) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/_/attach/repositories/%s", ps.getProjectsBaseUrl(), repoName)
	resp, body, err := ps.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusNoContent)
}

func (ps *ProjectService) GetGroups(projectKey string) (*[]ProjectGroup, error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s/groups", ps.getProjectsBaseUrl(), projectKey)
	resp, body, _, err := ps.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// In case the requested project is not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var projectGroups ProjectGroups
	err = json.Unmarshal(body, &projectGroups)
	return &projectGroups.Members, errorutils.CheckError(err)
}

func (ps *ProjectService) GetGroup(projectKey string, groupName string) (*ProjectGroup, error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s/groups/%s", ps.getProjectsBaseUrl(), projectKey, groupName)
	resp, body, _, err := ps.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// In case the requested project or group in project is not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var projectGroup ProjectGroup
	err = json.Unmarshal(body, &projectGroup)
	return &projectGroup, errorutils.CheckError(err)
}

func (ps *ProjectService) UpdateGroup(projectKey string, groupName string, group ProjectGroup) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s/groups/%s", ps.getProjectsBaseUrl(), projectKey, groupName)
	requestContent, err := json.Marshal(group)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	resp, body, err := ps.client.SendPut(url, requestContent, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (ps *ProjectService) DeleteExistingGroup(projectKey string, groupName string) error {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s/groups/%s", ps.getProjectsBaseUrl(), projectKey, groupName)
	resp, body, err := ps.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}

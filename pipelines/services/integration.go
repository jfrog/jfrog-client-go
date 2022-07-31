package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type IntegrationsService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewIntegrationsService(client *jfroghttpclient.JfrogHttpClient) *IntegrationsService {
	return &IntegrationsService{client: client}
}

const (
	ArtifactoryName      = "artifactory"
	artifactoryId        = 98
	GithubName           = "github"
	githubId             = 20
	GithubEnterpriseName = "githubEnterprise"
	githubEnterpriseId   = 15
	githubDefaultUrl     = "https://api.github.com"
	BitbucketName        = "bitbucket"
	bitbucketId          = 16
	bitbucketDefaultUrl  = "https://api.bitbucket.org"
	BitbucketServerName  = "bitbucketServerBasic"
	bitbucketServerId    = 90
	GitlabName           = "gitlab"
	gitlabId             = 21

	defaultProjectId = 1

	urlLabel      = "url"
	tokenLabel    = "token"
	userLabel     = "user"
	usernameLabel = "username"
	apikeyLabel   = "apikey"
	passwordLabel = "password"

	integrationsRestApi = "api/v1/projectIntegrations/"
)

func (is *IntegrationsService) CreateGithubIntegration(integrationName, token string) (id int, err error) {
	integration := IntegrationCreation{
		Integration: Integration{
			Name:                  integrationName,
			MasterIntegrationId:   githubId,
			MasterIntegrationName: GithubName,
			ProjectId:             defaultProjectId,
		},
		FormJSONValues: []jsonValues{
			{urlLabel, githubDefaultUrl},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateGithubEnterpriseIntegration(integrationName, url, token string) (id int, err error) {
	integration := IntegrationCreation{
		Integration: Integration{
			Name:                  integrationName,
			MasterIntegrationId:   githubEnterpriseId,
			MasterIntegrationName: GithubEnterpriseName,
			ProjectId:             defaultProjectId,
		},
		FormJSONValues: []jsonValues{
			{urlLabel, url},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateBitbucketIntegration(integrationName, username, token string) (id int, err error) {
	integration := IntegrationCreation{
		Integration: Integration{
			Name:                  integrationName,
			MasterIntegrationId:   bitbucketId,
			MasterIntegrationName: BitbucketName,
			ProjectId:             defaultProjectId},
		FormJSONValues: []jsonValues{
			{urlLabel, bitbucketDefaultUrl},
			{usernameLabel, username},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateBitbucketServerIntegration(integrationName, url, username, passwordOrToken string) (id int, err error) {
	integration := IntegrationCreation{
		Integration: Integration{
			Name:                  integrationName,
			MasterIntegrationId:   bitbucketServerId,
			MasterIntegrationName: BitbucketServerName,
			ProjectId:             defaultProjectId},
		FormJSONValues: []jsonValues{
			{urlLabel, url},
			{usernameLabel, username},
			{passwordLabel, passwordOrToken},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateGitlabIntegration(integrationName, url, token string) (id int, err error) {
	integration := IntegrationCreation{
		Integration: Integration{
			Name:                  integrationName,
			MasterIntegrationId:   gitlabId,
			MasterIntegrationName: GitlabName,
			ProjectId:             defaultProjectId},
		FormJSONValues: []jsonValues{
			{urlLabel, url},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateArtifactoryIntegration(integrationName, url, user, apikey string) (id int, err error) {
	integration := IntegrationCreation{
		Integration: Integration{
			Name:                  integrationName,
			MasterIntegrationId:   artifactoryId,
			MasterIntegrationName: ArtifactoryName,
			ProjectId:             defaultProjectId},
		FormJSONValues: []jsonValues{
			{urlLabel, strings.TrimSuffix(url, "/")},
			{userLabel, user},
			{apikeyLabel, apikey},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) createIntegration(integration IntegrationCreation) (id int, err error) {
	log.Debug("Creating " + integration.MasterIntegrationName + " integration...")
	content, err := json.Marshal(integration)
	if err != nil {
		return -1, errorutils.CheckError(err)
	}
	httpDetails := is.ServiceDetails.CreateHttpClientDetails()
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	utils.MergeMaps(httpDetails.Headers, headers)
	httpDetails.Headers = headers

	resp, body, err := is.client.SendPost(is.ServiceDetails.GetUrl()+integrationsRestApi, content, &httpDetails)
	if err != nil {
		return -1, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		if resp.StatusCode == http.StatusConflict {
			return -1, &IntegrationAlreadyExistsError{InnerError: err}
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return -1, &IntegrationUnauthorizedError{InnerError: err}
		}
		return -1, err
	}

	created := &Integration{}
	err = json.Unmarshal(body, created)
	return created.Id, errorutils.CheckError(err)
}

type Integration struct {
	Name                  string   `json:"name,omitempty"`
	MasterIntegrationId   int      `json:"masterIntegrationId,omitempty"`
	MasterIntegrationName string   `json:"masterIntegrationName,omitempty"`
	ProjectId             int      `json:"projectId,omitempty"`
	Environments          []string `json:"environments,omitempty"`
	// Following fields returned when fetching or creating integration:
	Id int `json:"id,omitempty"`
}

// Using this separate struct for creation because FormJSONValues may have values of any type.
type IntegrationCreation struct {
	Integration
	FormJSONValues []jsonValues `json:"formJSONValues,omitempty"`
}

type jsonValues struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

func (is *IntegrationsService) DeleteIntegration(integrationId int) error {
	log.Debug("Deleting integration by id '" + strconv.Itoa(integrationId) + "'...")
	httpDetails := is.ServiceDetails.CreateHttpClientDetails()
	resp, body, err := is.client.SendDelete(is.ServiceDetails.GetUrl()+integrationsRestApi+strconv.Itoa(integrationId), nil, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatus(resp, body, http.StatusOK)
}

func (is *IntegrationsService) GetIntegrationById(integrationId int) (*Integration, error) {
	log.Debug("Getting integration by id '" + strconv.Itoa(integrationId) + "'...")
	httpDetails := is.ServiceDetails.CreateHttpClientDetails()
	url := is.ServiceDetails.GetUrl() + integrationsRestApi + strconv.Itoa(integrationId)
	resp, body, _, err := is.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	integration := &Integration{}
	err = json.Unmarshal(body, integration)
	return integration, errorutils.CheckError(err)
}

func (is *IntegrationsService) GetIntegrationByName(name string) (*Integration, error) {
	log.Debug("Getting integration by name '" + name + "'...")
	integrations, err := is.GetAllIntegrations()
	if err != nil {
		return nil, err
	}
	for _, integration := range integrations {
		if integration.Name == name {
			return &integration, nil
		}
	}
	return nil, errorutils.CheckErrorf("integration with provided name was not found in pipelines")
}

func (is *IntegrationsService) GetAllIntegrations() ([]Integration, error) {
	log.Debug("Fetching all integrations...")
	httpDetails := is.ServiceDetails.CreateHttpClientDetails()
	url := is.ServiceDetails.GetUrl() + integrationsRestApi
	resp, body, _, err := is.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	integrations := &[]Integration{}
	err = json.Unmarshal(body, integrations)
	return *integrations, errorutils.CheckError(err)
}

type IntegrationAlreadyExistsError struct {
	InnerError error
}

func (*IntegrationAlreadyExistsError) Error() string {
	return "Pipelines: Integration already exists."
}

type IntegrationUnauthorizedError struct {
	InnerError error
}

func (*IntegrationUnauthorizedError) Error() string {
	return "Pipelines: Integration failed, received 401 Unauthorized."
}

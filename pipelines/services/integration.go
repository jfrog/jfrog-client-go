package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
)

type IntegrationsService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewIntegrationsService(client *jfroghttpclient.JfrogHttpClient) *IntegrationsService {
	return &IntegrationsService{client: client}
}

const (
	ArtifactoryName     = "artifactory"
	artifactoryId       = 98
	GithubName          = "github"
	githubId            = 20
	githubDefaultUrl    = "https://api.github.com"
	BitbucketName       = "bitbucket"
	bitbucketId         = 16
	bitbucketDefaultUrl = "https://api.bitbucket.org"
	GitlabName          = "gitlab"
	gitlabId            = 21

	defaultProjectId = 1

	urlLabel      = "url"
	tokenLabel    = "token"
	userLabel     = "user"
	usernameLabel = "username"
	apikeyLabel   = "apikey"

	integrationsRestApi = "api/v1/projectIntegrations/"
)

func (is *IntegrationsService) CreateGithubIntegration(integrationName, token string) (id int, err error) {
	integration := Integration{
		Name:                  integrationName,
		MasterIntegrationId:   githubId,
		MasterIntegrationName: GithubName,
		ProjectId:             defaultProjectId,
		FormJSONValues: []jsonValues{
			{urlLabel, githubDefaultUrl},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateBitbucketIntegration(integrationName, username, token string) (id int, err error) {
	integration := Integration{
		Name:                  integrationName,
		MasterIntegrationId:   bitbucketId,
		MasterIntegrationName: BitbucketName,
		ProjectId:             defaultProjectId,
		FormJSONValues: []jsonValues{
			{urlLabel, bitbucketDefaultUrl},
			{usernameLabel, username},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateGitlabIntegration(integrationName, url, token string) (id int, err error) {
	integration := Integration{
		Name:                  integrationName,
		MasterIntegrationId:   gitlabId,
		MasterIntegrationName: GitlabName,
		ProjectId:             defaultProjectId,
		FormJSONValues: []jsonValues{
			{urlLabel, url},
			{tokenLabel, token},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) CreateArtifactoryIntegration(integrationName, url, user, apikey string) (id int, err error) {
	integration := Integration{
		Name:                  integrationName,
		MasterIntegrationId:   artifactoryId,
		MasterIntegrationName: ArtifactoryName,
		ProjectId:             defaultProjectId,
		FormJSONValues: []jsonValues{
			{urlLabel, url},
			{userLabel, user},
			{apikeyLabel, apikey},
		},
	}
	return is.createIntegration(integration)
}

func (is *IntegrationsService) createIntegration(integration Integration) (id int, err error) {
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err := errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body))
		if resp.StatusCode == http.StatusConflict {
			return -1, errorutils.CheckError(&IntegrationAlreadyExistsError{InnerError: err})
		}
		return -1, errorutils.CheckError(err)
	}

	created := &Integration{}
	err = json.Unmarshal(body, created)
	if err != nil {
		return -1, err
	}
	return created.Id, nil
}

type Integration struct {
	Name                  string       `json:"name,omitempty"`
	MasterIntegrationId   int          `json:"masterIntegrationId,omitempty"`
	MasterIntegrationName string       `json:"masterIntegrationName,omitempty"`
	ProjectId             int          `json:"projectId,omitempty"`
	Environments          []string     `json:"environments,omitempty"`
	FormJSONValues        []jsonValues `json:"formJSONValues,omitempty"`

	// Following fields returned when fetching or creating integration:
	Id int `json:"id,omitempty"`
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
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	return nil
}

func (is *IntegrationsService) GetIntegrationById(integrationId int) (*Integration, error) {
	log.Debug("Getting integration by id '" + strconv.Itoa(integrationId) + "'...")
	httpDetails := is.ServiceDetails.CreateHttpClientDetails()
	url := is.ServiceDetails.GetUrl() + integrationsRestApi + strconv.Itoa(integrationId)
	resp, body, _, err := is.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	integration := &Integration{}
	err = json.Unmarshal(body, integration)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	return integration, nil
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
	return nil, errorutils.CheckError(errors.New("integration with provided name was not found in pipelines"))
}

func (is *IntegrationsService) GetAllIntegrations() ([]Integration, error) {
	log.Debug("Fetching all integrations...")
	httpDetails := is.ServiceDetails.CreateHttpClientDetails()
	url := is.ServiceDetails.GetUrl() + integrationsRestApi
	resp, body, _, err := is.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	integrations := &[]Integration{}
	err = json.Unmarshal(body, integrations)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	return *integrations, nil
}

type IntegrationAlreadyExistsError struct {
	InnerError error
}

func (*IntegrationAlreadyExistsError) Error() string {
	return "Pipelines: Integration already exists."
}

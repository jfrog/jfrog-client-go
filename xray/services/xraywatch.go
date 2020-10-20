package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"

	artUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

const (
	watchAPIURL = "api/v2/watches"
)

// XrayWatchService defines the http client and xray details
type XrayWatchService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	XrayDetails auth.ServiceDetails
}

// NewXrayWatchService creates a new Xray Watch Service
func NewXrayWatchService(client *rthttpclient.ArtifactoryHttpClient) *XrayWatchService {
	return &XrayWatchService{client: client}
}

// GetJfrogHttpClient returns the http client
func (xws *XrayWatchService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return xws.client
}

// The getXrayWatchURL does not end with a slash
// So, calling functions will need to add it
func (xws *XrayWatchService) getXrayWatchURL() string {
	return clientutils.AddTrailingSlashIfNeeded(xws.XrayDetails.GetUrl()) + watchAPIURL
}

// Delete will delete an existing watch by name
// It will error if no watch can be found by that name.
func (xws *XrayWatchService) Delete(watchName string) error {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Deleting watch...")
	resp, body, err := xws.client.SendDelete(xws.getXrayWatchURL()+"/"+watchName, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Xray response:", resp.Status)
	log.Info("Done deleting watch.")
	return nil
}

// Create will create a new xray watch
func (xws *XrayWatchService) Create(params utils.XrayWatchParams) error {
	payloadBody, err := utils.CreateBody(params)
	if err != nil {
		return errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xws.getXrayWatchURL()
	var resp *http.Response
	var respBody []byte

	log.Info("Creating watch...")
	resp, respBody, err = xws.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done creating watch.")
	return nil
}

// Update will update an existing Xray watch by name
// It will error if no watch can be found by that name.
func (xws *XrayWatchService) Update(params utils.XrayWatchParams) error {
	payloadBody, err := utils.CreateBody(params)
	if err != nil {
		return errorutils.CheckError(err)
	}

	// Xray does not allow you to update a watch's name
	// The endpoint throws an error when the name is specified and the method is update.
	// Therefore, remove the name before sending the update payload
	payloadBody.GeneralData.Name = ""

	if err != nil {
		return errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xws.getXrayWatchURL() + "/" + params.Name
	var resp *http.Response
	var respBody []byte

	log.Info("Updating watch...")
	resp, respBody, err = xws.client.SendPut(url, content, &httpClientsDetails)

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done updating watch.")
	return nil
}

// Get retrieves the details about an Xray watch by its name
// It will error if no watch can be found by that name.
func (xws *XrayWatchService) Get(watchName string) (watchResp *utils.XrayWatchParams, err error) {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Getting watch...")
	resp, body, _, err := xws.client.SendGet(xws.getXrayWatchURL()+"/"+watchName, true, &httpClientsDetails)
	watch := utils.XrayWatchBody{}

	if err != nil {
		return &utils.XrayWatchParams{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &utils.XrayWatchParams{}, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &watch)

	if err != nil {
		return &utils.XrayWatchParams{}, errors.New("failed unmarshalling watch " + watchName)
	}

	result := utils.NewXrayWatchParams()
	result.Name = watch.GeneralData.Name
	result.Description = watch.GeneralData.Description
	result.Active = watch.GeneralData.Active
	result.Repositories = utils.XrayWatchRepositoriesParams{
		All:          utils.XrayWatchRepositoryAll{},
		Repositories: map[string]utils.XrayWatchRepository{},
		XrayWatchPathFilters: utils.XrayWatchPathFilters{
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
		},
	}
	result.Builds = utils.XrayWatchBuildsParams{
		All:     utils.XrayWatchBuildsAllParams{},
		ByNames: map[string]utils.XrayWatchBuildsByNameParams{},
	}
	result.Policies = watch.AssignedPolicies

	utils.UnpackWatchBody(&result, &watch)

	log.Debug("Xray response:", resp.Status)
	log.Info("Done getting watch.")

	return &result, nil
}

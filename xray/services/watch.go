package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"

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

// WatchService defines the http client and xray details
type WatchService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewWatchService creates a new Xray Watch Service
func NewWatchService(client *jfroghttpclient.JfrogHttpClient) *WatchService {
	return &WatchService{client: client}
}

// GetXrayDetails returns the xray details
func (vs *WatchService) GetXrayDetails() auth.ServiceDetails {
	return vs.XrayDetails
}

// GetJfrogHttpClient returns the http client
func (xws *WatchService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return xws.client
}

// The getWatchURL does not end with a slash
// So, calling functions will need to add it
func (xws *WatchService) getWatchURL() string {
	return clientutils.AddTrailingSlashIfNeeded(xws.XrayDetails.GetUrl()) + watchAPIURL
}

// Delete will delete an existing watch by name
// It will error if no watch can be found by that name.
func (xws *WatchService) Delete(watchName string) (*http.Response, error) {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Deleting watch...")
	resp, body, err := xws.client.SendDelete(xws.getWatchURL()+"/"+watchName, nil, &httpClientsDetails)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Xray response:", resp.Status)
	log.Info("Done deleting watch.")
	return resp, nil
}

// Create will create a new xray watch
func (xws *WatchService) Create(params utils.WatchParams) (*http.Response, error) {
	payloadBody, err := utils.CreateBody(params)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xws.getWatchURL()
	var resp *http.Response
	var respBody []byte

	log.Info("Creating watch...")
	resp, respBody, err = xws.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done creating watch.")
	return resp, nil
}

// Update will update an existing Xray watch by name
// It will error if no watch can be found by that name.
func (xws *WatchService) Update(params utils.WatchParams) (*http.Response, error) {
	payloadBody, err := utils.CreateBody(params)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	// Xray does not allow you to update a watch's name
	// The endpoint throws an error when the name is specified and the method is update.
	// Therefore, remove the name before sending the update payload
	payloadBody.GeneralData.Name = ""

	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xws.getWatchURL() + "/" + params.Name
	var resp *http.Response
	var respBody []byte

	log.Info("Updating watch...")
	resp, respBody, err = xws.client.SendPut(url, content, &httpClientsDetails)

	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done updating watch.")
	return resp, nil
}

// Get retrieves the details about an Xray watch by its name
// It will error if no watch can be found by that name.
func (xws *WatchService) Get(watchName string) (watchResp *utils.WatchParams, resp *http.Response, err error) {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Getting watch...")
	resp, body, _, err := xws.client.SendGet(xws.getWatchURL()+"/"+watchName, true, &httpClientsDetails)
	watch := utils.WatchBody{}

	if err != nil {
		return &utils.WatchParams{}, resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return &utils.WatchParams{}, resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &watch)

	if err != nil {
		return &utils.WatchParams{}, resp, errors.New("failed unmarshalling watch " + watchName)
	}

	result := utils.NewWatchParams()
	result.Name = watch.GeneralData.Name
	result.Description = watch.GeneralData.Description
	result.Active = watch.GeneralData.Active
	result.Repositories = utils.WatchRepositoriesParams{
		All:          utils.WatchRepositoryAll{},
		Repositories: map[string]utils.WatchRepository{},
		WatchPathFilters: utils.WatchPathFilters{
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
		},
	}
	result.Builds = utils.WatchBuildsParams{
		All:     utils.WatchBuildsAllParams{},
		ByNames: map[string]utils.WatchBuildsByNameParams{},
	}
	result.Policies = watch.AssignedPolicies

	utils.UnpackWatchBody(&result, &watch)

	log.Debug("Xray response:", resp.Status)
	log.Info("Done getting watch.")

	return &result, resp, nil
}

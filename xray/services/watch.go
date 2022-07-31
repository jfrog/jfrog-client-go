package services

import (
	"encoding/json"
	"errors"
	"fmt"
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

// WatchService defines the http client and Xray details
type WatchService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

type WatchAlreadyExistsError struct {
	InnerError error
}

func (*WatchAlreadyExistsError) Error() string {
	return "Xray: Watch already exists."
}

// NewWatchService creates a new Xray Watch Service
func NewWatchService(client *jfroghttpclient.JfrogHttpClient) *WatchService {
	return &WatchService{client: client}
}

// GetXrayDetails returns the Xray details
func (xws *WatchService) GetXrayDetails() auth.ServiceDetails {
	return xws.XrayDetails
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
func (xws *WatchService) Delete(watchName string) error {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Deleting watch...")
	resp, body, err := xws.client.SendDelete(xws.getWatchURL()+"/"+watchName, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	log.Debug("Xray response:", resp.Status)
	log.Info("Done deleting watch.")
	return nil
}

// Create will create a new Xray watch
func (xws *WatchService) Create(params utils.WatchParams) error {
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
	var url = xws.getWatchURL()

	log.Info(fmt.Sprintf("Creating a new Watch named %s on JFrog Xray....", params.Name))
	resp, body, err := xws.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		if resp.StatusCode == http.StatusConflict {
			return &WatchAlreadyExistsError{InnerError: err}
		}
		return err
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done creating watch.")
	return nil
}

// Update will update an existing Xray watch by name
// It will error if no watch can be found by that name.
func (xws *WatchService) Update(params utils.WatchParams) error {
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
	var url = xws.getWatchURL() + "/" + params.Name

	log.Info("Updating watch...")
	resp, body, err := xws.client.SendPut(url, content, &httpClientsDetails)

	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		return err
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done updating watch.")
	return nil
}

// Get retrieves the details about an Xray watch by its name
// It will error if no watch can be found by that name.
func (xws *WatchService) Get(watchName string) (watchResp *utils.WatchParams, err error) {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Getting watch...")
	resp, body, _, err := xws.client.SendGet(xws.getWatchURL()+"/"+watchName, true, &httpClientsDetails)
	watch := utils.WatchBody{}

	if err != nil {
		return &utils.WatchParams{}, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return &utils.WatchParams{}, err
	}
	err = json.Unmarshal(body, &watch)

	if err != nil {
		return &utils.WatchParams{}, errors.New("failed unmarshalling watch " + watchName)
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

	return &result, nil
}

package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type PermissionTargetService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewPermissionTargetService(client *jfroghttpclient.JfrogHttpClient) *PermissionTargetService {
	return &PermissionTargetService{client: client}
}

func (pts *PermissionTargetService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return pts.client
}

func (pts *PermissionTargetService) Delete(permissionTargetName string) error {
	httpClientsDetails := pts.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting permission target...")
	resp, body, err := pts.client.SendDelete(pts.ArtDetails.GetUrl()+"api/v2/security/permissions/"+permissionTargetName, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done deleting permission target.")
	return nil
}

func (pts *PermissionTargetService) Get(permissionTargetName string) (*PermissionTargetParams, error) {
	httpClientsDetails := pts.ArtDetails.CreateHttpClientDetails()
	log.Info("Getting permission target...")
	resp, body, _, err := pts.client.SendGet(pts.ArtDetails.GetUrl()+"api/v2/security/permissions/"+permissionTargetName, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	// The case the requested permission target is not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	log.Debug("Artifactory response:", resp.Status)
	permissionTarget := &PermissionTargetParams{}
	if err := json.Unmarshal(body, permissionTarget); err != nil {
		return nil, err
	}
	return permissionTarget, nil
}

func (pts *PermissionTargetService) Create(params PermissionTargetParams) error {
	return pts.performRequest(params, false)
}

func (pts *PermissionTargetService) Update(params PermissionTargetParams) error {
	return pts.performRequest(params, true)
}

func (pts *PermissionTargetService) performRequest(params PermissionTargetParams, update bool) error {
	content, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	httpClientsDetails := pts.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = pts.ArtDetails.GetUrl() + "api/v2/security/permissions/" + params.Name
	var operationString string
	var resp *http.Response
	var body []byte
	if update {
		log.Info("Updating permission target...")
		operationString = "updating"
		resp, body, err = pts.client.SendPut(url, content, &httpClientsDetails)
	} else {
		log.Info("Creating permission target...")
		operationString = "creating"
		resp, body, err = pts.client.SendPost(url, content, &httpClientsDetails)
	}
	if err != nil {
		return err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		if resp.StatusCode == http.StatusConflict {
			return &PermissionTargetAlreadyExistsError{InnerError: err}
		}
		return err
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done " + operationString + " permission target.")
	return nil
}

func NewPermissionTargetParams() PermissionTargetParams {
	return PermissionTargetParams{}
}

// Using struct pointers to keep the fields null if they are empty.
// Artifactory evaluates inner struct typed fields if they are not null, which can lead to failures in the request.
type PermissionTargetParams struct {
	Name          string                   `json:"name"`
	Repo          *PermissionTargetSection `json:"repo,omitempty"`
	Build         *PermissionTargetSection `json:"build,omitempty"`
	ReleaseBundle *PermissionTargetSection `json:"releaseBundle,omitempty"`
}

type PermissionTargetSection struct {
	IncludePatterns []string `json:"include-patterns,omitempty"`
	ExcludePatterns []string `json:"exclude-patterns,omitempty"`
	Repositories    []string `json:"repositories"`
	Actions         *Actions `json:"actions,omitempty"`
}

type Actions struct {
	Users  map[string][]string `json:"users,omitempty"`
	Groups map[string][]string `json:"groups,omitempty"`
}

type PermissionTargetAlreadyExistsError struct {
	InnerError error
}

func (*PermissionTargetAlreadyExistsError) Error() string {
	return "Artifactory: Permission target already exists."
}

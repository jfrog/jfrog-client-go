package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

// VersionService returns the https client and Xsc details
type VersionService struct {
	client     *jfroghttpclient.JfrogHttpClient
	XscDetails auth.ServiceDetails
}

// NewVersionService creates a new service to retrieve the version of Xsc
func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

// GetXscDetails returns the Xsc details
func (vs *VersionService) GetXscDetails() auth.ServiceDetails {
	return vs.XscDetails
}

// GetVersion returns the version of Xsc
func (vs *VersionService) GetVersion() (string, error) {
	httpDetails := vs.XscDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.XscDetails.GetUrl()+"api/v1/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var version xscVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

type xscVersion struct {
	Version  string `json:"xsc_version,omitempty"`
	Revision string `json:"xray_version,omitempty"`
}

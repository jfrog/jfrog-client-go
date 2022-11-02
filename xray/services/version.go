package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

// VersionService returns the https client and Xray details
type VersionService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewVersionService creates a new service to retrieve the version of Xray
func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

// GetXrayDetails returns the Xray details
func (vs *VersionService) GetXrayDetails() auth.ServiceDetails {
	return vs.XrayDetails
}

// GetVersion returns the version of Xray
func (vs *VersionService) GetVersion() (string, error) {
	httpDetails := vs.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.XrayDetails.GetUrl()+"api/v1/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", errors.New("failed while attempting to get Xray version:\n" + err.Error())
	}
	var version xrayVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckErrorf("couldn't parse Xray server response: " + err.Error())
	}
	return strings.TrimSpace(version.Version), nil
}

type xrayVersion struct {
	Version  string `json:"xray_version,omitempty"`
	Revision string `json:"xray_revision,omitempty"`
}

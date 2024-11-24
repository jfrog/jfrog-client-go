package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const (
	xscVersionApiSuffix           = "system/version"
	xscDeprecatedVersionApiSuffix = "api/v1/" + xscVersionApiSuffix
)

// VersionService returns the https client and Xsc details
type VersionService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XscDetails  auth.ServiceDetails
	XrayDetails auth.ServiceDetails
}

// NewVersionService creates a new service to retrieve the version of Xsc
func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

func (vs *VersionService) sendVersionRequest() (resp *http.Response, body []byte, err error) {
	if vs.XrayDetails != nil {
		httpDetails := vs.XrayDetails.CreateHttpClientDetails()
		resp, body, _, err = vs.client.SendGet(vs.XrayDetails.GetUrl()+xscutils.XscInXraySuffix+xscVersionApiSuffix, true, &httpDetails)
		return
	}
	// Backward compatibility
	httpDetails := vs.XscDetails.CreateHttpClientDetails()
	resp, body, _, err = vs.client.SendGet(vs.XscDetails.GetUrl()+xscDeprecatedVersionApiSuffix, true, &httpDetails)
	return
}

// GetVersion returns the version of Xsc
func (vs *VersionService) GetVersion() (string, error) {
	resp, body, err := vs.sendVersionRequest()
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		if resp.StatusCode == http.StatusNotFound {
			// When XSC is disabled, StatusNotFound is expected. Don't GenerateResponseError in this case.
			return "", fmt.Errorf("xsc is not available for this server")
		}
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

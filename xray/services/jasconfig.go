package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	jasConfigApiURL = "api/v1/configuration/jas"
)

// JasConfigService returns the https client and Xray details
type JasConfigService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewJasConfigService creates a new service to retrieve the version of Xray
func NewJasConfigService(client *jfroghttpclient.JfrogHttpClient) *JasConfigService {
	return &JasConfigService{client: client}
}

// GetXrayDetails returns the Xray details
func (jcs *JasConfigService) GetXrayDetails() auth.ServiceDetails {
	return jcs.XrayDetails
}

// GetJasConfigTokenValidation returns token validation status in xray
func (jcs *JasConfigService) GetJasConfigTokenValidation() (bool, error) {
	httpDetails := jcs.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := jcs.client.SendGet(jcs.XrayDetails.GetUrl()+jasConfigApiURL, true, &httpDetails)
	if err != nil {
		return false, errors.New("failed while attempting to get JFrog Xray Jas Configuration: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return false, errors.New("got unexpected server response while attempting to get JFrog Xray Jas Configuration:\n" + err.Error())
	}
	var jasConfig JasConfig
	if err = json.Unmarshal(body, &jasConfig); err != nil {
		return false, errorutils.CheckErrorf("couldn't parse JFrog Xray server Jas Configuration response: " + err.Error())
	}
	return *jasConfig.TokenValidationToggle, nil
}

type JasConfig struct {
	TokenValidationToggle *bool `json:"enable_token_validation_scanning,omitempty"`
}

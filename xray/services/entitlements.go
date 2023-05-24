package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type EntitlementsService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewEntitlementsService creates a new service to retrieve the entitlement data from Xray
func NewEntitlementsService(client *jfroghttpclient.JfrogHttpClient) *EntitlementsService {
	return &EntitlementsService{client: client}
}

// GetXrayDetails returns the Xray details
func (es *EntitlementsService) GetXrayDetails() auth.ServiceDetails {
	return es.XrayDetails
}

// IsEntitled returns true if the user is entitled for the requested feature ID
func (es *EntitlementsService) IsEntitled(featureId string) (entitled bool, err error) {
	httpDetails := es.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := es.client.SendGet(es.XrayDetails.GetUrl()+"api/v1/entitlements/feature/"+featureId, true, &httpDetails)
	if err != nil {
		return
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return false, fmt.Errorf("failed while attempting to get Xray entitlements response for %s:\n%s", featureId, err.Error())
	}
	var userEntitlements entitlements
	if err = json.Unmarshal(body, &userEntitlements); err != nil {
		return false, errorutils.CheckErrorf("couldn't parse Xray server response: " + err.Error())
	}
	return userEntitlements.Entitled, nil
}

type entitlements struct {
	FeatureId string `json:"feature_id,omitempty"`
	Entitled  bool   `json:"entitled,omitempty"`
}

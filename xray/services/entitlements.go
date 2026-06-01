package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type EntitlementsService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

// NewEntitlementsService creates a new service to retrieve the entitlement data from Xray
func NewEntitlementsService(client *jfroghttpclient.JfrogHttpClient) *EntitlementsService {
	return &EntitlementsService{client: client}
}

// GetXrayDetails returns the Xray details
func (es *EntitlementsService) GetXrayDetails() auth.ServiceDetails {
	return es.XrayDetails
}

func (es *EntitlementsService) getUrlForEntitlementApi(featureId string) string {
	return clientutils.AppendScopedProjectKeyParam(es.XrayDetails.GetUrl()+"api/v1/entitlements/feature/"+featureId, es.ScopeProjectKey)
}

// IsEntitled returns true if the user is entitled for the requested feature ID
func (es *EntitlementsService) IsEntitled(featureId string) (entitled bool, err error) {
	httpDetails := es.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := es.client.SendGet(es.getUrlForEntitlementApi(featureId), true, &httpDetails)
	if err != nil {
		err = errors.New("failed while attempting to get JFrog Xray entitlements response: " + err.Error())
		return
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		err = fmt.Errorf("got unexpected server response while attempting to get JFrog Xray entitlements response for %s:\n%s", featureId, err.Error())
		return
	}
	var userEntitlements entitlements
	if err = json.Unmarshal(body, &userEntitlements); err != nil {
		err = errorutils.CheckErrorf("couldn't parse JFrog Xray server entitlements response: %s", err.Error())
		return
	}
	entitled = userEntitlements.Entitled
	return
}

type entitlements struct {
	FeatureId string `json:"feature_id,omitempty"`
	Entitled  bool   `json:"entitled,omitempty"`
}

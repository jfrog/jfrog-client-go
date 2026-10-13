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

const applicabilityInputApi = "api/v1/internal/cve_applicability_input_vulnerabilities"

type ApplicabilityService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

// NewApplicabilityService creates a new service to retrieve the contextual analysis input from Xray
func NewApplicabilityService(client *jfroghttpclient.JfrogHttpClient) *ApplicabilityService {
	return &ApplicabilityService{client: client}
}

// GetApplicabilityInput returns the contextual analysis input of the given vulnerabilities
func (as *ApplicabilityService) GetApplicabilityInput(vulnerabilityIds []string) (*ApplicabilityInput, error) {
	requestBody, err := json.Marshal(applicabilityInputRequest{Vulnerabilities: vulnerabilityIds})
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	httpDetails := as.XrayDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()
	resp, body, err := as.client.SendPost(clientutils.AppendScopedProjectKeyParam(as.XrayDetails.GetUrl()+applicabilityInputApi, as.ScopeProjectKey), requestBody, &httpDetails)
	if err != nil {
		return nil, errors.New("failed while attempting to get JFrog Xray applicability input: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("got unexpected server response while attempting to get JFrog Xray applicability input:\n%s", err.Error())
	}
	var applicabilityInput ApplicabilityInput
	if err = json.Unmarshal(body, &applicabilityInput); err != nil {
		return nil, errorutils.CheckErrorf("couldn't parse JFrog Xray applicability input response: %s", err.Error())
	}
	return &applicabilityInput, nil
}

type applicabilityInputRequest struct {
	Vulnerabilities []string `json:"vulnerabilities"`
}

type ApplicabilityInput struct {
	Explanation map[string]string `json:"explanation,omitempty"`
}

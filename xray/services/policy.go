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
	policyAPIURL = "api/v2/policies"
)

// PolicyService defines the http client and xray details
type PolicyService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewPolicyService creates a new Xray Policy Service
func NewPolicyService(client *jfroghttpclient.JfrogHttpClient) *PolicyService {
	return &PolicyService{client: client}
}

// GetXrayDetails returns the xray details
func (xps *PolicyService) GetXrayDetails() auth.ServiceDetails {
	return xps.XrayDetails
}

// GetJfrogHttpClient returns the http client
func (xps *PolicyService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return xps.client
}

// The getPolicyURL does not end with a slash
// So, calling functions will need to add it
func (xps *PolicyService) getPolicyURL() string {
	return clientutils.AddTrailingSlashIfNeeded(xps.XrayDetails.GetUrl()) + policyAPIURL
}

// Delete will delete an existing policy by name
// It will error if no policy can be found by that name.
func (xps *PolicyService) Delete(policyName string) (*http.Response, error) {
	httpClientsDetails := xps.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)

	log.Info("Deleting policy...")
	resp, body, err := xps.client.SendDelete(xps.getPolicyURL()+"/"+policyName, nil, &httpClientsDetails)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Xray response:", resp.Status)
	log.Info("Done deleting policy.")
	return resp, nil
}

// Create will create a new xray policy
func (xps *PolicyService) Create(params utils.PolicyParams) (*http.Response, error) {
	body := utils.CreatePolicyBody(params)
	content, err := json.Marshal(body)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientsDetails := xps.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xps.getPolicyURL()
	var resp *http.Response
	var respBody []byte

	log.Info("Creating policy...")
	resp, respBody, err = xps.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done creating policy.")
	return resp, nil
}

// Update will update an existing Xray policy by name
// It will error if no policy can be found by that name.
func (xps *PolicyService) Update(params utils.PolicyParams) (*http.Response, error) {
	body := utils.CreatePolicyBody(params)
	content, err := json.Marshal(body)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientsDetails := xps.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xps.getPolicyURL() + "/" + params.Name
	var resp *http.Response
	var respBody []byte

	log.Info("Updating policy...")
	resp, respBody, err = xps.client.SendPut(url, content, &httpClientsDetails)

	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done updating policy.")
	return resp, nil
}

// Get retrieves the details about an Xray policy by its name
// It will error if no policy can be found by that name.
func (xps *PolicyService) Get(policyName string) (policyResp *utils.PolicyParams, resp *http.Response, err error) {
	httpClientsDetails := xps.XrayDetails.CreateHttpClientDetails()
	log.Info("Getting policy...")
	resp, body, _, err := xps.client.SendGet(xps.getPolicyURL()+"/"+policyName, true, &httpClientsDetails)
	policy := &utils.PolicyBody{}

	if err != nil {
		return &utils.PolicyParams{}, resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return &utils.PolicyParams{}, resp, errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, policy)

	if err != nil {
		return &utils.PolicyParams{}, resp, errors.New("failed unmarshalling policy " + policyName)
	}

	log.Debug("Xray response:", resp.Status)
	log.Info("Done getting policy.")

	return &utils.PolicyParams{
		Name:        policy.Name,
		Type:        policy.Type,
		Description: policy.Description,
		Rules:       policy.Rules,
	}, resp, nil
}

package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	artUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

const (
	ignoreRuleAPIURL = "api/v1/ignore_rules"
)

// IgnoreRuleService defines the http client and Xray details
type IgnoreRuleService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewIgnoreRuleService creates a new Xray Ignore Rule Service
func NewIgnoreRuleService(client *jfroghttpclient.JfrogHttpClient) *IgnoreRuleService {
	return &IgnoreRuleService{client: client}
}

// GetXrayDetails returns the Xray details
func (xirs *IgnoreRuleService) GetXrayDetails() auth.ServiceDetails {
	return xirs.XrayDetails
}

// GetJfrogHttpClient returns the http client
func (xirs *IgnoreRuleService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return xirs.client
}

// The getIgnoreRuleURL does not end with a slash
// So, calling functions will need to add it
func (xirs *IgnoreRuleService) getIgnoreRuleURL() string {
	return clientutils.AddTrailingSlashIfNeeded(xirs.XrayDetails.GetUrl()) + ignoreRuleAPIURL
}

// Delete will delete an ignore rule by id
func (xirs *IgnoreRuleService) Delete(ignoreRuleId string) error {
	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)

	log.Info("Deleting ignore rule...")
	resp, body, err := xirs.client.SendDelete(xirs.getIgnoreRuleURL()+"/"+ignoreRuleId, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Debug("Xray response:", resp.Status)
	log.Info("Done deleting ignore rule.")
	return nil
}

// Create will create a new Xray ignore rule
// The function creates the ignore rule and returns its id which is recieved after post
func (xirs *IgnoreRuleService) Create(params utils.IgnoreRuleParams, ignoreRuleId *string) error {
	ignoreRuleBody := utils.CreateIgnoreRuleBody(params)
	content, err := json.Marshal(ignoreRuleBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xirs.getIgnoreRuleURL()

	log.Info("Create new ignore rule...")
	resp, body, err := xirs.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		return err
	}

	*ignoreRuleId = getIgnoreRuleIdFromBody(body)
	log.Debug("Xray response:", resp.Status)
	log.Info("Done creating ignore rule.")
	return nil
}

func getIgnoreRuleIdFromBody(body []byte) string {
	str := string(body)

	re := regexp.MustCompile(`id:\s*([a-f0-9-]+)`)
	match := re.FindStringSubmatch(str)

	if len(match) != 0 {
		return match[1]
	}

	return ""
}

// Get retrieves the details about an Xray ignore rule by its id
// It will error if no ignore rule can be found by that id.
func (xirs *IgnoreRuleService) Get(ignoreRuleId string) (ignoreRuleResp *utils.IgnoreRuleParams, err error) {
	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	log.Info("Getting ignore rule...")
	resp, body, _, err := xirs.client.SendGet(xirs.getIgnoreRuleURL()+"/"+ignoreRuleId, true, &httpClientsDetails)
	ignoreRule := &utils.IgnoreRuleBody{}

	if err != nil {
		return &utils.IgnoreRuleParams{}, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return &utils.IgnoreRuleParams{}, err
	}

	err = json.Unmarshal(body, ignoreRule)

	if err != nil {
		return &utils.IgnoreRuleParams{}, errors.New("failed unmarshalling ignore rule " + ignoreRuleId)
	}

	log.Debug("Xray response:", resp.Status)
	log.Info("Done getting ignore rule.")

	return &utils.IgnoreRuleParams{
		Notes:         ignoreRule.Notes,
		ExpiresAt:     ignoreRule.ExpiresAt,
		IgnoreFilters: ignoreRule.IgnoreFilters,
	}, nil
}

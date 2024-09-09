package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

const (
	ignoreRuleAPIURL          = "api/v1/ignore_rules"
	minXrayIgnoreRulesVersion = "3.11"
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

func (xirs *IgnoreRuleService) CheckMinimumVersion() error {
	xrDetails := xirs.GetXrayDetails()
	if xrDetails == nil {
		return errorutils.CheckErrorf("Xray details not configured.")
	}
	version, err := xrDetails.GetVersion()
	if err != nil {
		return fmt.Errorf("couldn't get Xray version. Error: %w", err)
	}

	return clientutils.ValidateMinimumVersion(clientutils.Xray, version, minXrayIgnoreRulesVersion)
}

// The getIgnoreRuleURL does not end with a slash
// So, calling functions will need to add it
func (xirs *IgnoreRuleService) getIgnoreRuleURL() string {
	return clientutils.AddTrailingSlashIfNeeded(xirs.XrayDetails.GetUrl()) + ignoreRuleAPIURL
}

// Delete will delete an ignore rule by id
func (xirs *IgnoreRuleService) Delete(ignoreRuleId string) error {
	if err := xirs.CheckMinimumVersion(); err != nil {
		return err
	}
	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	log.Info("Deleting ignore rule...")
	resp, body, err := xirs.client.SendDelete(xirs.getIgnoreRuleURL()+"/"+ignoreRuleId, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent); err != nil {
		return err
	}
	log.Debug("Xray response status:", resp.Status)
	log.Info("Done deleting ignore rule.")
	return nil
}

// Create will create a new Xray ignore rule
// The function creates the ignore rule and returns its id which is received after post
func (xirs *IgnoreRuleService) Create(params utils.IgnoreRuleParams) (ignoreRuleId string, err error) {
	if err := xirs.CheckMinimumVersion(); err != nil {
		return "", err
	}
	ignoreRuleBody := utils.CreateIgnoreRuleBody(params)
	content, err := json.Marshal(ignoreRuleBody)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()
	var url = xirs.getIgnoreRuleURL()

	log.Info("Create new ignore rule...")
	resp, body, err := xirs.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return "", err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
		return "", err
	}
	log.Debug("Xray response status:", resp.Status)

	ignoreRuleId, err = getIgnoreRuleIdFromBody(body)
	if err != nil {
		return "", err
	}

	log.Info("Done creating ignore rule.")
	log.Debug("Ignore rule id is: ", ignoreRuleId)

	return ignoreRuleId, nil
}

func getIgnoreRuleIdFromBody(body []byte) (string, error) {
	str := string(body)

	re := regexp.MustCompile(`id:\s*([a-f0-9-]+)`)
	match := re.FindStringSubmatch(str)

	if len(match) <= 1 {
		return "", errorutils.CheckErrorf("couldn't find id for ignore rule in str: %s", str)
	}

	return match[1], nil
}

// Get retrieves the details about an Xray ignore rule by its id
// It will error if the ignore rule id can't be found.
func (xirs *IgnoreRuleService) Get(ignoreRuleId string) (ignoreRuleResp *utils.IgnoreRuleBody, err error) {
	if err = xirs.CheckMinimumVersion(); err != nil {
		return nil, err
	}
	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	log.Info(fmt.Sprintf("Getting ignore rule '%s'...", ignoreRuleId))
	resp, body, _, err := xirs.client.SendGet(xirs.getIgnoreRuleURL()+"/"+ignoreRuleId, true, &httpClientsDetails)
	ignoreRule := &utils.IgnoreRuleBody{}

	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		log.Debug("Xray response:", string(body), resp.Status)
		return nil, err
	}

	if err = json.Unmarshal(body, ignoreRule); err != nil {
		return nil, errorutils.CheckErrorf("failed unmarshalling %s for ignore rule %s", string(body), ignoreRuleId)
	}

	log.Debug("Xray response status:", resp.Status)
	log.Info("Done getting ignore rule.")

	return ignoreRule, nil
}

// GetAll retrieves the details about all Xray ignore rules that match the given filters
func (xirs *IgnoreRuleService) GetAll(params *utils.IgnoreRulesGetAllParams) (ignoreRules *utils.IgnoreRuleResponse, err error) {
	if err = xirs.CheckMinimumVersion(); err != nil {
		return nil, err
	}
	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	url, err := clientutils.BuildUrl(clientutils.AddTrailingSlashIfNeeded(xirs.XrayDetails.GetUrl()), ignoreRuleAPIURL, params.GetParamMap())
	if err != nil {
		return nil, err
	}
	resp, body, _, err := xirs.client.SendGet(url, true, &httpClientsDetails)
	ignoreRules = &utils.IgnoreRuleResponse{}
	if err != nil {
		return nil, err
	}
	log.Debug("Xray response:", resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, ignoreRules)
	if err != nil {
		return nil, errorutils.CheckErrorf("failed unmarshalling ignoreRules")
	}

	return ignoreRules, nil
}

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
	ignoreRuleBody := utils.CreateIgnoreRuleBody(params)
	if err = validateIgnoreFilters(ignoreRuleBody.IgnoreFilters); err != nil {
		return "", err
	}
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

func validateIgnoreFilters(ignoreFilters utils.IgnoreFilters) error {
	filters := []string{}
	if len(ignoreFilters.CVEs) > 0 {
		filters = append(filters, "CVEs")
	}
	if ignoreFilters.Exposures != nil {
		filters = append(filters, "Exposures")
	}
	if ignoreFilters.Sast != nil {
		filters = append(filters, "Sast")
	}
	// if more than one filter is set, notify the user
	if len(filters) > 1 {
		return errorutils.CheckErrorf("more than one ignore filter is set, split them to multiple ignore rules: %v", filters)
	}
	return nil
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
func (xirs *IgnoreRuleService) Get(ignoreRuleId string) (ignoreRuleResp *utils.IgnoreRuleParams, err error) {
	httpClientsDetails := xirs.XrayDetails.CreateHttpClientDetails()
	log.Info(fmt.Sprintf("Getting ignore rule '%s'...", ignoreRuleId))
	resp, body, _, err := xirs.client.SendGet(xirs.getIgnoreRuleURL()+"/"+ignoreRuleId, true, &httpClientsDetails)
	ignoreRule := &utils.IgnoreRuleBody{}

	if err != nil {
		return &utils.IgnoreRuleParams{}, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		log.Debug("Xray response:", string(body), resp.Status)
		return &utils.IgnoreRuleParams{}, err
	}

	if err = json.Unmarshal(body, ignoreRule); err != nil {
		return &utils.IgnoreRuleParams{}, errorutils.CheckErrorf("failed unmarshalling %s for ignore rule %s", string(body), ignoreRuleId)
	}

	log.Debug("Xray response status:", resp.Status)
	log.Info("Done getting ignore rule.")

	return &ignoreRule.IgnoreRuleParams, nil
}

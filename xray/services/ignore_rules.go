package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	artUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	ignoreRulesUrl            = "api/v1/ignore_rules"
	uuidRegEx                 = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"
	minXrayIgnoreRulesVersion = "3.11"
)

// IgnoreRulesService defines the http client and Xray details
type IgnoreRulesService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

type IgnoreRuleNotFoundError struct {
	InnerError error
}

func (e IgnoreRuleNotFoundError) Error() string {
	innerErrorText := ""
	if e.InnerError != nil {
		innerErrorText = e.InnerError.Error()
	}
	return fmt.Sprintf("Xray: ignore rule not found. %s", innerErrorText)
}

// NewIgnoreRulesService creates a new Xray Policy Service
func NewIgnoreRulesService(client *jfroghttpclient.JfrogHttpClient) *IgnoreRulesService {
	return &IgnoreRulesService{client: client}
}

// GetXrayDetails returns the Xray details
func (irs *IgnoreRulesService) GetXrayDetails() auth.ServiceDetails {
	return irs.XrayDetails
}

// GetJfrogHttpClient returns the http client
func (irs *IgnoreRulesService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return irs.client
}

func (irs *IgnoreRulesService) CheckMinimumVersion() error {
	xrDetails := irs.GetXrayDetails()
	if xrDetails == nil {
		return errorutils.CheckErrorf("Xray details not configured.")
	}
	version, err := xrDetails.GetVersion()
	if err != nil {
		return fmt.Errorf("couldn't get Xray version. Error: %w", err)
	}

	return clientUtils.ValidateMinimumVersion(clientUtils.Xray, version, minXrayIgnoreRulesVersion)
}

func (irs *IgnoreRulesService) getIgnoreRulesURL() string {
	return fmt.Sprintf("%s%s", irs.XrayDetails.GetUrl(), ignoreRulesUrl)
}

// Delete will delete an existing ignore rule by ruleId
// It will error if no ignore rule can be found by that ruleId.
func (irs *IgnoreRulesService) Delete(ruleId string) error {
	if err := irs.CheckMinimumVersion(); err != nil {
		return err
	}
	httpClientsDetails := irs.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)

	resp, body, err := irs.client.SendDelete(irs.getRuleIdUrl(ruleId), nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			notFound := IgnoreRuleNotFoundError{InnerError: err}
			return notFound
		}
		return err
	}
	return nil
}

// Create will create a new Xray ignore rule
func (irs *IgnoreRulesService) Create(ignoreRule IgnoreRule) (string, error) {
	if err := irs.CheckMinimumVersion(); err != nil {
		return "", err
	}
	content, err := json.Marshal(ignoreRule)
	if err != nil {
		return "", errorutils.CheckErrorf("error unmarshalling ignore rule %w", err)
	}

	httpClientsDetails := irs.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)

	resp, body, err := irs.client.SendPost(irs.getIgnoreRulesURL(), content, &httpClientsDetails)
	if err != nil {
		return "", err
	}
	if resp != nil {
		log.Debug("Xray response:", resp.Status)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
		return "", err
	}
	responseBody := struct {
		Info string `json:"info"`
	}{}
	if err = json.Unmarshal(body, &responseBody); err != nil {
		return "", errorutils.CheckError(err)
	}
	uuid := regexp.MustCompile(uuidRegEx).FindString(responseBody.Info)

	return uuid, nil
}

// Get retrieves the details about an Xray ignore rule by its id
// It will error if no policy can be found by that name.
func (irs *IgnoreRulesService) Get(ruleId string) (ignoreRule *IgnoreRuleDetail, err error) {
	if err = irs.CheckMinimumVersion(); err != nil {
		return nil, err
	}
	httpClientsDetails := irs.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := irs.client.SendGet(irs.getRuleIdUrl(ruleId), true, &httpClientsDetails)
	ignoreRule = &IgnoreRuleDetail{}
	if err != nil {
		return nil, err
	}
	log.Debug("Xray response:", resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, ignoreRule)
	if err != nil {
		return nil, errors.New("failed unmarshalling ignore rules " + ruleId)
	}

	return ignoreRule, nil
}

// GetAll retrieves the details about all Xray ignore rules that match the given parameters
func (irs *IgnoreRulesService) GetAll(params *IgnoreRulesGetAllParams) (ignoreRules *IgnoreRuleResponse, err error) {
	if err = irs.CheckMinimumVersion(); err != nil {
		return nil, err
	}
	httpClientsDetails := irs.XrayDetails.CreateHttpClientDetails()
	url, err := clientUtils.BuildUrl(irs.XrayDetails.GetUrl(), ignoreRulesUrl, params.getParamMap())
	if err != nil {
		return nil, err
	}
	resp, body, _, err := irs.client.SendGet(url, true, &httpClientsDetails)
	ignoreRules = &IgnoreRuleResponse{}
	if err != nil {
		return nil, err
	}
	log.Debug("Xray response:", resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, ignoreRules)
	if err != nil {
		return nil, errors.New("failed unmarshalling ignoreRules")
	}

	return ignoreRules, nil
}

func (irs *IgnoreRulesService) getRuleIdUrl(ruleId string) string {
	return fmt.Sprintf("%s/%s", irs.getIgnoreRulesURL(), ruleId)
}

func (p *IgnoreRulesGetAllParams) getParamMap() map[string]string {
	params := make(map[string]string)
	if p == nil {
		return params
	}
	if p.Vulnerability != "" {
		params["vulnerability"] = p.Vulnerability
	}
	if p.License != "" {
		params["license"] = p.License
	}
	if p.Policy != "" {
		params["policy"] = p.Policy
	}
	if p.Watch != "" {
		params["watch"] = p.Watch
	}
	if p.ComponentName != "" {
		params["component_name"] = p.ComponentName
	}
	if p.ComponentVersion != "" {
		params["component_version"] = p.ComponentVersion
	}
	if p.ArtifactName != "" {
		params["artifact_name"] = p.ArtifactName
	}
	if p.ArtifactVersion != "" {
		params["artifact_version"] = p.ArtifactVersion
	}
	if p.BuildName != "" {
		params["build_name"] = p.BuildName
	}
	if p.BuildVersion != "" {
		params["build_version"] = p.BuildVersion
	}
	if p.ReleaseBundleName != "" {
		params["release_bundle_name"] = p.ReleaseBundleName
	}
	if p.ReleaseBundleVersion != "" {
		params["release_bundle_version"] = p.ReleaseBundleVersion
	}
	if p.DockerLayer != "" {
		params["docker_layer"] = p.DockerLayer
	}
	if p.OrderBy != "" {
		params["order_by"] = p.OrderBy
	}
	if p.Direction != "" {
		params["direction"] = p.Direction
	}
	if p.PageNum != 0 {
		params["page_num"] = fmt.Sprintf("%d", p.PageNum)
	}
	if p.NumOfRows != 0 {
		params["num_of_rows"] = fmt.Sprintf("%d", p.NumOfRows)
	}
	if !p.ExpiresBefore.IsZero() {
		params["expires_before"] = p.ExpiresBefore.UTC().Format(time.RFC3339)
	}
	if !p.ExpiresAfter.IsZero() {
		params["expires_after"] = p.ExpiresAfter.UTC().Format(time.RFC3339)
	}
	if p.ProjectKey != "" {
		params["project_key"] = p.ProjectKey
	}

	return params
}

// IgnoreRuleResponse struct representing the entire JSON
type IgnoreRuleResponse struct {
	Data       []IgnoreRuleDetail `json:"data"`
	TotalCount int                `json:"total_count"`
}

// IgnoreRule struct representing an Ignore Rule
type IgnoreRule struct {
	Notes     string        `json:"notes"`
	ExpiresAt *time.Time    `json:"expires_at,omitempty"`
	Filters   IgnoreFilters `json:"ignore_filters"`
}

// IgnoreRuleDetail struct representing an Ignore Rule as returned by the API
type IgnoreRuleDetail struct {
	IgnoreRule
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Created   time.Time `json:"created"`
	IsExpired bool      `json:"is_expired"`
}

// IgnoreFilters struct representing the "ignore_filters" object
type IgnoreFilters struct {
	ReleaseBundles  []NameVersion        `json:"release_bundles,omitempty"`
	Builds          []NameVersion        `json:"builds,omitempty"`
	Components      []NameVersion        `json:"components,omitempty"`
	Artifacts       []ArtifactDescriptor `json:"artifacts,omitempty"`
	Policies        []string             `json:"policies,omitempty"`
	DockerLayers    []string             `json:"docker_layers,omitempty"`
	Vulnerabilities []string             `json:"vulnerabilities,omitempty"`
	Licenses        []string             `json:"licenses,omitempty"`
	CVEs            []string             `json:"cves,omitempty"`
	Watches         []string             `json:"watches,omitempty"`
	OperationalRisk []string             `json:"operational_risk,omitempty"`
}

// NameVersion struct representing items with a Name / Version combo
type NameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ArtifactDescriptor struct representing each item in the "artifacts" array
type ArtifactDescriptor struct {
	NameVersion
	Path string `json:"path"`
}

type IgnoreRulesGetAllParams struct {
	Vulnerability        string    `json:"vulnerability"`
	License              string    `json:"license"`
	Policy               string    `json:"policy"`
	Watch                string    `json:"watch"`
	ComponentName        string    `json:"component_name"`
	ComponentVersion     string    `json:"component_version"`
	ArtifactName         string    `json:"artifact_name"`
	ArtifactVersion      string    `json:"artifact_version"`
	BuildName            string    `json:"build_name"`
	BuildVersion         string    `json:"build_version"`
	ReleaseBundleName    string    `json:"release_bundle_name"`
	ReleaseBundleVersion string    `json:"release_bundle_version"`
	DockerLayer          string    `json:"docker_layer"`
	ExpiresBefore        time.Time `json:"expires_before"`
	ExpiresAfter         time.Time `json:"expires_after"`
	ProjectKey           string    `json:"project_key"`
	OrderBy              string    `json:"order_by"`
	Direction            string    `json:"direction"`
	PageNum              int       `json:"page_num"`
	NumOfRows            int       `json:"num_of_rows"`
}

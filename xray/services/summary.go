package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	servicesutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	summaryAPI = "api/v2/summary/"
)

func (ss *SummaryService) getSummaryUrl() string {
	return ss.XrayDetails.GetUrl() + summaryAPI
}

// SummaryService returns the https client and Xray details
type SummaryService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewSummaryService creates a new service to retrieve the version of Xray
func NewSummaryService(client *jfroghttpclient.JfrogHttpClient) *SummaryService {
	return &SummaryService{client: client}
}

func (ss *SummaryService) GetBuildSummary(params XrayBuildParams) (*SummaryResponse, error) {
	httpDetails := ss.XrayDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sbuild?build_name=%s&build_number=%s", ss.getSummaryUrl(), params.BuildName, params.BuildNumber)
	if params.Project != "" {
		url += "&" + projectKeyQueryParam + params.Project
	}
	resp, body, _, err := ss.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var summaryResponse SummaryResponse
	err = json.Unmarshal(body, &summaryResponse)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	if summaryResponse.Errors != nil && len(summaryResponse.Errors) > 0 {
		return nil, errorutils.CheckErrorf("Getting build-summary for build: %s failed with error: %s", summaryResponse.Errors[0].Identifier, summaryResponse.Errors[0].Error)
	}
	return &summaryResponse, nil
}

func (ss *SummaryService) GetArtifactSummary(params ArtifactSummaryParams) (*ArtifactSummaryResponse, error) {
	httpDetails := ss.XrayDetails.CreateHttpClientDetails()
	servicesutils.SetContentType("application/json", &httpDetails.Headers)

	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	url := fmt.Sprintf("%sartifact", ss.getSummaryUrl())
	resp, body, err := ss.client.SendPost(url, requestBody, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var response ArtifactSummaryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	if response.Errors != nil && len(response.Errors) > 0 {
		return nil, errorutils.CheckErrorf("Getting artifact-summary for artifact: %s failed with error: %s", response.Errors[0].Identifier, response.Errors[0].Error)
	}
	return &response, nil
}

type ArtifactSummaryParams struct {
	Checksums []string `json:"checksums,omitempty"`
	Paths     []string `json:"paths,omitempty"`
}

type ArtifactSummaryResponse struct {
	Artifacts []Artifact `json:"artifacts,omitempty"`
	Errors    []Error    `json:"errors,omitempty"`
}

type Artifact struct {
	General  General          `json:"general,omitempty"`
	Issues   []Issue          `json:"issues,omitempty"`
	Licenses []SummaryLicense `json:"licenses,omitempty"`
}

type General struct {
	ComponentId string `json:"component_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	PkgType     string `json:"pkg_type,omitempty"`
	Sha256      string `json:"sha256,omitempty"`
}

type SummaryResponse struct {
	Issues []Issue
	Errors []Error
}

type Issue struct {
	IssueId                string             `json:"issue_id,omitempty"`
	Summary                string             `json:"summary,omitempty"`
	Description            string             `json:"description,omitempty"`
	IssueType              string             `json:"issue_type,omitempty"`
	Severity               string             `json:"severity,omitempty"`
	Provider               string             `json:"provider,omitempty"`
	Cves                   []SummaryCve       `json:"cves,omitempty"`
	Created                string             `json:"created,omitempty"`
	ImpactPath             []string           `json:"impact_path,omitempty"`
	Components             []SummaryComponent `json:"components,omitempty"`
	ComponentPhysicalPaths []string           `json:"component_physical_paths,omitempty"`
}

type SummaryLicense struct {
	Components  []string `json:"components,omitempty"`
	FullName    string   `json:"full_name,omitempty"`
	MoreInfoUrl []string `json:"more_info_url,omitempty"`
	Name        string   `json:"name,omitempty"`
}

type Error struct {
	Error      string `json:"error,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type SummaryCve struct {
	Id          string   `json:"cve,omitempty"`
	CvssV2Score string   `json:"cvss_v2,omitempty"`
	CvssV3Score string   `json:"cvss_v3,omitempty"`
	Cwe         []string `json:"cwe,omitempty"`
}

type SummaryComponent struct {
	ComponentId   string   `json:"component_id,omitempty"`
	FixedVersions []string `json:"fixed_versions,omitempty"`
}

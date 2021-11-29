package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	summaryAPI = "api/v2/summary/"
)

func (ss *SummaryService) getSummeryUrl() string {
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
	url := fmt.Sprintf("%sbuild?build_name=%s&build_number=%s", ss.getSummeryUrl(), params.BuildName, params.BuildNumber)
	if params.Project != "" {
		url += "&" + projectKeyQueryParam + params.Project
	}
	resp, body, _, err := ss.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var summaryResponse SummaryResponse
	err = json.Unmarshal(body, &summaryResponse)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	if summaryResponse.Errors != nil && len(summaryResponse.Errors) > 0 {
		return nil, errorutils.CheckErrorf("Getting build-summery for build: %s failed with error: %s", summaryResponse.Errors[0].Identifier, summaryResponse.Errors[0].Error)
	}
	return &summaryResponse, nil
}

type SummaryResponse struct {
	Issues []Issue
	Errors []Error
}

type Issue struct {
	IssueId     string             `json:"issue_id,omitempty"`
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	IssueType   string             `json:"issue_type,omitempty"`
	Severity    string             `json:"severity,omitempty"`
	Provider    string             `json:"provider,omitempty"`
	Cves        []SummeryCve       `json:"cves,omitempty"`
	Created     string             `json:"created,omitempty"`
	ImpactPath  []string           `json:"impact_path,omitempty"`
	Components  []SummeryComponent `json:"components,omitempty"`
}

type Error struct {
	Error      string `json:"error,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type SummeryCve struct {
	Id          string `json:"cve,omitempty"`
	CvssV2Score string `json:"cvss_v2,omitempty"`
	CvssV3Score string `json:"cvss_v3,omitempty"`
}

type SummeryComponent struct {
	ComponentId   string   `json:"component_id,omitempty"`
	FixedVersions []string `json:"fixed_versions,omitempty"`
}

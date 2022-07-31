package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	// ReportsAPI refer to: https://www.jfrog.com/confluence/display/JFROG/Xray+REST+API#XrayRESTAPI-REPORTS
	ReportsAPI         = "api/v1/reports"
	VulnerabilitiesAPI = ReportsAPI + "/vulnerabilities"
)

// ReportService defines the Http client and Xray details
type ReportService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// ReportDetails defines the detail response for an Xray report
type ReportDetails struct {
	Id                 int    `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Type               string `json:"report_type,omitempty"`
	Status             string `json:"status,omitempty"`
	TotalArtifacts     int    `json:"total_artifacts,omitempty"`
	ProcessedArtifacts int    `json:"num_of_processed_artifacts,omitempty"`
	Progress           int    `json:"progress,omitempty"`
	RowCount           int    `json:"number_of_rows,omitempty"`
	StartTime          string `json:"start_time,omitempty"`
	EndTime            string `json:"end_time,omitempty"`
	Author             string `json:"author,omitempty"`
}

// ReportContentRequestParams defines a report content request
type ReportContentRequestParams struct {
	ReportId  string
	Direction string
	PageNum   int
	NumRows   int
	OrderBy   string
}

// ReportContent defines a report content response
type ReportContent struct {
	TotalRows int   `json:"total_rows"`
	Rows      []Row `json:"rows"`
}

// Row defines an entry of the report content
type Row struct {
	Cves                     []ReportCve `json:"cves,omitempty"`
	Cvsv2MaxScore            float64     `json:"cvss2_max_score,omitempty"`
	Cvsv3MaxScore            float64     `json:"cvss3_max_score,omitempty"`
	Summary                  string      `json:"summary,omitempty"`
	Severity                 string      `json:"severity,omitempty"`
	SeveritySource           string      `json:"severity_source,omitempty"`
	VulnerableComponent      string      `json:"vulnerable_component,omitempty"`
	ImpactedArtifact         string      `json:"impacted_artifact,omitempty"`
	ImpactPath               []string    `json:"impact_path,omitempty"`
	Path                     string      `json:"path,omitempty"`
	FixedVersions            []string    `json:"fixed_versions,omitempty"`
	Published                string      `json:"published,omitempty"`
	IssueId                  string      `json:"issue_id,omitempty"`
	PackageType              string      `json:"package_type,omitempty"`
	Provider                 string      `json:"provider,omitempty"`
	Description              string      `json:"description,omitempty"`
	References               []string    `json:"references,omitempty"`
	ExternalAdvisorySource   string      `json:"external_advisory_source,omitempty"`
	ExternalAdvisorySeverity string      `json:"external_advisory_severity,omitempty"`
}

type ReportCve struct {
	Id           string  `json:"cve,omitempty"`
	CvssV2Score  float64 `json:"cvss_v2_score,omitempty"`
	CvssV2Vector string  `json:"cvss_v2_vector,omitempty"`
	CvssV3Score  float64 `json:"cvss_v3_score,omitempty"`
	CvssV3Vector string  `json:"cvss_v3_vector,omitempty"`
}

// ReportRequestParams defines a report request
type ReportRequestParams struct {
	Name      string   `json:"name,omitempty"`
	Filters   Filter   `json:"filters,omitempty"`
	Resources Resource `json:"resources,omitempty"`
}

type Filter struct {
	HasRemediation *bool     `json:"has_remediation,omitempty"`
	CvssScore      CvssScore `json:"cvss_score,omitempty"`
	Severity       []string  `json:"severities,omitempty"`
}

type CvssScore struct {
	MinScore float32 `json:"min_score,omitempty"`
	MaxScore float32 `json:"max_score,omitempty"`
}

type Resource struct {
	IncludePathPatterns []string     `json:"include_path_patterns,omitempty"`
	Repositories        []Repository `json:"repositories,omitempty"`
}

type Repository struct {
	Name string `json:"name,omitempty"`
}

// ReportResponse defines a report request response
type ReportResponse struct {
	ReportId int    `json:"report_id"`
	Status   string `json:"status"`
}

// NewReportService creates a new Xray Report Service
func NewReportService(client *jfroghttpclient.JfrogHttpClient) *ReportService {
	return &ReportService{client: client}
}

// Vulnerabilities requests a new Xray scan for vulnerabilities
func (rs *ReportService) Vulnerabilities(req ReportRequestParams) (*ReportResponse, error) {
	retVal := ReportResponse{}
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s", rs.XrayDetails.GetUrl(), VulnerabilitiesAPI)
	content, err := json.Marshal(req)
	if err != nil {
		return &retVal, errorutils.CheckError(err)
	}

	resp, body, err := rs.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return &retVal, err
	}

	err = json.Unmarshal(body, &retVal)
	if err != nil {
		return &retVal, errorutils.CheckError(err)
	}

	return &retVal, nil
}

// Details retrieves the details for a report
func (rs *ReportService) Details(reportId string) (*ReportDetails, error) {
	retVal := ReportDetails{}
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s/%s", rs.XrayDetails.GetUrl(), ReportsAPI, reportId)
	resp, body, _, err := rs.client.SendGet(url, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return &retVal, err
	}

	err = json.Unmarshal(body, &retVal)
	if err != nil {
		return &retVal, errorutils.CheckError(err)
	}

	return &retVal, nil
}

// Content retrieves the report content for the provided request
func (rs *ReportService) Content(request ReportContentRequestParams) (*ReportContent, error) {
	retVal := ReportContent{}
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s/%s?direction=%s&page_num=%d&num_of_rows=%d&order_by=%s",
		rs.XrayDetails.GetUrl(), VulnerabilitiesAPI, request.ReportId, request.Direction, request.PageNum, request.NumRows, request.OrderBy)
	resp, body, err := rs.client.SendPost(url, nil, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return &retVal, err
	}

	err = json.Unmarshal(body, &retVal)
	return &retVal, errorutils.CheckError(err)
}

// Delete deletes the report that has an id matching reportId
func (rs *ReportService) Delete(reportId string) error {
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s/%s", rs.XrayDetails.GetUrl(), ReportsAPI, reportId)
	resp, body, err := rs.client.SendDelete(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	return nil
}

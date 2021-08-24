package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	reportsAPI         = "api/v1/reports"
	vulnerabilitiesAPI = reportsAPI + "/vulnerabilities"
)

// ReportService defines the Http client and XRay details
type ReportService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// ReportDetails defines the detail response for an XRay report
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
	Cves                []Cve    `json:"cves,omitempty"`
	Summary             string   `json:"summary,omitempty"`
	Severity            string   `json:"severity,omitempty"`
	VulnerableComponent string   `json:"vulnerable_component,omitempty"`
	ImpactedArtifact    string   `json:"impacted_artifact,omitempty"`
	Path                string   `json:"path,omitempty"`
	FixedVersions       []string `json:"fixed_versions,omitempty"`
	Published           string   `json:"published,omitempty"`
	IssueId             string   `json:"issue_id,omitempty"`
	PackageType         string   `json:"package_type,omitempty"`
	Provider            string   `json:"provider,omitempty"`
	Description         string   `json:"description,omitempty"`
	References          []string `json:"references,omitempty"`
}

// ReportRequestParams defines a report request
type ReportRequestParams struct {
	Name      string   `json:"name,omitempty"`
	Filters   Filter   `json:"filters,omitempty"`
	Resources Resource `json:"resources,omitempty"`
}

type Filter struct {
	HasRemediation bool      `json:"has_remediation,omitempty"`
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

// NewReportService creates a new XRay Report Service
func NewReportService(client *jfroghttpclient.JfrogHttpClient) *ReportService {
	return &ReportService{client: client}
}

// Vulnerabilities requests a new XRay scan for vulnerabilities
func (rs *ReportService) Vulnerabilities(req ReportRequestParams) (*ReportResponse, error) {
	retVal := ReportResponse{}
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s", rs.XrayDetails.GetUrl(), vulnerabilitiesAPI)
	content, err := json.Marshal(req)
	if err != nil {
		return &retVal, errorutils.CheckError(err)
	}

	resp, body, err := rs.client.SendPost(url, content, &httpClientsDetails)
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return &retVal, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &retVal)
	if err != nil {
		return &retVal, errors.New("failed unmarshalling report response")
	}

	return &retVal, nil
}

// Details retrieves the details for a report
func (rs *ReportService) Details(reportId string) (*ReportDetails, error) {
	retVal := ReportDetails{}
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s/%s", rs.XrayDetails.GetUrl(), reportsAPI, reportId)
	resp, body, _, err := rs.client.SendGet(url, true, &httpClientsDetails)
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return &retVal, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &retVal)
	if err != nil {
		return &retVal, errors.New("failed unmarshalling report details " + reportId)
	}

	return &retVal, nil
}

// Content retrieves the report content for the provided request
func (rs *ReportService) Content(request ReportContentRequestParams) (*ReportContent, error) {
	retVal := ReportContent{}
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s/%s?direction=%s&page_num=%d&num_of_rows=%d&order_by=%s",
		rs.XrayDetails.GetUrl(), vulnerabilitiesAPI, request.ReportId, request.Direction, request.PageNum, request.NumRows, request.OrderBy)
	resp, body, err := rs.client.SendPost(url, nil, &httpClientsDetails)
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return &retVal, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &retVal)
	if err != nil {
		return &retVal, errors.New("failed unmarshalling content for report " + request.ReportId)
	}
	return &retVal, nil
}

// Delete deletes the report that has an id matching reportId
func (rs *ReportService) Delete(reportId string) error {
	httpClientsDetails := rs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	url := fmt.Sprintf("%s/%s/%s", rs.XrayDetails.GetUrl(), reportsAPI, reportId)
	resp, body, err := rs.client.SendDelete(url, nil, &httpClientsDetails)
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}

	return nil
}

package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

const (
	violationsAPI = "api/v1/violations"
)

const (
	Applicable      ApplicabilityStatus = "applicable"
	NotApplicable   ApplicabilityStatus = "not_applicable"
	NotCovered      ApplicabilityStatus = "not_covered"
	Undetermined    ApplicabilityStatus = "undetermined"
	NotScanned      ApplicabilityStatus = "not_scanned"
	NotSupported    ApplicabilityStatus = "technology_unsupported"
	RescanRequired  ApplicabilityStatus = "rescan_required"
	UpgradeRequired ApplicabilityStatus = "upgrade_required"
)

type ApplicabilityStatus string

type ViolationsService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewViolationsService(client *jfroghttpclient.JfrogHttpClient) *ViolationsService {
	return &ViolationsService{client: client}
}

// Gets the Xray violations based on a set of search criteria: https://jfrog.com/help/r/xray-rest-apis/get-violations
func (vs *ViolationsService) GetViolations(params utils.ViolationsRequest) (response *ViolationsResponse, err error) {
	httpClientsDetails := vs.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	requestBody, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return
	}

	resp, body, err := vs.client.SendPost(clientutils.AppendScopedProjectKeyParam(vs.XrayDetails.GetUrl()+violationsAPI, vs.ScopeProjectKey), requestBody, &httpClientsDetails)
	if err != nil {
		return
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		err = errorutils.CheckErrorf("got unexpected server response while attempting to get violations:\n%s", err.Error())
		return nil, err
	}

	response = &ViolationsResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		return nil, errorutils.CheckErrorf("couldn't parse JFrog Xray server violations response: %s", err.Error())
	}
	return response, nil
}

type ViolationsResponse struct {
	Total      int             `json:"total_violations,omitempty"`
	Violations []XrayViolation `json:"violations,omitempty"`
}

type XrayViolation struct {
	// General violation properties
	IssueId  string              `json:"issue_id"`
	Type     utils.ViolationType `json:"type"`
	Watch    string              `json:"watch_name"`
	Severity utils.Severity      `json:"severity"`
	Created  string              `json:"created"`
	// Additional violation details (depending on IncludeDetails flag, the violation type and configuration)
	Id                   string              `json:"violation_id,omitempty"`
	Description          string              `json:"description,omitempty"`
	Summary              string              `json:"summary,omitempty"`
	Provider             string              `json:"provider,omitempty"`
	InfectedComponentIds []string            `json:"infected_components"`
	InfectedVersions     []string            `json:"infected_versions,omitempty"`
	InfectedFilePaths    []string            `json:"infected_file_path,omitempty"`
	PhysicalPaths        []string            `json:"component_physical_paths,omitempty"`
	Url                  string              `json:"violation_details_url,omitempty"`
	ImpactArtifacts      []string            `json:"impact_artifacts,omitempty"`
	GitRepository        string              `json:"impacted_git_repository,omitempty"`
	Policies             []ViolationPolicies `json:"matched_policies,omitempty"`
	// Optional Ignore information if exists
	IgnoreInfo *IgnoreRuleInfo `json:"ignore_rule_info,omitempty"`
	// Security Violations details (based on scan type)
	Cves                     []CveDetails              `json:"properties,omitempty"`
	FixVersions              []string                  `json:"fix_versions,omitempty"`
	JfrogResearchInformation *ExtendedInformation      `json:"extended_information,omitempty"`
	OperationalRisk          *OperationalRiskDetails   `json:"operational_risk,omitempty"`
	Applicability            []CveApplicability        `json:"applicability,omitempty"`
	ApplicabilityDetails     []CveApplicabilityDetails `json:"applicability_details,omitempty"`
	ExposureDetails          *ExposureDetails          `json:"details,omitempty"`
	SastDetails              *BaseJasDetails           `json:"sast_details,omitempty"`
}

type ViolationPolicies struct {
	Policy     string `json:"policy"`
	Rule       string `json:"rule"`
	IsBlocking bool   `json:"is_blocking"`

	IsIgnored    bool `json:"is_ignored,omitempty"`
	BlockingMask int  `json:"blocking_mask,omitempty"`
}

type IgnoreRuleInfo struct {
	Id             string `json:"id"`
	Type           string `json:"ignore_rule_type"`
	Author         string `json:"author"`
	Created        string `json:"created"`
	Notes          string `json:"notes"`
	IsExpired      bool   `json:"is_expired"`
	ExpirationDate string `json:"expires_at,omitempty"`
	DeletedBy      string `json:"deleted_by,omitempty"`
	DeletedAt      string `json:"deleted_at,omitempty"`
}

type OperationalRiskDetails struct {
	Risk          string   `json:"risk"`
	RiskReason    string   `json:"risk_reason"`
	IsEol         *bool    `json:"is_eol,omitempty"`
	EolMessage    string   `json:"eol_message,omitempty"`
	LatestVersion string   `json:"latest_version,omitempty"`
	NewerVersions *int     `json:"newer_versions,omitempty"`
	Cadence       *float64 `json:"cadence"`
	Commits       *int64   `json:"commits,omitempty"`
	Committers    *int     `json:"committers,omitempty"`
	Released      string   `json:"released,omitempty"`
}

type CveDetails struct {
	Id           string         `json:"cve,omitempty"`
	CvssV2Vector string         `json:"cvss_v2,omitempty"`
	CvssV3Vector string         `json:"cvss_v3,omitempty"`
	Cwe          []string       `json:"cwe,omitempty"`
	CweDetails   map[string]Cwe `json:"cwe_details,omitempty"`
}

type CveApplicability struct {
	ScannerAvailable      bool                       `json:"scanner_available"`
	ComponentId           string                     `json:"component_id"`
	VulnerableComponentId string                     `json:"source_comp_id"`
	CveId                 string                     `json:"cve_id"`
	ScanStatus            int                        `json:"scan_status"`
	Applicability         *bool                      `json:"applicability,omitempty"`
	ScannerDescription    string                     `json:"scanner_explanation,omitempty"`
	Reason                *string                    `json:"info,omitempty"`
	Evidence              []CveApplicabilityEvidence `json:"evidence,omitempty"`
}

// Information is returned dynamically based values in columns and not by attributes
// Values: [Path (relative path), Evidence (snippet), Line Number (starts from 1) and issue Found (Reason)]
type CveApplicabilityEvidence struct {
	ColumnNames []string   `json:"column_names"`
	Rows        [][]string `json:"rows"`
}

type CveApplicabilityDetails struct {
	ComponentId           string              `json:"component_id"`
	VulnerableComponentId string              `json:"source_comp_id"`
	CveId                 string              `json:"vulnerability_id"`
	Status                ApplicabilityStatus `json:"result"`
}

type BaseJasDetails struct {
	Id           string         `json:"id"`
	Status       string         `json:"status"`
	Severity     utils.Severity `json:"jfrog_severity,omitempty"`
	Reason       string         `json:"description,omitempty"`
	Abbreviation string         `json:"abbreviation"`
	CWE          *JasCwe        `json:"cwe,omitempty"`
	Findings     JasFindings    `json:"findings,omitempty"`
}

type JasCwe struct {
	Id          string `json:"cwe_id"`
	Description string `json:"cwe_name"`
	Link        string `json:"cwe_link"`
}

type JasFindings struct {
	ScannerDescription string `json:"explanation,omitempty"`
}

type ExposureDetails struct {
	BaseJasDetails
	Origin string `json:"origin,omitempty"`
}

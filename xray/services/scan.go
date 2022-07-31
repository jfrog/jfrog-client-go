package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const (
	scanGraphAPI = "api/v1/scan/graph"

	// Graph scan query params
	repoPathQueryParam = "repo_path="
	projectQueryParam  = "project="
	watchesQueryParam  = "watch="
	scanTypeQueryParam = "scan_type="

	// Get scan results query params
	includeVulnerabilitiesParam = "?include_vulnerabilities=true"
	includeLicensesParam        = "?include_licenses=true"
	andIncludeLicensesParam     = "&include_licenses=true"

	// Get scan results timeouts
	defaultMaxWaitMinutes    = 45 * time.Minute // 45 minutes
	defaultSyncSleepInterval = 5 * time.Second  // 5 seconds

	// ScanType values
	Dependency ScanType = "dependency"
	Binary     ScanType = "binary"

	xrayScanStatusFailed = "failed"
)

type ScanType string

type ScanService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewScanService creates a new service to scan binaries and audit code projects' dependencies.
func NewScanService(client *jfroghttpclient.JfrogHttpClient) *ScanService {
	return &ScanService{client: client}
}

func createScanGraphQueryParams(scanParams XrayGraphScanParams) string {
	var params []string
	if scanParams.ProjectKey != "" {
		params = append(params, projectQueryParam+scanParams.ProjectKey)
	} else if scanParams.RepoPath != "" {
		params = append(params, repoPathQueryParam+scanParams.RepoPath)
	} else if len(scanParams.Watches) > 0 {
		for _, watch := range scanParams.Watches {
			if watch != "" {
				params = append(params, watchesQueryParam+watch)
			}
		}
	}

	if scanParams.ScanType != "" {
		params = append(params, scanTypeQueryParam+string(scanParams.ScanType))
	}

	if len(params) == 0 {
		return ""
	}
	return "?" + strings.Join(params, "&")
}

func (ss *ScanService) ScanGraph(scanParams XrayGraphScanParams) (string, error) {
	httpClientsDetails := ss.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	requestBody, err := json.Marshal(scanParams.Graph)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	url := ss.XrayDetails.GetUrl() + scanGraphAPI
	url += createScanGraphQueryParams(scanParams)
	resp, body, err := ss.client.SendPost(url, requestBody, &httpClientsDetails)
	if err != nil {
		return "", err
	}

	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		scanErrorJson := ScanErrorJson{}
		if e := json.Unmarshal(body, &scanErrorJson); e == nil {
			return "", errorutils.CheckErrorf(scanErrorJson.Error)
		}
		return "", err
	}
	scanResponse := RequestScanResponse{}
	if err = json.Unmarshal(body, &scanResponse); err != nil {
		return "", errorutils.CheckError(err)
	}
	return scanResponse.ScanId, err
}

func (ss *ScanService) GetScanGraphResults(scanId string, includeVulnerabilities, includeLicenses bool) (*ScanResponse, error) {
	httpClientsDetails := ss.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	//The scan request may take some time to complete. We expect to receive a 202 response, until the completion.
	endPoint := ss.XrayDetails.GetUrl() + scanGraphAPI + "/" + scanId
	if includeVulnerabilities {
		endPoint += includeVulnerabilitiesParam
		if includeLicenses {
			endPoint += andIncludeLicensesParam
		}
	} else if includeLicenses {
		endPoint += includeLicensesParam
	}
	log.Info("Waiting for scan to complete...")
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		resp, body, _, err := ss.client.SendGet(endPoint, true, &httpClientsDetails)
		if err != nil {
			return true, nil, err
		}
		if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
			return true, nil, err
		}
		// Got the full valid response.
		if resp.StatusCode == http.StatusOK {
			return true, body, nil
		}
		return false, nil, nil
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get Dependencies Scan results... ",
	}

	body, err := pollingExecutor.Execute()
	if err != nil {
		return nil, err
	}
	scanResponse := ScanResponse{}
	if err = json.Unmarshal(body, &scanResponse); err != nil {
		return nil, errorutils.CheckError(err)
	}
	if scanResponse.ScannedStatus == xrayScanStatusFailed {
		return nil, errorutils.CheckErrorf("Xray scan failed")
	}
	return &scanResponse, err
}

type XrayGraphScanParams struct {
	// A path in Artifactory that this Artifact is intended to be deployed to.
	// This will provide a way to extract the watches that should be applied on this graph
	RepoPath               string
	ProjectKey             string
	Watches                []string
	ScanType               ScanType
	Graph                  *GraphNode
	IncludeVulnerabilities bool
	IncludeLicenses        bool
}

type GraphNode struct {
	// Component Id in the JFrog standard.
	// For instance, for maven: gav://<groupId>:<artifactId>:<version>
	// For detailed format examples please see:
	// https://www.jfrog.com/confluence/display/JFROG/Xray+REST+API#XrayRESTAPI-ComponentIdentifiers
	Id string `json:"component_id,omitempty"`
	// Sha of the binary representing the component.
	Sha256 string `json:"sha256,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
	// For root file shall be the file name.
	// For internal components shall be the internal path. (Relevant only for binary scan).
	Path string `json:"path,omitempty"`
	// List of license names
	Licenses []string `json:"licenses,omitempty"`
	// Component properties
	Properties map[string]string `json:"properties,omitempty"`
	// List of subcomponents.
	Nodes []*GraphNode `json:"nodes,omitempty"`
	// Node parent (for internal use)
	Parent *GraphNode `json:"-"`
}

type RequestScanResponse struct {
	ScanId string `json:"scan_id,omitempty"`
}

type ScanErrorJson struct {
	Error string `json:"error"`
}

type ScanResponse struct {
	ScanId             string          `json:"scan_id,omitempty"`
	XrayDataUrl        string          `json:"xray_data_url,omitempty"`
	Violations         []Violation     `json:"violations,omitempty"`
	Vulnerabilities    []Vulnerability `json:"vulnerabilities,omitempty"`
	Licenses           []License       `json:"licenses,omitempty"`
	ScannedComponentId string          `json:"component_id,omitempty"`
	ScannedPackageType string          `json:"package_type,omitempty"`
	ScannedStatus      string          `json:"status,omitempty"`
}

type Violation struct {
	Summary             string               `json:"summary,omitempty"`
	Severity            string               `json:"severity,omitempty"`
	ViolationType       string               `json:"type,omitempty"`
	Components          map[string]Component `json:"components,omitempty"`
	WatchName           string               `json:"watch_name,omitempty"`
	IssueId             string               `json:"issue_id,omitempty"`
	Cves                []Cve                `json:"cves,omitempty"`
	References          []string             `json:"references,omitempty"`
	FailBuild           bool                 `json:"fail_build,omitempty"`
	LicenseKey          string               `json:"license_key,omitempty"`
	LicenseName         string               `json:"license_name,omitempty"`
	IgnoreUrl           string               `json:"ignore_url,omitempty"`
	RiskReason          string               `json:"risk_reason,omitempty"`
	IsEol               *bool                `json:"is_eol,omitempty"`
	EolMessage          string               `json:"eol_message,omitempty"`
	LatestVersion       string               `json:"latest_version,omitempty"`
	NewerVersions       *int                 `json:"newer_versions,omitempty"`
	Cadence             *float64             `json:"cadence,omitempty"`
	Commits             *int64               `json:"commits,omitempty"`
	Committers          *int                 `json:"committers,omitempty"`
	ExtendedInformation *ExtendedInformation `json:"extended_information,omitempty"`
	Technology          string               `json:"-"`
}

type Vulnerability struct {
	Cves                []Cve                `json:"cves,omitempty"`
	Summary             string               `json:"summary,omitempty"`
	Severity            string               `json:"severity,omitempty"`
	Components          map[string]Component `json:"components,omitempty"`
	IssueId             string               `json:"issue_id,omitempty"`
	References          []string             `json:"references,omitempty"`
	ExtendedInformation *ExtendedInformation `json:"extended_information,omitempty"`
	Technology          string               `json:"-"`
}

type License struct {
	Key        string               `json:"license_key,omitempty"`
	Name       string               `json:"name,omitempty"`
	Components map[string]Component `json:"components,omitempty"`
	Custom     bool                 `json:"custom,omitempty"`
	References []string             `json:"references,omitempty"`
}

type Component struct {
	FixedVersions []string           `json:"fixed_versions,omitempty"`
	ImpactPaths   [][]ImpactPathNode `json:"impact_paths,omitempty"`
}

type ImpactPathNode struct {
	ComponentId string `json:"component_id,omitempty"`
	FullPath    string `json:"full_path,omitempty"`
}

type Cve struct {
	Id           string `json:"cve,omitempty"`
	CvssV2Score  string `json:"cvss_v2_score,omitempty"`
	CvssV2Vector string `json:"cvss_v2_vector,omitempty"`
	CvssV3Score  string `json:"cvss_v3_score,omitempty"`
	CvssV3Vector string `json:"cvss_v3_vector,omitempty"`
}

type ExtendedInformation struct {
	ShortDescription             string                        `json:"short_description,omitempty"`
	FullDescription              string                        `json:"full_description,omitempty"`
	JfrogResearchSeverity        string                        `json:"jfrog_research_severity,omitempty"`
	JfrogResearchSeverityReasons []JfrogResearchSeverityReason `json:"jfrog_research_severity_reasons,omitempty"`
	Remediation                  string                        `json:"remediation,omitempty"`
}

type JfrogResearchSeverityReason struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsPositive  bool   `json:"is_positive,omitempty"`
}

func (gp *XrayGraphScanParams) GetProjectKey() string {
	return gp.ProjectKey
}

func (currNode *GraphNode) NodeHasLoop() bool {
	parent := currNode.Parent
	for parent != nil {
		if currNode.Id == parent.Id {
			return true
		}
		parent = parent.Parent
	}
	return false
}

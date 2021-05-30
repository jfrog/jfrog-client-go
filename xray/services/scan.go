package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	scanGraphAPI             = "api/v1/scan/graph"
	repoPathQueryParam       = "?repo_path="
	projectQueryParam        = "?project="
	defaultMaxWaitMinutes    = 5 * time.Minute // 5 minutes
	defaultSyncSleepInterval = 5               // 5 seconds
)

type ScanService struct {
	client         *jfroghttpclient.JfrogHttpClient
	XrayDetails    auth.ServiceDetails
	MaxWaitMinutes time.Duration
}

// NewScanService creates a new service to scan Binaries and Sources.
func NewScanService(client *jfroghttpclient.JfrogHttpClient) *ScanService {
	return &ScanService{client: client}
}

func (ss *ScanService) ScanGraph(scanParams XrayGraphScanParams) (string, error) {
	httpClientsDetails := ss.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	requestBody, err := json.Marshal(scanParams.Graph)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	url := ss.XrayDetails.GetUrl() + scanGraphAPI
	if scanParams.ProjectKey != "" {
		url += projectQueryParam + scanParams.ProjectKey
	} else if scanParams.RepoPath != "" {
		url += repoPathQueryParam + scanParams.RepoPath
	}
	resp, body, err := ss.client.SendPost(url, requestBody, &httpClientsDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", errorutils.CheckError(errors.New("Server Response: " + resp.Status))
	}
	scanResponse := RequestScanResponse{}
	json.Unmarshal(body, &scanResponse)
	return scanResponse.ScanId, err
}

func (ss *ScanService) GetScanGraphResults(scanId string) (*ScanResponse, error) {
	maxWaitMinutes := defaultMaxWaitMinutes
	if ss.MaxWaitMinutes > 0 {
		maxWaitMinutes = ss.MaxWaitMinutes
	}
	httpClientsDetails := ss.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	message := fmt.Sprintf("Sync: Get Scan Graph Results. Scan ID:%s...", scanId)
	// The scan request may take time to be calculated, we expect to receive a 202 response until the compleation of the requested scan.
	ticker := time.NewTicker(defaultSyncSleepInterval)
	timeout := make(chan bool)
	errChan := make(chan error)
	resultChan := make(chan []byte)
	endPoint := ss.XrayDetails.GetUrl() + scanGraphAPI + "/" + scanId
	go func() {
		for {
			select {
			case <-timeout:
				errChan <- errorutils.CheckError(errors.New("Timeout for sync get scan graph results."))
				resultChan <- nil
				return
			case _ = <-ticker.C:
				log.Debug(message)
				resp, body, _, err := ss.client.SendGet(endPoint, true, &httpClientsDetails)
				if err != nil {
					errChan <- err
					resultChan <- nil
					return
				}
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
					errChan <- errorutils.CheckError(errors.New("Server Response: " + resp.Status))
					resultChan <- nil
					return
				}
				// Got the full valid response.
				if resp.StatusCode == http.StatusOK {
					errChan <- nil
					resultChan <- body
					return
				}
			}
		}
	}()
	// Make sure we don't wait forever
	go func() {
		time.Sleep(maxWaitMinutes)
		timeout <- true
	}()
	// Wait for result or error
	err := <-errChan
	body := <-resultChan
	ticker.Stop()
	if err != nil {
		return nil, err
	}
	scanResponse := ScanResponse{}
	json.Unmarshal(body, &scanResponse)
	return &scanResponse, err
}

type XrayGraphScanParams struct {
	// A path in Artifactory that this Artifact is intended to be deployed to.
	// This will provide a way to extract the watches that should be applied on this graph
	RepoPath   string
	ProjectKey string
	Graph      *GraphNode
}

type GraphNode struct {
	// Component Id in the JFrog standard.
	Id string `json:"id,omitempty"`
	// Sha of the binary representing the component.
	Sha256 string `json:"sha256,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
	// For root file shall be the file name.
	// For internal components shall be the internal path. (Relevant only for binary scan).
	Path string `json:"path,omitempty"`
	// List of license name
	Licenses []string `json:"licenses,omitempty"`
	// List of sub components.
	Nodes []*GraphNode `json:"nodes,omitempty"`
}

type RequestScanResponse struct {
	ScanId string `json:"scan_id,omitempty"`
}

type ScanResponse struct {
	ScanId          string          `json:"scan_id,omitempty"`
	Violations      []Violation     `json:"violations,omitempty"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	Licenses        []License       `json:"licenses,omitempty"`
}

type Violation struct {
	Summary       string               `json:"summary,omitempty"`
	Severity      string               `json:"severity,omitempty"`
	ViolationType string               `json:"type,omitempty"`
	Components    map[string]Component `json:"components,omitempty"`
	WatchName     string               `json:"watch_name,omitempty"`
	IssueId       string               `json:"issue_id,omitempty"`
	Cves          []Cve                `json:"cves,omitempty"`
	References    []string             `json:"references,omitempty"`
	FailBuild     bool                 `json:"fail_build,omitempty"`
	LicenseKey    string               `json:"license_key,omitempty"`
	LicenseName   string               `json:"license_name,omitempty"`
	IgnoreUrl     string               `json:"ignore_url,omitempty"`
}

type Vulnerability struct {
	Cves                 []Cve                `json:"cves,omitempty"`
	Summary              string               `json:"summary,omitempty"`
	Severity             string               `json:"severity,omitempty"`
	VulnerableComponents []string             `json:"vulnerable_components,omitempty"`
	Components           map[string]Component `json:"components,omitempty"`
}

type License struct {
	Key        string               `json:"key,omitempty"`
	Name       string               `json:"name,omitempty"`
	Components map[string]Component `json:"components,omitempty"`
	Custom     bool                 `json:"custom,omitempty"`
	References []string             `json:"references,omitempty"`
}
type Component struct {
	FixedVersions []string `json:"fixed_versions,omitempty"`
	ImpactPaths   []string `json:"impact_paths,omitempty"`
}
type Cve struct {
	Id           string `json:"cve,omitempty"`
	CvssV2Score  string `json:"cvss_v2_score,omitempty"`
	CvssV2Vector string `json:"cvss_v2_vector,omitempty"`
	CvssV3Score  string `json:"cvss_v3_score,omitempty"`
	CvssV3Vector string `json:"cvss_v3_vector,omitempty"`
}

func (gp *XrayGraphScanParams) GetProjectKey() string {
	return gp.ProjectKey
}

func NewXrayGraphScanParams() XrayGraphScanParams {
	return XrayGraphScanParams{}
}

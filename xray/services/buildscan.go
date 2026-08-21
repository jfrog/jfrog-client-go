package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jfrog/gofrog/version"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	BuildScanAPI                        = "api/v2/ci/build"
	xrayScanBuildNotSelectedForIndexing = "is not selected for indexing"

	XrayScanBuildNoFailBuildPolicy = "No Xray “Fail build in case of a violation” policy rule has been defined on this build"
	XrayScanBuildNotFoundFormat    = "Build %s number %d wasn't found in Artifactory"

	projectKeyQueryParam                  = "projectKey="
	includeVulnerabilitiesQueryParam      = "include_vulnerabilities="
	buildScanResultsPostApiMinXrayVersion = "3.77.0"
	buildScanResultsPostApi               = "scanResult"
)

var (
	// Phrases Xray returns (in the "info" or "error" field of a 404) while the build has not finished
	// asynchronous indexing yet — a transient, retryable condition. Matched case-insensitively because the
	// exact wording has varied across Xray versions, so we match a set of stable substrings rather than one
	// exact format:
	//   - "not indexed" covers the issue-596 string ("build doesn't exist or not indexed in Xray").
	//   - "build doesn't exist" is defensive coverage for phrasings that omit "not indexed".
	//   - "wasn't found in artifactory" covers the legacy "Build <name> number <num> wasn't found in
	//     Artifactory" message — the case-insensitive substring subsumes that older structured wording.
	buildNotIndexedRegex = regexp.MustCompile(`(?i)not indexed|build doesn't exist|wasn't found in artifactory`)
)

type BuildScanService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

// NewBuildScanService creates a new service to scan build dependencies.
func NewBuildScanService(client *jfroghttpclient.JfrogHttpClient) *BuildScanService {
	return &BuildScanService{client: client}
}

func (bs *BuildScanService) ScanBuild(params XrayBuildParams, includeVulnerabilities bool, triggerRetries int) (scanResponse *BuildScanResponse, noFailBuildPolicy bool, err error) {
	if err = bs.triggerScan(params, triggerRetries); err != nil {
		// If the includeVulnerabilities flag is true and error is "No Xray Fail build...." continue to getBuildScanResults to get vulnerabilities
		if includeVulnerabilities && strings.Contains(err.Error(), XrayScanBuildNoFailBuildPolicy) {
			noFailBuildPolicy = true
		} else {
			return
		}
	}
	getResultsReqFunc, err := bs.prepareGetResultsRequest(params, includeVulnerabilities)
	if err != nil {
		return
	}
	scanResponse, err = bs.getBuildScanResults(getResultsReqFunc, params)
	return
}

func isArtifactoryBuildNotFoundError(resp *http.Response, body []byte) error {
	if resp.StatusCode != http.StatusNotFound {
		return nil
	}
	buildScanResponse := RequestBuildScanResponse{}
	if err := json.Unmarshal(body, &buildScanResponse); err != nil {
		// Unable to parse response body = actual 404 error.
		log.Debug("Failed to parse Xray build scan response:", err)
		return nil
	}
	// Xray returns the transient "still indexing" message in either the "info" or the "error" field,
	// and the exact wording has changed across Xray versions. Treat any known phrasing as a retryable
	// not-found so the trigger waits for asynchronous indexing to complete instead of failing.
	if isBuildNotIndexedMessage(buildScanResponse.Info) {
		return errors.New(buildScanResponse.Info)
	}
	if isBuildNotIndexedMessage(buildScanResponse.Error) {
		return errors.New(buildScanResponse.Error)
	}
	return nil
}

// isBuildNotIndexedMessage reports whether an Xray 404 message indicates the build is not yet indexed
// (a transient, retryable condition) rather than a permanent error such as a genuinely wrong build name.
func isBuildNotIndexedMessage(message string) bool {
	return buildNotIndexedRegex.MatchString(message)
}

func (bs *BuildScanService) triggerScan(params XrayBuildParams, retries int) error {
	paramsBytes, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := bs.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()
	url := bs.XrayDetails.GetUrl() + BuildScanAPI

	if retries <= 0 {
		retries = 1
	}
	retryExecutor := utils.RetryExecutor{
		MaxRetries:               retries,
		RetriesIntervalMilliSecs: int(defaultSyncSleepInterval.Milliseconds()),
		LogMsgPrefix:             "trigger build scan ",
		ExecutionHandler: func() (shouldRetry bool, err error) {
			resp, body, err := bs.client.SendPost(utils.AppendScopedProjectKeyParam(url, bs.ScopeProjectKey), paramsBytes, &httpClientsDetails)
			if err != nil {
				return false, err
			}
			if err = isArtifactoryBuildNotFoundError(resp, body); err != nil {
				return true, err
			}
			if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
				return false, err
			}
			buildScanResponse := RequestBuildScanResponse{}
			if err = json.Unmarshal(body, &buildScanResponse); err != nil {
				return false, errorutils.CheckError(err)
			}
			buildScanInfo := buildScanResponse.Info
			if strings.Contains(buildScanInfo, xrayScanBuildNotSelectedForIndexing) ||
				strings.Contains(buildScanInfo, XrayScanBuildNoFailBuildPolicy) {
				return false, errors.New(buildScanResponse.Info)
			}
			log.Info(buildScanInfo)
			return false, nil
		},
	}
	return retryExecutor.Execute()
}

// prepareGetResultsRequest creates a function that requests for the scan results from Xray.
// Starting from Xray version 3.77.0, there's a new POST API that supports special characters in the build-name and build-number fields.
func (bs *BuildScanService) prepareGetResultsRequest(params XrayBuildParams, includeVulnerabilities bool) (getResultsReqFunc func() (*http.Response, []byte, error), err error) {
	paramsBytes, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return
	}
	xrayVer, err := bs.XrayDetails.GetVersion()
	if err != nil {
		return
	}
	var queryParams []string
	if includeVulnerabilities {
		queryParams = append(queryParams, includeVulnerabilitiesQueryParam+"true")
	}
	httpClientsDetails := bs.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()
	if version.NewVersion(xrayVer).AtLeast(buildScanResultsPostApiMinXrayVersion) {
		getResultsReqFunc = bs.getResultsPostRequestFunc(params, paramsBytes, &httpClientsDetails, queryParams)
		return
	}
	getResultsReqFunc = bs.getResultsGetRequestFunc(params, &httpClientsDetails, queryParams)
	return
}

func (bs *BuildScanService) getBuildScanResults(reqFunc func() (*http.Response, []byte, error), params XrayBuildParams) (*BuildScanResponse, error) {
	log.Info("Waiting for Build Scan to complete...")
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		resp, body, err := reqFunc()
		if err != nil {
			return true, nil, err
		}
		if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
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
		MsgPrefix:       fmt.Sprintf("Get Build Scan results for Build: %s/%s...", params.BuildName, params.BuildNumber),
	}

	body, err := pollingExecutor.Execute()
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errorutils.CheckErrorf(
			"Received empty response from Xray server (HTTP 200). " +
				"This may indicate a server-side timeout during authentication. Please retry.",
		)
	}
	buildScanResponse := BuildScanResponse{}
	if err = json.Unmarshal(body, &buildScanResponse); err != nil {
		return nil, errorutils.CheckError(err)
	}
	if buildScanResponse.Status == xrayScanStatusFailed {
		return nil, errorutils.CheckErrorf("Xray build scan failed")
	}
	return &buildScanResponse, err
}

func (bs *BuildScanService) getResultsGetRequestFunc(params XrayBuildParams, httpClientsDetails *httputils.HttpClientDetails, queryParams []string) func() (*http.Response, []byte, error) {
	endPoint := fmt.Sprintf("%s%s/%s/%s", bs.XrayDetails.GetUrl(), BuildScanAPI, params.BuildName, params.BuildNumber)
	if params.Project != "" {
		queryParams = append(queryParams, projectKeyQueryParam+params.Project)
	}
	if len(queryParams) > 0 {
		endPoint += "?" + strings.Join(queryParams, "&")
	}
	return func() (*http.Response, []byte, error) {
		resp, body, _, err := bs.client.SendGet(utils.AppendScopedProjectKeyParam(endPoint, bs.ScopeProjectKey), true, httpClientsDetails)
		return resp, body, err
	}
}

func (bs *BuildScanService) getResultsPostRequestFunc(params XrayBuildParams, paramsBytes []byte, httpClientsDetails *httputils.HttpClientDetails, queryParams []string) func() (*http.Response, []byte, error) {
	endPoint := fmt.Sprintf("%s%s/%s", bs.XrayDetails.GetUrl(), BuildScanAPI, buildScanResultsPostApi)
	if params.Project != "" {
		queryParams = append(queryParams, projectKeyQueryParam+params.Project)
	}
	if len(queryParams) > 0 {
		endPoint += "?" + strings.Join(queryParams, "&")
	}
	return func() (*http.Response, []byte, error) {
		resp, body, err := bs.client.SendPost(utils.AppendScopedProjectKeyParam(endPoint, bs.ScopeProjectKey), paramsBytes, httpClientsDetails)
		return resp, body, err
	}
}

type XrayBuildParams struct {
	BuildName   string `json:"build_name,omitempty"`
	BuildNumber string `json:"build_number,omitempty"`
	Project     string `json:"project,omitempty"`
	Rescan      bool   `json:"rescan,omitempty"`
}

type RequestBuildScanResponse struct {
	Info  string `json:"info,omitempty"`
	Error string `json:"error,omitempty"`
}

type BuildScanResponse struct {
	Status          string          `json:"status,omitempty"`
	MoreDetailsUrl  string          `json:"more_details_url,omitempty"`
	FailBuild       bool            `json:"fail_build,omitempty"`
	Violations      []Violation     `json:"violations,omitempty"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	Info            string          `json:"info,omitempty"`
}

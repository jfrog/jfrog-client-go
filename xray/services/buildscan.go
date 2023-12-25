package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/gofrog/version"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
)

const (
	BuildScanAPI                          = "api/v2/ci/build"
	xrayScanBuildNotSelectedForIndexing   = "is not selected for indexing"
	XrayScanBuildNoFailBuildPolicy        = "No Xray “Fail build in case of a violation” policy rule has been defined on this build"
	projectKeyQueryParam                  = "projectKey="
	includeVulnerabilitiesQueryParam      = "include_vulnerabilities="
	buildScanResultsPostApiMinXrayVersion = "3.77.0"
	buildScanResultsPostApi               = "scanResult"
)

type BuildScanService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewBuildScanService creates a new service to scan build dependencies.
func NewBuildScanService(client *jfroghttpclient.JfrogHttpClient) *BuildScanService {
	return &BuildScanService{client: client}
}

func (bs *BuildScanService) ScanBuild(params XrayBuildParams, includeVulnerabilities bool) (scanResponse *BuildScanResponse, noFailBuildPolicy bool, err error) {
	paramsBytes, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return
	}
	err = bs.triggerScan(paramsBytes)
	if err != nil {
		// If the includeVulnerabilities flag is true and error is "No Xray Fail build...." continue to getBuildScanResults to get vulnerabilities
		if includeVulnerabilities && strings.Contains(err.Error(), XrayScanBuildNoFailBuildPolicy) {
			noFailBuildPolicy = true
		} else {
			return
		}
	}
	getResultsReqFunc, err := bs.prepareGetResultsRequest(params, paramsBytes, includeVulnerabilities)
	if err != nil {
		return
	}
	scanResponse, err = bs.getBuildScanResults(getResultsReqFunc, params)
	return
}

func (bs *BuildScanService) triggerScan(paramsBytes []byte) error {
	httpClientsDetails := bs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	url := bs.XrayDetails.GetUrl() + BuildScanAPI

	resp, body, err := bs.client.SendPost(url, paramsBytes, &httpClientsDetails)
	if err != nil {
		return err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		return err
	}
	buildScanResponse := RequestBuildScanResponse{}
	if err = json.Unmarshal(body, &buildScanResponse); err != nil {
		return errorutils.CheckError(err)
	}
	buildScanInfo := buildScanResponse.Info
	if strings.Contains(buildScanInfo, xrayScanBuildNotSelectedForIndexing) ||
		strings.Contains(buildScanInfo, XrayScanBuildNoFailBuildPolicy) {
		return errors.New(buildScanResponse.Info)
	}
	log.Info(buildScanInfo)
	return nil
}

// prepareGetResultsRequest creates a function that requests for the scan results from Xray.
// Starting from Xray version 3.77.0, there's a new POST API that supports special characters in the build-name and build-number fields.
func (bs *BuildScanService) prepareGetResultsRequest(params XrayBuildParams, paramsBytes []byte, includeVulnerabilities bool) (getResultsReqFunc func() (*http.Response, []byte, error), err error) {
	xrayVer, err := bs.XrayDetails.GetVersion()
	if err != nil {
		return
	}
	var queryParams []string
	if includeVulnerabilities {
		queryParams = append(queryParams, includeVulnerabilitiesQueryParam+"true")
	}
	httpClientsDetails := bs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
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
		resp, body, _, err := bs.client.SendGet(endPoint, true, httpClientsDetails)
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
		resp, body, err := bs.client.SendPost(endPoint, paramsBytes, httpClientsDetails)
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
	Info string `json:"info,omitempty"`
}

type BuildScanResponse struct {
	Status          string          `json:"status,omitempty"`
	MoreDetailsUrl  string          `json:"more_details_url,omitempty"`
	FailBuild       bool            `json:"fail_build,omitempty"`
	Violations      []Violation     `json:"violations,omitempty"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	Info            string          `json:"info,omitempty"`
}

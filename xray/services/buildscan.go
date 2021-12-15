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
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
)

const (
	buildScanAPI                        = "api/v2/ci/build"
	xrayScanBuildNotSelectedForIndexing = "is not selected for indexing"
	XrayScanBuildNoFailBuildPolicy      = "No Xray “Fail build in case of a violation” policy rule has been defined on this build"
	projectKeyQueryParam                = "projectKey="
)

type BuildScanService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewBuildScanService creates a new service to scan build dependencies.
func NewBuildScanService(client *jfroghttpclient.JfrogHttpClient) *BuildScanService {
	return &BuildScanService{client: client}
}

func (bs *BuildScanService) Scan(params XrayBuildParams) error {
	httpClientsDetails := bs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	requestBody, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := bs.XrayDetails.GetUrl() + buildScanAPI

	resp, body, err := bs.client.SendPost(url, requestBody, &httpClientsDetails)
	if err != nil {
		return err
	}

	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
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

func (bs *BuildScanService) GetBuildScanResults(params XrayBuildParams) (*BuildScanResponse, error) {
	endPoint := fmt.Sprintf("%s%s/%s/%s", bs.XrayDetails.GetUrl(), buildScanAPI, params.BuildName, params.BuildNumber)
	if params.Project != "" {
		endPoint += "?" + projectKeyQueryParam + params.Project
	}
	syncMessage := fmt.Sprintf("Sync: Get Build Scan results. Build:%s/%s...", params.BuildName, params.BuildNumber)
	httpClientsDetails := bs.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		log.Debug(syncMessage)
		resp, body, _, err := bs.client.SendGet(endPoint, true, &httpClientsDetails)
		if err != nil {
			return true, nil, err
		}
		if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusAccepted); err != nil {
			err = errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
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
	}

	body, err := pollingExecutor.Execute()
	if err != nil {
		return nil, err
	}
	buildScanResponse := BuildScanResponse{}
	if err = json.Unmarshal(body, &buildScanResponse); err != nil {
		return nil, errorutils.CheckError(err)
	}
	if &buildScanResponse == nil || buildScanResponse.Status == xrayScanStatusFailed {
		return nil, errorutils.CheckErrorf("Xray build scan failed")
	}
	return &buildScanResponse, err
}

type XrayBuildParams struct {
	BuildName   string `json:"build_name,omitempty"`
	BuildNumber string `json:"build_number,omitempty"`
	Project     string `json:"project,omitempty"`
}

type RequestBuildScanResponse struct {
	Info string `json:"info,omitempty"`
}

type BuildScanResponse struct {
	Status         string      `json:"status,omitempty"`
	MoreDetailsUrl string      `json:"more_details_url,omitempty"`
	FailBuild      bool        `json:"fail_build,omitempty"`
	Violations     []Violation `json:"violations,omitempty"`
	Info           string      `json:"info,omitempty"`
}

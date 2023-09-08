package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const (
	postScanContextAPI = "api/v1/gitinfo"

	XscGraphAPI = "api/v1/sca/scan/graph"

	multiScanIdParam = "multi_scan_id="

	scanTechQueryParam = "tech="

	XscVersionAPI = "api/v1/system/version"
)

type XscScanService struct {
	ScanService
}

func NewXscScanService(client *jfroghttpclient.JfrogHttpClient, details auth.ServiceDetails) *XscScanService {
	return &XscScanService{ScanService{client: client, XrayDetails: details}}
}

func (xsc *XscScanService) SendScanContext(details *XscGitInfoContext) (multiScanId string, err error) {
	// XscGitInfoContext is optional
	if details == nil {
		return
	}
	httpClientsDetails := xsc.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	requestBody, err := json.Marshal(details)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	url := xsc.XrayDetails.GetXscUrl() + postScanContextAPI
	resp, body, err := xsc.client.SendPost(url, requestBody, &httpClientsDetails)
	if err != nil {
		return
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		scanErrorJson := ScanErrorJson{}
		if e := json.Unmarshal(body, &scanErrorJson); e == nil {
			return "", errorutils.CheckErrorf(scanErrorJson.Error)
		}
		return
	}
	scanResponse := XscPostContextResponse{}
	if err = json.Unmarshal(body, &scanResponse); err != nil {
		return "", errorutils.CheckError(err)
	}
	return scanResponse.MultiScanId, err
}

func (xsc *XscScanService) ScanGraph(scanParams *XrayGraphScanParams) (string, error) {
	httpClientsDetails := xsc.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	requestBody, err := json.Marshal(scanParams.DependenciesGraph)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	url := xsc.XrayDetails.GetXscUrl() + XscGraphAPI
	url += createScanGraphQueryParams(*scanParams)

	resp, body, err := xsc.client.SendPost(url, requestBody, &httpClientsDetails)
	if err != nil {
		return "", err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated); err != nil {
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

func (xsc *XscScanService) GetScanGraphResults(scanId string, _, _ bool) (*ScanResponse, error) {
	httpClientsDetails := xsc.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	// The scan request may take some time to complete. We expect to receive a 202 response, until the completion.
	endPoint := xsc.XrayDetails.GetXscUrl() + XscGraphAPI + "/" + scanId
	log.Info("Waiting for scan to complete on JFrog Xray...")
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		resp, body, _, err := xsc.client.SendGet(endPoint, true, &httpClientsDetails)
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
		Timeout:         DefaultMaxWaitMinutes,
		PollingInterval: DefaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get Dependencies Scan results... ",
	}

	body, err := pollingExecutor.Execute()
	if err != nil {
		return nil, err
	}
	scanResponse := ScanResponse{}
	if err = json.Unmarshal(body, &scanResponse); err != nil {
		return nil, errorutils.CheckErrorf("couldn't parse JFrog Xray server response: " + err.Error())
	}
	if scanResponse.ScannedStatus == XrayScanStatusFailed {
		// Failed due to an internal Xray error
		return nil, errorutils.CheckErrorf("received a failure status from JFrog Xray server:\n%s", errorutils.GenerateErrorString(body))
	}
	return &scanResponse, err
}

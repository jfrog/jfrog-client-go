package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const (
	importGraph    = "api/v1/scan/import"
	importGraphXML = "api/v1/scan/import_xml"
)

func (ss *ScanService) ImportGraph(scanParams XrayGraphImportParams) (string, error) {
	httpClientsDetails := ss.XrayDetails.CreateHttpClientDetails()
	var v interface{}
	err := json.Unmarshal(scanParams.SBOMInput, &v)
	url := ""
	if err == nil {
		utils.SetContentType("application/json", &httpClientsDetails.Headers)
		url = ss.XrayDetails.GetUrl() + importGraph
	} else {
		utils.SetContentType("application/xml", &httpClientsDetails.Headers)
		url = ss.XrayDetails.GetUrl() + importGraphXML
	}
	requestBody := scanParams.SBOMInput
	resp, body, err := ss.client.SendPost(url, requestBody, &httpClientsDetails)
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

func (ss *ScanService) GetImportGraphResults(scanId string) (*ScanResponse, error) {
	httpClientsDetails := ss.XrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	// The scan request may take some time to complete. We expect to receive a 202 response, until the completion.
	endPoint := ss.XrayDetails.GetUrl() + scanGraphAPI
	// Modify endpoint if XSC is enabled
	endPoint += "/" + scanId
	endPoint += includeVulnerabilitiesParam
	log.Info("Waiting for enrich process to complete on JFrog Xray...")
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		resp, body, _, err := ss.client.SendGet(endPoint, true, &httpClientsDetails)
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
	if scanResponse.ScannedStatus == xrayScanStatusFailed {
		// Failed due to an internal Xray error
		return nil, errorutils.CheckErrorf("received a failure status from JFrog Xray server:\n%s", errorutils.GenerateErrorString(body))
	}
	return &scanResponse, err
}

type XrayGraphImportParams struct {
	// A path in Artifactory that this Artifact is intended to be deployed to.
	// This will provide a way to extract the watches that should be applied on this graph
	ScanType          ScanType
	SBOMInput         []byte
	XscGitInfoContext *XscGitInfoContext
	XscVersion        string
}

package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	xrayUtils "github.com/jfrog/jfrog-client-go/xray/services/utils"
	"net/http"
)

const (
	importGraph    = "api/v1/scan/import"
	importGraphXML = "api/v1/scan/import_xml"
)

type EnrichService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewEnrichService creates a new service to enrich CycloneDX xml and jsons.
func NewEnrichService(client *jfroghttpclient.JfrogHttpClient) *EnrichService {
	return &EnrichService{client: client}
}

func (es *EnrichService) ImportGraph(importParams XrayGraphImportParams) (string, error) {
	httpClientsDetails := es.XrayDetails.CreateHttpClientDetails()
	var v interface{}
	// There's an option to run on XML or JSON file so we need to call the correct API accordingly.
	err := json.Unmarshal(importParams.SBOMInput, &v)
	var url string
	if err != nil {
		utils.SetContentType("application/xml", &httpClientsDetails.Headers)
		url = es.XrayDetails.GetUrl() + importGraphXML
	} else {
		httpClientsDetails.SetContentTypeApplicationJson()
		url = es.XrayDetails.GetUrl() + importGraph
	}

	requestBody := importParams.SBOMInput
	resp, body, err := es.client.SendPost(url, requestBody, &httpClientsDetails)
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

func (es *EnrichService) GetImportGraphResults(scanId string) (*ScanResponse, error) {
	httpClientsDetails := es.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	// Getting the import graph results is from the same api but with some parameters always initialized.
	endPoint := es.XrayDetails.GetUrl() + scanGraphAPI + "/" + scanId + includeVulnerabilitiesParam
	log.Info("Waiting for enrich process to complete on JFrog Xray...")
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   xrayUtils.PollingAction(es.client, endPoint, httpClientsDetails),
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

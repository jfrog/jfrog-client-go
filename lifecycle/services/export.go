package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
)

const (
	exportReleaseBundleEndpoint       = "export"
	exportReleaseBundleStatusEndpoint = "export/status"
)

type RlExportedStatus string

const (
	ExportCompleted    RlExportedStatus = "COMPLETED"
	ExportProcessing   RlExportedStatus = "IN_PROGRESS"
	ExportFailed       RlExportedStatus = "FAILED"
	ExportNotTriggered RlExportedStatus = "NOT_TRIGGERED"
)

type ReleaseBundleExportedStatusResponse struct {
	Status      RlExportedStatus `json:"status,omitempty"`
	RelativeUrl string           `json:"relative_download_url,omitempty"`
	DownloadUrl string           `json:"download_url,omitempty"`
}

type ReleaseBundleExportParams struct {
	ReleaseBundleDetails
	Modifications `json:"modifications"`
}

type RbExportBody struct {
	ReleaseBundleExportParams
}

type exportStatusOperation struct {
	reqBody     RbExportBody
	queryParams CommonOptionalQueryParams
}

func (exs *exportStatusOperation) getOperationRestApi() string {
	fullUrl, err := url.JoinPath(distributionBaseApi, exportReleaseBundleStatusEndpoint, exs.reqBody.ReleaseBundleName, exs.reqBody.ReleaseBundleVersion)
	if err != nil {
		panic(err)
	}
	return fullUrl
}

func (exs *exportStatusOperation) getRequestBody() any {
	return exs.reqBody
}

func (exs *exportStatusOperation) getOperationSuccessfulMsg() string {
	return "Successfully received Release Bundle export status"
}

func (exs *exportStatusOperation) getOperationParams() CommonOptionalQueryParams {
	return exs.queryParams
}

func (exs *exportStatusOperation) getSigningKeyName() string {
	return ""
}

type exportOperation struct {
	reqBody     RbExportBody
	queryParams CommonOptionalQueryParams
}

func (exp *exportOperation) getOperationRestApi() string {
	fullUrl, err := url.JoinPath(distributionBaseApi, exportReleaseBundleEndpoint, exp.reqBody.ReleaseBundleName, exp.reqBody.ReleaseBundleVersion)
	if err != nil {
		panic(err)
	}
	return fullUrl
}

func (exp *exportOperation) getRequestBody() any {
	return exp.reqBody
}

func (exp *exportOperation) getOperationSuccessfulMsg() string {
	return "Release Bundle successfully exported"
}

func (exp *exportOperation) getOperationParams() CommonOptionalQueryParams {
	return exp.queryParams
}

func (exp *exportOperation) getSigningKeyName() string {
	return ""
}

func (rbs *ReleaseBundlesService) ExportReleaseBundle(rlExportParams *ReleaseBundleExportParams, queryParams CommonOptionalQueryParams) (err error) {
	var rlExportedResponse ReleaseBundleExportedStatusResponse
	//Check the current status
	if rlExportedResponse, _, err = rbs.getExportedReleaseBundleStatus(rlExportParams, queryParams); err != nil {
		return
	}
	// Trigger export if needed
	if rlExportedResponse.Status == ExportNotTriggered {
		if err = rbs.triggerRlExport(rlExportParams, queryParams); err != nil {
			return
		}
	}
	// Wait for export to finish
	// TODO what if the status is on different state? like failure? check this issue.
	if rlExportedResponse.Status == ExportProcessing {
		if rlExportedResponse, err = rbs.checkExportedStatusWithRetries(rlExportParams, queryParams); err != nil {
			return
		}
	}
	return
}

func (rbs *ReleaseBundlesService) checkExportedStatusWithRetries(rlExportParams *ReleaseBundleExportParams, queryParams CommonOptionalQueryParams) (response ReleaseBundleExportedStatusResponse, err error) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		response, responseBody, err = rbs.getExportedReleaseBundleStatus(rlExportParams, queryParams)
		if err != nil {
			return
		}
		if err != nil {
			return true, nil, err
		}
		switch response.Status {
		case ExportProcessing:
			return false, nil, nil
		case ExportCompleted, ExportFailed:
			return true, responseBody, nil
		default:
			return true, nil, errorutils.CheckErrorf("received unexpected status: '%s'", response.Status)
		}
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWait,
		PollingInterval: SyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       fmt.Sprintf("Getting Exported Release Bundle %s/%s status...", rlExportParams.ReleaseBundleName, rlExportParams.ReleaseBundleVersion),
	}
	_, err = pollingExecutor.Execute()
	return
}

func (rbs *ReleaseBundlesService) triggerRlExport(rlExportParams *ReleaseBundleExportParams, queryParams CommonOptionalQueryParams) (err error) {
	operation := &exportOperation{
		reqBody:     RbExportBody{*rlExportParams},
		queryParams: queryParams,
	}
	log.Debug("Triggering Release Bundle Export process...")
	_, err = rbs.doOperation(operation)
	return
}

func (rbs *ReleaseBundlesService) getExportedReleaseBundleStatus(rlExportParams *ReleaseBundleExportParams, queryParams CommonOptionalQueryParams) (exportedStatusResponse ReleaseBundleExportedStatusResponse, body []byte, err error) {
	operation := &exportStatusOperation{
		reqBody:     RbExportBody{*rlExportParams},
		queryParams: queryParams,
	}
	log.Debug("Getting Release Bundle Export status...")
	respBody, err := rbs.doGetOperation(operation)
	err = json.Unmarshal(respBody, &exportedStatusResponse)
	return
}

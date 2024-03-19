package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"path"
)

const (
	releaseBundleExportEndpoint = "export"
	releaseBundleStatusEndpoint = "status"
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
type exportStatusOperation struct {
	ReleaseBundleDetails
	queryParams CommonOptionalQueryParams
}

func (exs *exportStatusOperation) getOperationRestApi() string {
	return path.Join(distributionBaseApi, releaseBundleExportEndpoint, releaseBundleStatusEndpoint, exs.ReleaseBundleName, exs.ReleaseBundleVersion)
}

func (exs *exportStatusOperation) getRequestBody() any { return nil }

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
	ReleaseBundleDetails
	modifications Modifications
	queryParams   CommonOptionalQueryParams
}

func (exp *exportOperation) getOperationRestApi() string {
	return path.Join(distributionBaseApi, releaseBundleExportEndpoint, exp.ReleaseBundleName, exp.ReleaseBundleVersion)
}

func (exp *exportOperation) getRequestBody() any { return exp.modifications }

func (exp *exportOperation) getOperationSuccessfulMsg() string {
	return "Release Bundle successfully exported"
}

func (exp *exportOperation) getOperationParams() CommonOptionalQueryParams {
	return exp.queryParams
}

func (exp *exportOperation) getSigningKeyName() string {
	return ""
}

func (rbs *ReleaseBundlesService) ExportReleaseBundle(rbDetails ReleaseBundleDetails, modifications Modifications, queryParams CommonOptionalQueryParams) (exportResponse ReleaseBundleExportedStatusResponse, err error) {
	// Check the current status
	if exportResponse, err = rbs.getExportedReleaseBundleStatus(rbDetails, queryParams); err != nil {
		return
	}
	if exportResponse.Status == ExportCompleted {
		return
	}
	// Trigger export
	if err = rbs.triggerReleaseBundleExportProcess(rbDetails, modifications, queryParams); err != nil {
		return
	}
	// Wait for export to finish
	exportResponse, err = rbs.waitForExport(rbDetails, queryParams)
	return
}

func (rbs *ReleaseBundlesService) waitForExport(rbDetails ReleaseBundleDetails, queryParams CommonOptionalQueryParams) (response ReleaseBundleExportedStatusResponse, err error) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		response, err = rbs.getExportedReleaseBundleStatus(rbDetails, queryParams)
		if err != nil {
			return
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
		MsgPrefix:       fmt.Sprintf("Getting Exported Release Bundle %s/%s status...", rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion),
	}
	_, err = pollingExecutor.Execute()
	return
}

func (rbs *ReleaseBundlesService) triggerReleaseBundleExportProcess(rbDetails ReleaseBundleDetails, modifications Modifications, queryParams CommonOptionalQueryParams) (err error) {
	operation := &exportOperation{
		ReleaseBundleDetails: rbDetails,
		modifications:        modifications,
		queryParams:          queryParams,
	}
	log.Debug("Triggering Release Bundle Export process...")
	_, err = rbs.doPostOperation(operation)
	return
}

func (rbs *ReleaseBundlesService) getExportedReleaseBundleStatus(rbDetails ReleaseBundleDetails, queryParams CommonOptionalQueryParams) (exportedStatusResponse ReleaseBundleExportedStatusResponse, err error) {
	operation := &exportStatusOperation{
		ReleaseBundleDetails: rbDetails,
		queryParams:          queryParams,
	}
	log.Debug("Getting Release Bundle Export status...")
	respBody, err := rbs.doGetOperation(operation)
	if err != nil {
		return
	}
	err = errorutils.CheckError(json.Unmarshal(respBody, &exportedStatusResponse))
	return
}

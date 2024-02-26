package services

import (
	"encoding/json"
	"fmt"
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
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

func (rbs *ReleaseBundlesService) ExportReleaseBundle(params distribution.ReleaseBundleExportParams) (err error) {
	var rlExportedResponse ReleaseBundleExportedStatusResponse
	// Check the current status
	if rlExportedResponse, _, err = rbs.getExportedReleaseBundleStatus(params); err != nil {
		return
	}
	// Trigger export if needed
	if rlExportedResponse.Status == ExportNotTriggered {
		if err = rbs.triggerRlExport(params); err != nil {
			return
		}
		// Wait for export to finish
		if rlExportedResponse, err = rbs.checkExportedStatusWithRetries(params); err != nil {
			return
		}
	}
	// Download the archive
	return rbs.downloadExportedReleaseBundle(rlExportedResponse, params)
}

func (rbs *ReleaseBundlesService) checkExportedStatusWithRetries(params distribution.ReleaseBundleExportParams) (response ReleaseBundleExportedStatusResponse, err error) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		response, responseBody, err = rbs.getExportedReleaseBundleStatus(params)
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
		MsgPrefix:       fmt.Sprintf("Getting Exported Release Bundle %s/%s status...", params.Name, params.Version),
	}
	_, err = pollingExecutor.Execute()
	return
}

func (rbs *ReleaseBundlesService) triggerRlExport(params distribution.ReleaseBundleExportParams) error {

	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), GetReleaseBundleExportRestApi(params), nil)
	if err != nil {
		return err
	}

	httpClientDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)
	resp, body, err := rbs.client.SendPost(requestFullUrl, nil, &httpClientDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted)
}

func (rbs *ReleaseBundlesService) getExportedReleaseBundleStatus(params distribution.ReleaseBundleExportParams) (exportedStatusResponse ReleaseBundleExportedStatusResponse, body []byte, err error) {

	httpClientDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)
	resp, body, _, err := rbs.client.SendGet(rbs.GetLifecycleDetails().GetUrl()+GetReleaseBundleExportStatusRestApi(params), false, &httpClientDetails)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}
	err = errorutils.CheckError(json.Unmarshal(body, &exportedStatusResponse))
	return
}
func (rbs *ReleaseBundlesService) downloadExportedReleaseBundle(rlExportedResponse ReleaseBundleExportedStatusResponse, params distribution.ReleaseBundleExportParams) (err error) {
	httpClientDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	// Build download URL while replacing endpoint of lifecycle to artifactory
	downloadUrl := strings.Replace(rbs.GetLifecycleDetails().GetUrl(), "lifecycle", "artifactory", 1)
	downloadPath, err := url.JoinPath(downloadUrl, rlExportedResponse.RelativeUrl)
	if err != nil {
		return
	}
	// Download file into current working dir
	fileName := fmt.Sprintf("%s-%s.zip", params.Name, params.Version)
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	fileDetails := httpclient.DownloadFileDetails{
		FileName:      fileName,
		DownloadPath:  downloadPath,
		LocalPath:     wd,
		LocalFileName: fileName,
	}
	_, err = rbs.client.DownloadFileWithProgress(&fileDetails, "Download Exported Bundle", &httpClientDetails, false, false, nil)
	return
}

func GetReleaseBundleExportRestApi(rbDetails distribution.ReleaseBundleExportParams) string {
	return path.Join(distributionBaseApi, exportReleaseBundleEndpoint, rbDetails.Name, rbDetails.Version)
}

func GetReleaseBundleExportStatusRestApi(rbDetails distribution.ReleaseBundleExportParams) string {
	return path.Join(distributionBaseApi, exportReleaseBundleStatusEndpoint, rbDetails.Name, rbDetails.Version)
}

package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
	"path"
	"time"
)

const (
	statusesApi              = "statuses"
	defaultMaxWait           = 60 * time.Minute
	DefaultSyncSleepInterval = 10 * time.Second
)

var SyncSleepInterval = DefaultSyncSleepInterval

type RbStatus string

const (
	Completed  RbStatus = "COMPLETED"
	Processing RbStatus = "PROCESSING"
	Pending    RbStatus = "PENDING"
	Failed     RbStatus = "FAILED"
	Rejected   RbStatus = "REJECTED"
	Deleting   RbStatus = "DELETING"
)

func (rbs *ReleaseBundlesService) GetReleaseBundleCreationStatus(rbDetails ReleaseBundleDetails, projectKey string, sync bool) (ReleaseBundleStatusResponse, error) {
	return rbs.getReleaseBundleOperationStatus(GetReleaseBundleCreationStatusRestApi(rbDetails), projectKey, sync, "creation")
}

func GetReleaseBundleCreationStatusRestApi(rbDetails ReleaseBundleDetails) string {
	return path.Join(releaseBundleBaseApi, statusesApi, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
}

func (rbs *ReleaseBundlesService) GetReleaseBundlePromotionStatus(rbDetails ReleaseBundleDetails, projectKey, createdMillis string, sync bool) (ReleaseBundleStatusResponse, error) {
	restApi := path.Join(promotionBaseApi, statusesApi, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion, createdMillis)
	return rbs.getReleaseBundleOperationStatus(restApi, projectKey, sync, "promotion")
}

func (rbs *ReleaseBundlesService) getReleaseBundleOperationStatus(restApi string, projectKey string, sync bool, operationStr string) (ReleaseBundleStatusResponse, error) {
	if sync {
		return rbs.waitForRbOperationCompletion(restApi, projectKey, operationStr)
	}
	statusResp, _, err := rbs.getReleaseBundleStatus(restApi, projectKey)
	return statusResp, err
}

func (rbs *ReleaseBundlesService) getReleaseBundleStatus(restApi string, projectKey string) (statusResp ReleaseBundleStatusResponse, body []byte, err error) {
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, getProjectQueryParam(projectKey))
	if err != nil {
		return
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	resp, body, _, err := rbs.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}
	err = errorutils.CheckError(json.Unmarshal(body, &statusResp))
	return
}

func getStatusResponse(respBody []byte) (ReleaseBundleStatusResponse, error) {
	var rbStatusResponse ReleaseBundleStatusResponse
	err := json.Unmarshal(respBody, &rbStatusResponse)
	return rbStatusResponse, errorutils.CheckError(err)
}

func (rbs *ReleaseBundlesService) waitForRbOperationCompletion(restApi, projectKey, operation string) (ReleaseBundleStatusResponse, error) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		var rbStatusResponse ReleaseBundleStatusResponse
		rbStatusResponse, responseBody, err = rbs.getReleaseBundleStatus(restApi, projectKey)
		if err != nil {
			return true, nil, err
		}
		switch rbStatusResponse.Status {
		case Pending, Processing:
			return false, nil, nil
		case Completed, Rejected, Failed, Deleting:
			return true, responseBody, nil
		default:
			return true, nil, errorutils.CheckErrorf("received unexpected status: '%s'", rbStatusResponse.Status)
		}
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWait,
		PollingInterval: SyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       fmt.Sprintf("Getting Release Bundle %s status...", operation),
	}
	finalRespBody, err := pollingExecutor.Execute()
	if err != nil {
		return ReleaseBundleStatusResponse{}, err
	}
	return getStatusResponse(finalRespBody)
}

type ReleaseBundleStatusResponse struct {
	Status   RbStatus  `json:"status,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type Message struct {
	Source  string `json:"source,omitempty"`
	Text    string `json:"text,omitempty"`
	Created string `json:"created,omitempty"`
}

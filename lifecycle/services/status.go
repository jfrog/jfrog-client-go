package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"time"
)

const (
	recordsApi               = "records"
	statusesApi              = "statuses"
	trackersApi              = "trackers"
	defaultMaxWait           = 60 * time.Minute
	DefaultSyncSleepInterval = 10 * time.Second
)

var SyncSleepInterval = DefaultSyncSleepInterval

type RbStatus string

const (
	Completed  RbStatus = "COMPLETED"
	Processing RbStatus = "PROCESSING"
	InProgress RbStatus = "IN_PROGRESS"
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

func GetReleaseBundleSpecificationRestApi(rbDetails ReleaseBundleDetails) string {
	return path.Join(releaseBundleBaseApi, recordsApi, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
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
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, distribution.GetProjectQueryParam(projectKey))
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

func (rbs *ReleaseBundlesService) GetReleaseBundleSpecification(rbDetails ReleaseBundleDetails) (specResp ReleaseBundleSpecResponse, err error) {
	restApi := GetReleaseBundleSpecificationRestApi(rbDetails)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, nil)
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
	err = errorutils.CheckError(json.Unmarshal(body, &specResp))
	return
}

func (dr *DistributeReleaseBundleService) getReleaseBundleDistributionStatus(distributeParams *distribution.DistributionParams, trackerId json.Number) (statusResp *distribution.DistributionStatusResponse, body []byte, err error) {
	restApi := path.Join(distributionBaseApi, trackersApi, distributeParams.Name, distributeParams.Version, trackerId.String())
	requestFullUrl, err := utils.BuildUrl(dr.LcDetails.GetUrl(), restApi, distribution.GetProjectQueryParam(dr.GetProjectKey()))
	if err != nil {
		return
	}
	httpClientsDetails := dr.LcDetails.CreateHttpClientDetails()
	resp, body, _, err := dr.client.SendGet(requestFullUrl, true, &httpClientsDetails)
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

func (dr *DistributeReleaseBundleService) waitForDistributionOperationCompletion(distributeParams *distribution.DistributionParams, trackerId json.Number) error {
	maxWait := time.Duration(dr.GetMaxWaitMinutes()) * time.Minute
	if maxWait.Minutes() < 1 {
		maxWait = defaultMaxWait
	}

	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		statusResponse, responseBody, err := dr.getReleaseBundleDistributionStatus(distributeParams, trackerId)
		if err != nil {
			return true, nil, err
		}

		switch statusResponse.Status {
		case distribution.NotDistributed, distribution.InProgress, distribution.InQueue:
			return false, nil, nil
		case distribution.Failed, distribution.Completed:
			return true, responseBody, nil
		default:
			return true, nil, errorutils.CheckErrorf("received unexpected status: '%s'", statusResponse.Status)
		}
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         maxWait,
		PollingInterval: SyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       fmt.Sprintf("Sync: Distributing %s/%s...", distributeParams.Name, distributeParams.Version),
	}
	finalRespBody, err := pollingExecutor.Execute()
	if err != nil {
		return err
	}

	var dsStatusResponse distribution.DistributionStatusResponse
	if err = json.Unmarshal(finalRespBody, &dsStatusResponse); err != nil {
		return errorutils.CheckError(err)
	}

	if dsStatusResponse.Status != distribution.Completed {
		for _, st := range dsStatusResponse.Sites {
			if st.Status != distribution.Completed {
				err = errors.Join(err, fmt.Errorf("target %s name:%s error:%s", st.TargetArtifactory.Type, st.TargetArtifactory.Name, st.Error))
			}
		}
		return errorutils.CheckError(err)
	}
	log.Info("Distribution Completed!")
	return nil
}

type ReleaseBundleStatusResponse struct {
	Status   RbStatus  `json:"status,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type ReleaseBundleSpecResponse struct {
	CreatedBy     string    `json:"created_by,omitempty"`
	Created       time.Time `json:"created"`
	CreatedMillis int       `json:"created_millis,omitempty"`
	Artifacts     []struct {
		Path                string `json:"path,omitempty"`
		Checksum            string `json:"checksum,omitempty"`
		SourceRepositoryKey string `json:"source_repository_key,omitempty"`
		PackageType         string `json:"package_type,omitempty"`
		Size                int    `json:"size,omitempty"`
		Properties          []struct {
			Key    string   `json:"key"`
			Values []string `json:"values"`
		} `json:"properties,omitempty"`
	} `json:"artifacts,omitempty"`
}

type Message struct {
	Source  string `json:"source,omitempty"`
	Text    string `json:"text,omitempty"`
	Created string `json:"created,omitempty"`
}

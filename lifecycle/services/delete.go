package services

import (
	"encoding/json"
	"fmt"
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"strconv"
	"time"
)

const (
	remoteDeleteEndpoint = "remote_delete"
)

func (rbs *ReleaseBundlesService) DeleteReleaseBundleVersion(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams) error {
	restApi := path.Join(releaseBundleBaseApi, records, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
	return rbs.deleteReleaseBundle(params, restApi)
}

func (rbs *ReleaseBundlesService) DeleteReleaseBundleVersionPromotion(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, createdMillis string) error {
	restApi := path.Join(promotionBaseApi, records, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion, createdMillis)
	return rbs.deleteReleaseBundle(params, restApi)
}

func (rbs *ReleaseBundlesService) deleteReleaseBundle(params CommonOptionalQueryParams, restApi string) error {
	queryParams := distribution.GetProjectQueryParam(params.ProjectKey)
	queryParams[async] = strconv.FormatBool(params.Async)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, queryParams)
	if err != nil {
		return err
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	resp, body, err := rbs.client.SendDelete(requestFullUrl, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if params.Async {
		return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}

func (rbs *ReleaseBundlesService) RemoteDeleteReleaseBundle(rbDetails ReleaseBundleDetails, params ReleaseBundleRemoteDeleteParams) error {
	dryRunStr := ""
	if params.DryRun {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Remote Deleting: " + rbDetails.ReleaseBundleName + "/" + rbDetails.ReleaseBundleVersion)

	rbBody := distribution.CreateDistributeV1Body(params.DistributionRules, params.DryRun, false)
	content, err := json.Marshal(rbBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	restApi := GetRemoteDeleteReleaseBundleApi(rbDetails)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, nil)
	if err != nil {
		return err
	}

	httpClientDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)
	resp, body, err := rbs.client.SendPost(requestFullUrl, content, &httpClientDetails)
	if err != nil {
		return err
	}

	log.Debug("Artifactory response:", resp.Status)
	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted)
	if err != nil || params.Async || params.DryRun {
		return err
	}

	return rbs.waitForRemoteDeletion(rbDetails, params)
}

func GetRemoteDeleteReleaseBundleApi(rbDetails ReleaseBundleDetails) string {
	return path.Join(distributionBaseApi, remoteDeleteEndpoint, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
}

func (rbs *ReleaseBundlesService) waitForRemoteDeletion(rbDetails ReleaseBundleDetails, params ReleaseBundleRemoteDeleteParams) error {
	maxWaitTime := defaultMaxWait
	if params.MaxWaitMinutes > 0 {
		maxWaitTime = time.Duration(params.MaxWaitMinutes) * time.Minute
	}

	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		resp, _, err := rbs.getReleaseBundleDistributions(rbDetails, params.ProjectKey)
		if err != nil {
			return true, nil, err
		}
		deletionStatus := resp[len(resp)-1].Status
		switch deletionStatus {
		case InProgress:
			return false, nil, nil
		case Completed:
			return true, nil, nil
		case Failed:
			return true, nil, errorutils.CheckErrorf("remote deletion failed!")
		default:
			return true, nil, errorutils.CheckErrorf("unexpected status for remote deletion: %s", deletionStatus)
		}
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         maxWaitTime,
		PollingInterval: SyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       fmt.Sprintf("Performing sync remote deletion of release bundle %s/%s...", rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion),
	}
	_, err := pollingExecutor.Execute()
	return err
}

type ReleaseBundleRemoteDeleteParams struct {
	DistributionRules []*distribution.DistributionCommonParams
	DryRun            bool
	// Max time in minutes to wait for sync distribution to finish.
	MaxWaitMinutes int
	CommonOptionalQueryParams
}

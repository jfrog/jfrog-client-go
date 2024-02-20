package services

import (
	"encoding/json"
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"path"
	"strconv"
)

const (
	remoteDeleteEndpoint = "remote_delete"
)

func (rbs *ReleaseBundlesService) DeleteReleaseBundle(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams) error {
	queryParams := getProjectQueryParam(params.ProjectKey)
	queryParams[async] = strconv.FormatBool(params.Async)
	restApi := path.Join(releaseBundleBaseApi, records, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, queryParams)
	if err != nil {
		return err
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	resp, body, err := rbs.client.SendDelete(requestFullUrl, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}

func (rbs *ReleaseBundlesService) RemoteDeleteReleaseBundle(params distribution.DistributionParams, dryRun bool) error {
	rbBody := distribution.CreateDistributeV1Body(params, dryRun, false)
	content, err := json.Marshal(rbBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	restApi := path.Join(distributionBaseApi, remoteDeleteEndpoint, params.Name, params.Version)
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

	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted)
}

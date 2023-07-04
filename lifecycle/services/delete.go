package services

import (
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"path"
	"strconv"
)

func (rbs *ReleaseBundlesService) DeleteReleaseBundle(rbDetails ReleaseBundleDetails, params ReleaseBundleQueryParams) error {
	queryParams := getProjectQueryParam(params.ProjectKey)
	queryParams[async] = strconv.FormatBool(params.Async)
	restApi := path.Join(releaseBundleBaseApi, records, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
	requestFullUrl, err := rtUtils.BuildArtifactoryUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, queryParams)
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

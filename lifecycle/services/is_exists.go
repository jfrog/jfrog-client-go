package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
)

const (
	isExistInRbV2Endpoint = "api/v2/release_bundle/existence"
)

func (rbs *ReleaseBundlesService) IsExists(projectName, releaseBundleNameAndVersion string) (bool, error) {
	queryParams := distribution.GetProjectQueryParam(projectName)
	restApi := path.Join(isExistInRbV2Endpoint, releaseBundleNameAndVersion)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, queryParams)

	if err != nil {
		return false, err
	}

	httpClientDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	resp, body, _, err := rbs.client.SendGet(requestFullUrl, true, &httpClientDetails)
	if err != nil {
		return false, err
	}
	log.Debug("Artifactory response:", resp.Status)

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted, http.StatusOK); err != nil {
		return false, err
	}

	response := &isReleaseBundleExistResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return false, err
	}

	return response.Exists, nil
}

func GetIsExistReleaseBundleApi(releaseBundleNameAndVersion string) string {
	return path.Join(isExistInRbV2Endpoint, releaseBundleNameAndVersion)
}

type isReleaseBundleExistResponse struct {
	Exists bool `json:"exists"`
}

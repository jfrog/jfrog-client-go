package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BuildInfoService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewBuildInfoService(client *rthttpclient.ArtifactoryHttpClient) *BuildInfoService {
	return &BuildInfoService{client: client}
}

func (bis *BuildInfoService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return bis.ArtDetails
}

func (bis *BuildInfoService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	bis.ArtDetails = rt
}

func (bis *BuildInfoService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return bis.client, nil
}

func (bis *BuildInfoService) IsDryRun() bool {
	return false
}

type BuildInfoParams struct {
	BuildName   string
	BuildNumber string
}

func (bis *BuildInfoService) GetBuildInfo(params BuildInfoParams) (*buildinfo.BuildInfo, error) {
	// Resolve LATEST build number from Artifactory if required.
	name, number, err := utils.GetBuildNameAndNumberFromArtifactory(params.BuildName, params.BuildNumber, bis)
	if err != nil {
		return nil, err
	}

	// Get build-info json from Artifactory.
	httpClientsDetails := bis.GetArtifactoryDetails().CreateHttpClientDetails()
	buildInfoUrl := fmt.Sprintf("%sapi/build/%s/%s", bis.GetArtifactoryDetails().GetUrl(), name, number)
	log.Debug("Getting build-info from: ", buildInfoUrl)
	_, body, _, err := bis.client.SendGet(buildInfoUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}

	// Build BuildInfo struct from json.
	var buildInfoJson struct {
		BuildInfo buildinfo.BuildInfo `json:"buildInfo,omitempty"`
	}
	if err := json.Unmarshal(body, &buildInfoJson); err != nil {
		return nil, err
	}

	return &buildInfoJson.BuildInfo, nil
}

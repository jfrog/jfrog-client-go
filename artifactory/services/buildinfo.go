package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
)

type BuildInfoService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
	DryRun     bool
}

func NewBuildInfoService(client *rthttpclient.ArtifactoryHttpClient) *BuildInfoService {
	return &BuildInfoService{client: client}
}

func (bis *BuildInfoService) GetArtifactoryDetails() auth.ServiceDetails {
	return bis.ArtDetails
}

func (bis *BuildInfoService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	bis.ArtDetails = rt
}

func (bis *BuildInfoService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return bis.client, nil
}

func (bis *BuildInfoService) IsDryRun() bool {
	return bis.DryRun
}

type BuildInfoParams struct {
	BuildName   string
	BuildNumber string
}

func NewBuildInfoParams() BuildInfoParams {
	return BuildInfoParams{}
}

// Returns the build info and it's uri of the provided parameters.
func (bis *BuildInfoService) GetBuildInfo(params BuildInfoParams) (*buildinfo.PublishedBuildInfo, error) {
	// Resolve LATEST build number from Artifactory if required.
	name, number, err := utils.GetBuildNameAndNumberFromArtifactory(params.BuildName, params.BuildNumber, bis)
	if err != nil {
		return nil, err
	}

	// Get build-info json from Artifactory.
	httpClientsDetails := bis.GetArtifactoryDetails().CreateHttpClientDetails()

	restApi := path.Join("api/build/", name, number)
	requestFullUrl, err := utils.BuildArtifactoryUrl(bis.GetArtifactoryDetails().GetUrl(), restApi, make(map[string]string))

	log.Debug("Getting build-info from: ", requestFullUrl)
	resp, body, _, err := bis.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	// Build BuildInfo struct from json.
	publishedBuildInfo := &buildinfo.PublishedBuildInfo{}
	if err := json.Unmarshal(body, publishedBuildInfo); err != nil {
		return nil, err
	}

	return publishedBuildInfo, nil
}

func (bis *BuildInfoService) PublishBuildInfo(build *buildinfo.BuildInfo) error {
	content, err := json.Marshal(build)
	if errorutils.CheckError(err) != nil {
		return err
	}
	if bis.IsDryRun() {
		log.Info("[Dry run] Logging Build info preview...")
		log.Output(clientutils.IndentJson(content))
		return nil
	}
	httpClientsDetails := bis.GetArtifactoryDetails().CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory+json", &httpClientsDetails.Headers)
	log.Info("Deploying build info...")
	resp, body, err := bis.client.SendPut(bis.ArtDetails.GetUrl()+"api/build/", content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Build info successfully deployed. Browse it in Artifactory under " + bis.GetArtifactoryDetails().GetUrl() + "webapp/builds/" + build.Name + "/" + build.Number)
	return nil
}

package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BuildInfoService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
	Project    string
	DryRun     bool
}

func NewBuildInfoService(client *jfroghttpclient.JfrogHttpClient) *BuildInfoService {
	return &BuildInfoService{client: client}
}

func (bis *BuildInfoService) GetArtifactoryDetails() auth.ServiceDetails {
	return bis.ArtDetails
}

func (bis *BuildInfoService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	bis.ArtDetails = rt
}

func (bis *BuildInfoService) GetJfrogHttpClient() (*jfroghttpclient.JfrogHttpClient, error) {
	return bis.client, nil
}

func (bis *BuildInfoService) GetProject() string {
	return bis.Project
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

// Returns the build info and its uri of the provided parameters.
// If build info was not found (404), returns found=false (with error nil).
// For any other response that isn't 200, an error is returned.
func (bis *BuildInfoService) GetBuildInfo(params BuildInfoParams) (pbi *buildinfo.PublishedBuildInfo, found bool, err error) {
	// Resolve LATEST build number from Artifactory if required.
	name, number, err := utils.GetBuildNameAndNumberFromArtifactory(params.BuildName, params.BuildNumber, bis)
	if err != nil {
		return nil, false, err
	}

	// Get build-info json from Artifactory.
	httpClientsDetails := bis.GetArtifactoryDetails().CreateHttpClientDetails()

	restApi := path.Join("api/build/", name, number)
	requestFullUrl, err := utils.BuildArtifactoryUrl(bis.GetArtifactoryDetails().GetUrl(), restApi, make(map[string]string))

	log.Debug("Getting build-info from: ", requestFullUrl)
	resp, body, _, err := bis.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Debug("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
		return nil, false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	// Build BuildInfo struct from json.
	publishedBuildInfo := &buildinfo.PublishedBuildInfo{}
	if err := json.Unmarshal(body, publishedBuildInfo); err != nil {
		return nil, true, err
	}

	return publishedBuildInfo, true, nil
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
	resp, body, err := bis.client.SendPut(bis.ArtDetails.GetUrl()+"api/build"+bis.getProjectQueryParam(), content, &httpClientsDetails)
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

func (bis *BuildInfoService) getProjectQueryParam() string {
	if bis.Project == "" {
		return ""
	}
	return "?project=" + bis.Project
}

package services

import (
	"encoding/json"
	"errors"
	"net/http"

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

func (bis *BuildInfoService) IsDryRun() bool {
	return bis.DryRun
}

type BuildInfoParams struct {
	BuildName   string
	BuildNumber string
	ProjectKey  string
}

func NewBuildInfoParams() BuildInfoParams {
	return BuildInfoParams{}
}

// Returns the build info and its uri of the provided parameters.
// If build info was not found (404), returns found=false (with error nil).
// For any other response that isn't 200, an error is returned.
func (bis *BuildInfoService) GetBuildInfo(params BuildInfoParams) (pbi *buildinfo.PublishedBuildInfo, found bool, err error) {
	return utils.GetBuildInfo(params.BuildName, params.BuildNumber, params.ProjectKey, bis)
}

func (bis *BuildInfoService) PublishBuildInfo(build *buildinfo.BuildInfo, projectKey string) error {
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
	resp, body, err := bis.client.SendPut(bis.ArtDetails.GetUrl()+"api/build"+utils.GetProjectQueryParam(projectKey), content, &httpClientsDetails)
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

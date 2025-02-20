package services

import (
	"encoding/json"
	"fmt"
	buildinfo "github.com/jfrog/build-info-go/entities"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BuildInfoService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
	DryRun     bool
}

type DeleteBuildInfoBody struct {
	BuildName       string   `json:"buildName,omitempty"`
	Project         string   `json:"project,omitempty"`
	BuildNumber     []string `json:"buildNumbers,omitempty"`
	DeleteArtifacts bool     `json:"deleteArtifacts,omitempty"`
	DeleteAll       bool     `json:"deleteAll,omitempty"`
}

func NewBuildInfoService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *BuildInfoService {
	return &BuildInfoService{artDetails: &artDetails, client: client}
}

func (bis *BuildInfoService) GetArtifactoryDetails() auth.ServiceDetails {
	return *bis.artDetails
}

func (bis *BuildInfoService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return bis.client
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

// Returns the build runs for the requested build info name.
// If build info was not found (404), returns found=false (with error nil).
// For any other response that isn't 200, an error is returned.
func (bis *BuildInfoService) GetBuildRuns(params BuildInfoParams) (runs *buildinfo.BuildRuns, found bool, err error) {
	return utils.GetBuildRuns(params.BuildName, params.ProjectKey, bis)
}

func (bis *BuildInfoService) PublishBuildInfo(build *buildinfo.BuildInfo, projectKey string) (*clientutils.Sha256Summary, error) {
	summary := clientutils.NewSha256Summary()
	content, err := json.Marshal(build)
	if errorutils.CheckError(err) != nil {
		return summary, err
	}
	if bis.IsDryRun() {
		log.Info("[Dry run] Logging Build info preview...")
		log.Output(clientutils.IndentJson(content))
		return summary, nil
	}
	httpClientsDetails := bis.GetArtifactoryDetails().CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory+json", &httpClientsDetails.Headers)
	log.Info(fmt.Sprintf("Publishing build info for <%s>/<%s>...", build.Name, build.Number))
	resp, body, err := bis.client.SendPut(bis.GetArtifactoryDetails().GetUrl()+"api/build"+utils.GetProjectQueryParam(projectKey), content, &httpClientsDetails)
	if err != nil {
		return summary, fmt.Errorf("error occurred while publishing build info: %s", err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated, http.StatusNoContent); err != nil {
		return summary, err
	}
	summary.SetSucceeded(true)
	summary.SetSha256(resp.Header.Get("X-Checksum-Sha256"))

	log.Debug("Artifactory response:", resp.Status)
	return summary, nil
}

func (bis *BuildInfoService) DeleteBuildInfo(build *buildinfo.BuildInfo, projectKey string, numberOfBuildOccurrencesToBeDeleted int) error {
	params := createDeleteBuildInfoBody(build, projectKey, numberOfBuildOccurrencesToBeDeleted)
	content, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	if bis.IsDryRun() {
		log.Info("[Dry run] Deleting build info preview...")
		log.Output(clientutils.IndentJson(content))
		return nil
	}
	httpClientsDetails := bis.GetArtifactoryDetails().CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()
	resp, body, err := bis.client.SendPost(bis.GetArtifactoryDetails().GetUrl()+"api/build/delete", content, &httpClientsDetails)
	if err != nil {
		return fmt.Errorf("error occurred while deleting build info: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusNoContent); err != nil {
		return err
	}
	log.Debug("Artifactory response", resp.Status)
	return nil
}

func createDeleteBuildInfoBody(build *buildinfo.BuildInfo, projectKey string, numberOfBuildOccurrencesToBeDeleted int) DeleteBuildInfoBody {
	buildNumbers := make([]string, 0, numberOfBuildOccurrencesToBeDeleted)
	for i := 0; i < numberOfBuildOccurrencesToBeDeleted; i++ {
		buildNumbers = append(buildNumbers, build.Number)
	}
	return DeleteBuildInfoBody{
		BuildName:       build.Name,
		BuildNumber:     buildNumbers,
		Project:         projectKey,
		DeleteArtifacts: false,
		DeleteAll:       false,
	}
}

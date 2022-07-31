package services

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DockerPromoteService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewDockerPromoteService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *DockerPromoteService {
	return &DockerPromoteService{artDetails: &artDetails, client: client}
}

func (ps *DockerPromoteService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ps.artDetails
}

func (ps *DockerPromoteService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ps.client
}

func (ps *DockerPromoteService) IsDryRun() bool {
	return false
}

func (ps *DockerPromoteService) PromoteDocker(params DockerPromoteParams) error {
	// Create URL
	restApi := path.Join("api/docker", params.SourceRepo, "v2", "promote")
	url, err := utils.BuildArtifactoryUrl(ps.GetArtifactoryDetails().GetUrl(), restApi, nil)
	if err != nil {
		return err
	}

	// Create body
	data := DockerPromoteBody{
		TargetRepo:             params.TargetRepo,
		DockerRepository:       params.SourceDockerImage,
		TargetDockerRepository: params.TargetDockerImage,
		Tag:                    params.SourceTag,
		TargetTag:              params.TargetTag,
		Copy:                   params.Copy,
	}
	requestContent, err := json.Marshal(data)
	if err != nil {
		return errorutils.CheckError(err)
	}

	// Send POST request
	httpClientsDetails := ps.GetArtifactoryDetails().CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := ps.client.SendPost(url, requestContent, &httpClientsDetails)
	if err != nil {
		return err
	}

	// Check results
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	log.Debug("Artifactory response: ", resp.Status)
	log.Info("Promoted image", params.SourceDockerImage, "to:", params.TargetRepo, "repository.")
	return nil
}

type DockerPromoteParams struct {
	// Mandatory:
	// The name of the source repository in Artifactory, e.g. "docker-local-1". Supported by local repositories only.
	SourceRepo string
	// The name of the target repository in Artifactory, e.g. "docker-local-2". Supported by local repositories only.
	TargetRepo string
	// The name of the source Docker image, e.g. "hello-world".
	SourceDockerImage string

	// Optional:
	// The name of the target Docker image, e.g "hello-world2". If not specified - will use the same name as 'SourceDockerImage'.
	TargetDockerImage string
	// The name of the source image tag. If not specified - the entire docker repository will be promoted.
	SourceTag string
	// The name of the target image tag. If not specified - will use the same tag as 'SourceTag'.
	TargetTag string
	// If set to true, will do copy instead of move.
	Copy bool
}

func (dp *DockerPromoteParams) GetTargetRepo() string {
	return dp.TargetRepo
}

func (dp *DockerPromoteParams) GetSourceDockerImage() string {
	return dp.SourceDockerImage
}

func (dp *DockerPromoteParams) GetTargetDockerRepository() string {
	return dp.TargetDockerImage
}

func (dp *DockerPromoteParams) GetSourceTag() string {
	return dp.SourceTag
}

func (dp *DockerPromoteParams) GetTargetTag() string {
	return dp.TargetTag
}

func (dp *DockerPromoteParams) IsCopy() bool {
	return dp.Copy
}

func NewDockerPromoteParams(sourceDockerImage, sourceRepo, targetRepo string) DockerPromoteParams {
	return DockerPromoteParams{
		SourceDockerImage: sourceDockerImage,
		SourceRepo:        sourceRepo,
		TargetRepo:        targetRepo,
	}
}

type DockerPromoteBody struct {
	TargetRepo             string `json:"targetRepo"`
	DockerRepository       string `json:"dockerRepository"`
	TargetDockerRepository string `json:"targetDockerRepository"`
	Tag                    string `json:"tag"`
	TargetTag              string `json:"targetTag"`
	Copy                   bool   `json:"copy"`
}

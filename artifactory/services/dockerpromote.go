package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DockerPromoteService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewDockerPromoteService(client *rthttpclient.ArtifactoryHttpClient) *DockerPromoteService {
	return &DockerPromoteService{client: client}
}

func (ps *DockerPromoteService) GetArtifactoryDetails() auth.ServiceDetails {
	return ps.ArtDetails
}

func (ps *DockerPromoteService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	ps.ArtDetails = rt
}

func (ps *DockerPromoteService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return ps.client, nil
}

func (ps *DockerPromoteService) IsDryRun() bool {
	return false
}

func (ps *DockerPromoteService) PromoteDocker(params DockerPromoteParams) error {
	// Create URL
	restApi := path.Join("api/docker", params.SourceRepo, "v2", "promote")
	url, err := utils.BuildArtifactoryUrl(ps.ArtDetails.GetUrl(), restApi, nil)
	if err != nil {
		return err
	}

	// Create body
	data := DockerPromoteBody{
		TargetRepo:             params.TargetRepo,
		DockerRepository:       params.DockerRepository,
		TargetDockerRepository: params.TargetDockerRepository,
		Tag:                    params.Tag,
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
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response: ", resp.Status)

	return nil
}

type DockerPromoteParams struct {
	SourceRepo             string
	TargetRepo             string
	DockerRepository       string
	TargetDockerRepository string
	Tag                    string
	TargetTag              string
	Copy                   bool
}

func (dp *DockerPromoteParams) GetTargetRepo() string {
	return dp.TargetRepo
}

func (dp *DockerPromoteParams) GetDockerRepository() string {
	return dp.DockerRepository
}

func (dp *DockerPromoteParams) GetTargetDockerRepository() string {
	return dp.TargetDockerRepository
}

func (dp *DockerPromoteParams) GetTag() string {
	return dp.Tag
}

func (dp *DockerPromoteParams) GetTargetTag() string {
	return dp.TargetTag
}

func (dp *DockerPromoteParams) IsCopy() bool {
	return dp.Copy
}

func NewDockerPromoteParams(sourceRepo, targetRepo, dockerRepository string) DockerPromoteParams {
	return DockerPromoteParams{
		SourceRepo:       sourceRepo,
		TargetRepo:       targetRepo,
		DockerRepository: dockerRepository,
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

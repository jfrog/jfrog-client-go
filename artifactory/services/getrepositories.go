package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type GetRepositoriesService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewGetRepositoriesService(client *rthttpclient.ArtifactoryHttpClient) *GetRepositoriesService {
	return &GetRepositoriesService{client: client}
}

func (grs *GetRepositoriesService) Get(repoKey string) (*RepositoryDetails, error) {
	httpClientsDetails := grs.ArtDetails.CreateHttpClientDetails()
	log.Info("Getting repository '" + repoKey + "' details ...")
	repoDetails := &RepositoryDetails{}
	resp, body, _, err := grs.client.SendGet(grs.ArtDetails.GetUrl()+"api/repositories/"+repoKey, true, &httpClientsDetails)
	if err != nil {
		return &RepositoryDetails{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &RepositoryDetails{}, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	if err := json.Unmarshal(body, &repoDetails); err != nil {
		return repoDetails, errorutils.CheckError(err)
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done getting repositories.")
	return repoDetails, nil
}

func (grs *GetRepositoriesService) GetAll() ([]*RepositoryDetails, error) {
	httpClientsDetails := grs.ArtDetails.CreateHttpClientDetails()
	log.Info("Getting repositories details ...")
	repoDetails := []*repositoriesDetailsBody{}
	resp, body, _, err := grs.client.SendGet(grs.ArtDetails.GetUrl()+"api/repositories", true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	if err := json.Unmarshal(body, &repoDetails); err != nil {
		return nil, errorutils.CheckError(err)
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done getting repository details.")
	return CreateRepositoryDetailsFromBody(repoDetails), nil
}

func CreateRepositoryDetailsFromBody(arg []*repositoriesDetailsBody) (result []*RepositoryDetails) {
	for i := 0; i < len(arg); i++ {
		result = append(result, &RepositoryDetails{
			Key:         arg[i].Key,
			Rclass:      arg[i].Rclass,
			Description: arg[i].Description,
			Url:         arg[i].Url,
			PackageType: arg[i].PackageType,
		})
	}
	return
}

type RepositoryDetails struct {
	Key          string
	Rclass       string
	Description  string
	Url          string
	PackageType  string
	Repositories []string
}

type repositoriesDetailsBody struct {
	Key         string `json:"key"`
	Rclass      string `json:"type"`
	Description string `json:"description"`
	Url         string `json:"url"`
	PackageType string `json:"packageType"`
}

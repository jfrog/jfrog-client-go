package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const apiRepositories = "api/repositories"

type GetRepositoryService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewGetRepositoryService(client *jfroghttpclient.JfrogHttpClient) *GetRepositoryService {
	return &GetRepositoryService{client: client}
}

func (grs *GetRepositoryService) Get(repoKey string) (*RepositoryDetails, error) {
	log.Info("Getting repository '" + repoKey + "' details ...")
	body, err := grs.sendGet(apiRepositories + "/" + repoKey)
	if err != nil {
		return nil, err
	}
	repoDetails := &RepositoryDetails{}
	if err := json.Unmarshal(body, &repoDetails); err != nil {
		return repoDetails, errorutils.CheckError(err)
	}
	return repoDetails, nil
}

func (grs *GetRepositoryService) GetAll() (*[]RepositoryDetails, error) {
	log.Info("Getting all repositories ...")
	body, err := grs.sendGet(apiRepositories)
	if err != nil {
		return nil, err
	}
	repoDetails := &[]RepositoryDetails{}
	if err := json.Unmarshal(body, &repoDetails); err != nil {
		return repoDetails, errorutils.CheckError(err)
	}
	return repoDetails, nil
}

func (grs *GetRepositoryService) sendGet(api string) ([]byte, error) {
	httpClientsDetails := grs.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := grs.client.SendGet(grs.ArtDetails.GetUrl()+api, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done getting repository details.")
	return body, nil
}

type RepositoryDetails struct {
	Key         string
	Rclass      string
	Description string
	Url         string
	PackageType string
}

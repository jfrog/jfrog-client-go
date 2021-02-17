package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const apiRepositories = "api/repositories"

type RepositoriesService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewRepositoriesService(client *jfroghttpclient.JfrogHttpClient) *RepositoriesService {
	return &RepositoriesService{client: client}
}

func (rs *RepositoriesService) Get(repoKey string) (*RepositoryDetails, error) {
	log.Info("Getting repository '" + repoKey + "' details ...")
	body, err := rs.sendGet(apiRepositories + "/" + repoKey)
	if err != nil {
		return nil, err
	}
	repoDetails := &RepositoryDetails{}
	if err := json.Unmarshal(body, &repoDetails); err != nil {
		return repoDetails, errorutils.CheckError(err)
	}
	return repoDetails, nil
}

func (rs *RepositoriesService) GetAll() (*[]RepositoryDetails, error) {
	log.Info("Getting all repositories ...")
	return rs.GetAllFromTypeAndPackage(RepositoriesFilterParams{RepoType: "", PackageType: ""})
}

func (rs *RepositoriesService) GetAllFromTypeAndPackage(params RepositoriesFilterParams) (*[]RepositoryDetails, error) {
	url := fmt.Sprintf("%s?type=%s&packageType=%s", apiRepositories, params.RepoType, params.PackageType)
	body, err := rs.sendGet(url)
	if err != nil {
		return nil, err
	}
	repoDetails := &[]RepositoryDetails{}
	if err := json.Unmarshal(body, &repoDetails); err != nil {
		return repoDetails, errorutils.CheckError(err)
	}
	return repoDetails, nil
}

func (rs *RepositoriesService) sendGet(api string) ([]byte, error) {
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := rs.client.SendGet(rs.ArtDetails.GetUrl()+api, true, &httpClientsDetails)
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
	Type        string
	Description string
	Url         string
	PackageType string
}

func (rd RepositoryDetails) getRepoType() string {
	// When getting All repos from artifactory the REST returns with Type field,
	// but when getting a specific repo it will return with the Rclass field.
	if rd.Rclass != "" {
		return rd.Rclass
	}
	return rd.Type
}

type RepositoriesFilterParams struct {
	RepoType    string
	PackageType string
}

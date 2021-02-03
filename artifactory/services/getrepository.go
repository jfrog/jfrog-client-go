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
	return grs.GetAllFromTypeAndPackage("", "")
}

func (grs *GetRepositoryService) GetAllFromTypeAndPackage(repoType, packageType string) (*[]RepositoryDetails, error) {
	url := fmt.Sprintf("%s?type=%s&packageType=%s", apiRepositories, repoType, packageType)
	body, err := grs.sendGet(url)
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

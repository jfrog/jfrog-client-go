package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
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

// Get fetches repository details from Artifactory using given repokey (name) into the given params struct.
// The function expect to get the repo key and a pointer to a param struct that will be filled up.
// The param struct should contains the desired param's fields corresponded to the Artifactory REST API, such as RepositoryDetails, LocalRepositoryBaseParams, etc.
func (rs *RepositoriesService) Get(repoKey string, repoDetails interface{}) error {
	log.Info("Getting repository '" + repoKey + "' details ...")
	body, err := rs.sendGet(apiRepositories + "/" + repoKey)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, repoDetails)
	return errorutils.CheckError(err)
}

func (rs *RepositoriesService) GetAll() (*[]RepositoryDetails, error) {
	log.Info("Getting all repositories ...")
	return rs.GetWithFilter(RepositoriesFilterParams{RepoType: "", PackageType: ""})
}

func (rs *RepositoriesService) GetWithFilter(params RepositoriesFilterParams) (*[]RepositoryDetails, error) {
	url := fmt.Sprintf("%s?type=%s&packageType=%s", apiRepositories, params.RepoType, params.PackageType)
	body, err := rs.sendGet(url)
	if err != nil {
		return nil, err
	}
	repoDetails := &[]RepositoryDetails{}
	err = json.Unmarshal(body, &repoDetails)
	return repoDetails, errorutils.CheckError(err)
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
	log.Debug("Done getting repository details.")
	return body, nil
}

func (rs *RepositoriesService) CreateRemote(params RemoteRepositoryBaseParams) error {
	return rs.createRepo(params, params.Key)
}

func (rs *RepositoriesService) CreateVirtual(params VirtualRepositoryBaseParams) error {
	return rs.createRepo(params, params.Key)
}

func (rs *RepositoriesService) CreateLocal(params LocalRepositoryBaseParams) error {
	return rs.createRepo(params, params.Key)
}

func (rs *RepositoriesService) createRepo(params interface{}, repoName string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := rs.client.SendPut(rs.ArtDetails.GetUrl()+"api/repositories/"+repoName, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info(fmt.Sprintf("Repository %q created.", repoName))
	return nil
}

type RepositoryDetails struct {
	Key         string
	Rclass      string
	Type        string
	Description string
	Url         string
	PackageType string
}

func (rd RepositoryDetails) GetRepoType() string {
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

func NewRepositoriesFilterParams() RepositoriesFilterParams {
	return RepositoriesFilterParams{}
}

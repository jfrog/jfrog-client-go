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

type GetRepositoryService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewGetRepositoryService(client *rthttpclient.ArtifactoryHttpClient) *GetRepositoryService {
	return &GetRepositoryService{client: client}
}

func (grs *GetRepositoryService) Get(repoKey string) (*GetRepositoryData, error) {
	httpClientsDetails := grs.ArtDetails.CreateHttpClientDetails()
	log.Info("Getting repository '" + repoKey + "' details ...")
	repoDetails := &GetRepositoryData{}
	resp, body, _, err := grs.client.SendGet(grs.ArtDetails.GetUrl()+"api/repositories/"+repoKey, true, &httpClientsDetails)
	if err != nil {
		return &GetRepositoryData{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &GetRepositoryData{}, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	err = json.Unmarshal(body, &repoDetails)
	if resp.StatusCode != http.StatusOK {
		return &GetRepositoryData{}, errorutils.CheckError(err)
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done getting repository details.")
	return repoDetails, nil
}

type GetRepositoryData struct {
	Key         string
	Rclass      string
	Description string
	Url         string
	PackageType string
}

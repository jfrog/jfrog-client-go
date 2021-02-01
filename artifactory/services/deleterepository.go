package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DeleteRepositoryService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewDeleteRepositoryService(client *jfroghttpclient.JfrogHttpClient) *DeleteRepositoryService {
	return &DeleteRepositoryService{client: client}
}

func (drs *DeleteRepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return drs.client
}

func (drs *DeleteRepositoryService) Delete(repoKey string) error {
	httpClientsDetails := drs.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting repository " + repoKey + "...")
	resp, body, err := drs.client.SendDelete(drs.ArtDetails.GetUrl()+"api/repositories/"+repoKey, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done deleting repository.")
	return nil
}

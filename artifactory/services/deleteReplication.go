package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"

	httpclient "github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DeleteReplicationService struct {
	client     *httpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewDeleteReplicationService(client *httpclient.JfrogHttpClient) *DeleteReplicationService {
	return &DeleteReplicationService{client: client}
}

func (drs *DeleteReplicationService) GetJfrogHttpClient() *httpclient.JfrogHttpClient {
	return drs.client
}

func (drs *DeleteReplicationService) DeleteReplication(repoKey string) error {
	httpClientsDetails := drs.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting replication job...")
	resp, body, err := drs.client.SendDelete(drs.ArtDetails.GetUrl()+"api/replications/"+repoKey, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done Deleting replication job.")
	return nil
}

package services

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DeleteReplicationService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewDeleteReplicationService(client *jfroghttpclient.JfrogHttpClient) *DeleteReplicationService {
	return &DeleteReplicationService{client: client}
}

func (drs *DeleteReplicationService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return drs.client
}

func (drs *DeleteReplicationService) DeleteReplication(repoKey string) error {
	httpClientsDetails := drs.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting replication job...")
	resp, body, err := drs.client.SendDelete(drs.ArtDetails.GetUrl()+"api/replications/"+repoKey, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done Deleting replication job.")
	return nil
}

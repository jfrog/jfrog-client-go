package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DeleteMultipleReplicationService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewDeleteMultipleReplicationService(client *rthttpclient.ArtifactoryHttpClient) *DeleteMultipleReplicationService {
	return &DeleteMultipleReplicationService{client: client}
}

func (drs *DeleteMultipleReplicationService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return drs.client
}

func (drs *DeleteMultipleReplicationService) DeleteMultipleReplication(repoKey, repoUrl string) error {
	httpClientsDetails := drs.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting multiple replication job...")
	resp, body, err := drs.client.SendDelete(drs.ArtDetails.GetUrl()+"api/replications/"+repoKey+"?url="+repoUrl, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done Deleting multiple replication job.")
	return nil
}

package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ShowReplicationService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewShowReplicationService(client *rthttpclient.ArtifactoryHttpClient) *ShowReplicationService {
	return &ShowReplicationService{client: client}
}

func (drs *ShowReplicationService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return drs.client
}

func (drs *ShowReplicationService) Show(repoKey string) ([]PushReplicationParams, error) {
	httpClientsDetails := drs.ArtDetails.CreateHttpClientDetails()
	log.Info("Retrive replication configuration...")
	resp, body, _, err := drs.client.SendGet(drs.ArtDetails.GetUrl()+"api/replications/"+repoKey, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done retrive replication job.")
	var replicationConf []PushReplicationParams
	if err := json.Unmarshal(body, &replicationConf); err != nil {
		return nil, err
	}
	return replicationConf, nil
}

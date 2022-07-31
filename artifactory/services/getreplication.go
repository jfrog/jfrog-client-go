package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type GetReplicationService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewGetReplicationService(client *jfroghttpclient.JfrogHttpClient) *GetReplicationService {
	return &GetReplicationService{client: client}
}

func (drs *GetReplicationService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return drs.client
}

func (drs *GetReplicationService) GetReplication(repoKey string) ([]utils.ReplicationParams, error) {
	body, err := drs.preform(repoKey)
	if err != nil {
		return nil, err
	}
	var replicationBody []utils.GetReplicationBody
	if err := json.Unmarshal(body, &replicationBody); err != nil {
		return nil, errorutils.CheckError(err)
	}

	var replicationConf = make([]utils.ReplicationParams, len(replicationBody))
	for i, body := range replicationBody {
		replicationConf[i] = *utils.CreateReplicationParams(body)
	}

	return replicationConf, nil
}

func (drs *GetReplicationService) preform(repoKey string) ([]byte, error) {
	httpClientsDetails := drs.ArtDetails.CreateHttpClientDetails()
	log.Info("Retrieve replication configuration...")
	resp, body, _, err := drs.client.SendGet(drs.ArtDetails.GetUrl()+"api/replications/"+repoKey, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done retrieve replication job.")
	return body, nil
}

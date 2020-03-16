package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ReplicationService struct {
	isUpdate   bool
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewReplicationService(client *rthttpclient.ArtifactoryHttpClient, isUpdate bool) *ReplicationService {
	return &ReplicationService{client: client, isUpdate: isUpdate}
}

func (rs *ReplicationService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return rs.client
}

func (rs *ReplicationService) performRequest(replicationParams []byte, repoKey string) error {
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.replications.ReplicationConfigRequest+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/replications/" + repoKey
	var operationString string
	var resp *http.Response
	var body []byte
	var err error
	if rs.isUpdate {
		log.Info("Update replication job...")
		operationString = "updating"
		resp, body, err = rs.client.SendPost(url, replicationParams, &httpClientsDetails)
	} else {
		log.Info("Creating replication job...")
		operationString = "creating"
		resp, body, err = rs.client.SendPut(url, replicationParams, &httpClientsDetails)
	}
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done " + operationString + " repository.")
	return nil
}

func (rs *ReplicationService) Push(params PushReplicationParams) error {
	content, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return rs.performRequest(content, params.RepoKey)
}

func (rs *ReplicationService) Pull(params PullReplicationParams) error {
	content, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return rs.performRequest(content, params.RepoKey)
}

type CommonReplicationParams struct {
	CronExp                string `json:"cronExp"`
	RepoKey                string `json:"repoKey"`
	EnableEventReplication bool   `json:"enableEventReplication"`
	SocketTimeoutMillis    int    `json:"socketTimeoutMillis"`
	Enabled                bool   `json:"enabled"`
	SyncDeletes            bool   `json:"syncDeletes"`
	SyncProperties         bool   `json:"syncProperties"`
	SyncStatistics         bool   `json:"syncStatistics"`
	PathPrefix             string `json:"pathPrefix"`
}

type PullReplicationParams struct {
	CommonReplicationParams
}

type PushReplicationParams struct {
	CommonReplicationParams
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

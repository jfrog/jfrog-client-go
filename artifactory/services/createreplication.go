package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	httpclient "github.com/jfrog/jfrog-client-go/httpclient/jfrog"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type CreateReplicationService struct {
	client     *httpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewCreateReplicationService(client *httpclient.JfrogHttpClient) *CreateReplicationService {
	return &CreateReplicationService{client: client}
}

func (rs *CreateReplicationService) GetJfrogHttpClient() *httpclient.JfrogHttpClient {
	return rs.client
}

func (rs *CreateReplicationService) performRequest(params *utils.ReplicationBody) error {
	content, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.replications.ReplicationConfigRequest+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/replications/" + params.RepoKey
	var resp *http.Response
	var body []byte
	log.Info("Creating replication..")
	operationString := "creating"
	resp, body, err = rs.client.SendPut(url, content, &httpClientsDetails)
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

func (rs *CreateReplicationService) CreateReplication(params CreateReplicationParams) error {
	return rs.performRequest(utils.CreateReplicationBody(params.ReplicationParams))
}

func NewCreateReplicationParams() CreateReplicationParams {
	return CreateReplicationParams{}
}

type CreateReplicationParams struct {
	utils.ReplicationParams
}

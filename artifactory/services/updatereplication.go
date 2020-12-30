package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type UpdateReplicationService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewUpdateReplicationService(client *jfroghttpclient.JfrogHttpClient) *UpdateReplicationService {
	return &UpdateReplicationService{client: client}
}

func (rs *UpdateReplicationService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return rs.client
}

func (rs *UpdateReplicationService) performRequest(params *utils.ReplicationBody) error {
	content, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.replications.ReplicationConfigRequest+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/replications/" + params.RepoKey
	var resp *http.Response
	var body []byte
	log.Info("Update replication...")
	operationString := "updating"
	resp, body, err = rs.client.SendPost(url, content, &httpClientsDetails)
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

func (rs *UpdateReplicationService) UpdateReplication(params UpdateReplicationParams) error {
	return rs.performRequest(utils.CreateReplicationBody(params.ReplicationParams))
}

func NewUpdateReplicationParams() UpdateReplicationParams {
	return UpdateReplicationParams{}
}

type UpdateReplicationParams struct {
	utils.ReplicationParams
}

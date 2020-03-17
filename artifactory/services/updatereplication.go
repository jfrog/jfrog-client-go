package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type UpdateReplicationService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.CommonDetails
}

func NewUpdateReplicationService(client *rthttpclient.ArtifactoryHttpClient, isUpdate bool) *UpdateReplicationService {
	return &UpdateReplicationService{client: client}
}

func (rs *UpdateReplicationService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return rs.client
}

func (rs *UpdateReplicationService) performRequest(params utils.ReplicationParams) error {
	content, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.replications.ReplicationConfigRequest+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/replications/" + params.RepoKey
	var operationString string
	var resp *http.Response
	var body []byte
	log.Info("Update replication...")
	operationString = "updating"
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

func (rs *UpdateReplicationService) UpdateReplication(params utils.ReplicationParams) error {
	return rs.performRequest(params)
}

func NewUpdateReplicationParams() utils.ReplicationParams {
	return utils.ReplicationParams{}
}

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

type CreateMultipleReplicationService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewCreateMultipleReplicationService(client *rthttpclient.ArtifactoryHttpClient) *CreateMultipleReplicationService {
	return &CreateMultipleReplicationService{client: client}
}

func (rs *CreateMultipleReplicationService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return rs.client
}

func (rs *CreateMultipleReplicationService) performRequest(params *utils.MultipleReplicationBody) error {
	content, err := json.Marshal(params)
	if err != nil {
		return errorutils.CheckError(err)
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.replications.MultipleReplicationConfigRequest+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/replications/multiple/" + params.RepoKey
	var resp *http.Response
	var body []byte
	log.Info("Creating multple replications..")
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

func (rs *CreateMultipleReplicationService) CreateMultipleReplication(params CreateMultipleReplicationParams) error {
	return rs.performRequest(utils.CreateMultipleReplicationBody(params.MultipleReplicationParams))
}

func NewCreateMultipleReplicationParams() CreateMultipleReplicationParams {
	return CreateMultipleReplicationParams{}
}

type CreateMultipleReplicationParams struct {
	utils.MultipleReplicationParams
}

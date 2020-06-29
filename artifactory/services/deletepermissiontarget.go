package services

import (
	"errors"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type DeletePermissionTargetService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewDeletePermissionTargetService(client *rthttpclient.ArtifactoryHttpClient) *DeletePermissionTargetService {
	return &DeletePermissionTargetService{client: client}
}

func (dpts *DeletePermissionTargetService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return dpts.client
}

func (dpts *DeletePermissionTargetService) Delete(permissionTargetName string) error {
	httpClientsDetails := dpts.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting permission target...")
	resp, body, err := dpts.client.SendDelete(dpts.ArtDetails.GetUrl()+"api/security/permissions/"+permissionTargetName, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done deleting permission target.")
	return nil
}

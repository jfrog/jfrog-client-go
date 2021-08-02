package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type FederationService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewFederationService(client *jfroghttpclient.JfrogHttpClient) *FederationService {
	return &FederationService{client: client}
}

func (fs *FederationService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	fs.ArtDetails = rt
}

func (fs *FederationService) ConvertLocalToFederated(repoKey string) error {
	httpClientsDetails := fs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = fs.ArtDetails.GetUrl() + "api/federation/migrate/" + repoKey
	log.Info("Converting local repository to federated repository...")
	resp, body, err := fs.client.SendPost(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done converting repository.")
	return nil
}

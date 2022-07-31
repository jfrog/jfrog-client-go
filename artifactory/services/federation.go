package services

import (
	"net/http"
	"net/url"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
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
	var url = fs.ArtDetails.GetUrl() + "api/federation/migrate/" + url.PathEscape(repoKey)
	log.Info("Converting local repository to federated repository...")
	resp, body, err := fs.client.SendPost(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done converting repository.")
	return nil
}

func (fs *FederationService) TriggerFederatedFullSyncAll(repoKey string) error {
	httpClientsDetails := fs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = fs.ArtDetails.GetUrl() + "api/federation/fullSync/" + url.PathEscape(repoKey)
	log.Info("Triggering full federated repository synchronisation...")
	resp, body, err := fs.client.SendPost(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done triggering full federated repository synchronisation.")
	return nil
}

func (fs *FederationService) TriggerFederatedFullSyncMirror(repoKey string, mirrorUrl string) error {
	httpClientsDetails := fs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = fs.ArtDetails.GetUrl() + "api/federation/fullSync/" + url.PathEscape(repoKey) + "?mirror=" + url.QueryEscape(mirrorUrl)
	log.Info("Triggering federated repository synchronisation...")
	resp, body, err := fs.client.SendPost(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done triggering federated repository synchronisation.")
	return nil
}

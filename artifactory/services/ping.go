package services

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type PingService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewPingService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *PingService {
	return &PingService{artDetails: &artDetails, client: client}
}

func (ps *PingService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ps.artDetails
}

func (ps *PingService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ps.client
}

func (ps *PingService) IsDryRun() bool {
	return false
}

func (ps *PingService) Ping() ([]byte, error) {
	url, err := utils.BuildArtifactoryUrl(ps.GetArtifactoryDetails().GetUrl(), "api/system/ping", nil)
	if err != nil {
		return nil, err
	}
	httpClientDetails := ps.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := ps.client.SendGet(url, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return body, err
	}
	log.Debug("Artifactory response: ", resp.Status)
	return body, nil
}

package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type PingService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewPingService(client *jfroghttpclient.JfrogHttpClient) *PingService {
	return &PingService{client: client}
}

func (ps *PingService) GetArtifactoryDetails() auth.ServiceDetails {
	return ps.ArtDetails
}

func (ps *PingService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	ps.ArtDetails = rt
}

func (ps *PingService) GetJfrogHttpClient() (*jfroghttpclient.JfrogHttpClient, error) {
	return ps.client, nil
}

func (ps *PingService) IsDryRun() bool {
	return false
}

func (ps *PingService) Ping() ([]byte, error) {
	url, err := utils.BuildArtifactoryUrl(ps.ArtDetails.GetUrl(), "api/system/ping", nil)
	if err != nil {
		return nil, err
	}
	httpClientDetails := ps.ArtDetails.CreateHttpClientDetails()
	resp, respBody, _, err := ps.client.SendGet(url, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return respBody, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status))
	}
	log.Debug("Artifactory response: ", resp.Status)
	return respBody, nil
}

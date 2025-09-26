package services

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	jpdService "github.com/jfrog/jfrog-client-go/jpd"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	releaseBundlesAPI = "lifecycle/api/v2/release_bundle/names"
)

type ReleaseBundlesStatsService struct {
	client    *jfroghttpclient.JfrogHttpClient
	lcDetails *auth.ServiceDetails
}

func NewReleaseBundlesStatsService(lcDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *ReleaseBundlesStatsService {
	return &ReleaseBundlesStatsService{lcDetails: &lcDetails, client: client}
}

func (rbss *ReleaseBundlesStatsService) GetLifecycleDetails() auth.ServiceDetails {
	return *rbss.lcDetails
}

func (rbss *ReleaseBundlesStatsService) GetReleaseBundlesStats(serverUrl string) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, releaseBundlesAPI, nil)
	if err != nil {
		wrappedError := fmt.Errorf("failed to build RELEASE-BUNDLES API Endpoint: %w", err)
		return nil, jpdService.NewGenericError("RELEASE-BUNDLES", wrappedError)
	}
	httpClientsDetails := rbss.GetLifecycleDetails().CreateHttpClientDetails()
	resp, body, _, err := rbss.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		wrappedError := fmt.Errorf("failed to call RELEASE-BUNDLES API: %w", err)
		return nil, jpdService.NewGenericError("RELEASE-BUNDLES", wrappedError)
	}
	log.Debug("Release Bundle API response:", resp.Status)
	return body, err
}

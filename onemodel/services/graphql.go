package services

import (
	"fmt"
	"net/http"
	"net/url"

	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const queryUrl = "api/v1/graphql"

type Service interface {
	Query(query []byte) ([]byte, error)
}

type onemodelService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails *auth.ServiceDetails
}

func NewOnemodelService(serviceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) Service {
	return &onemodelService{serviceDetails: &serviceDetails, client: client}
}

func (m *onemodelService) GetOnemodelDetails() auth.ServiceDetails {
	return *m.serviceDetails
}

func (m *onemodelService) Query(query []byte) ([]byte, error) {
	graphqlUrl, err := url.Parse(m.GetOnemodelDetails().GetUrl() + queryUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	httpClientDetails := m.GetOnemodelDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)

	resp, body, err := m.client.SendPost(graphqlUrl.String(), query, &httpClientDetails)
	if err != nil {
		return []byte{}, err
	}
	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

package services

import (
	"fmt"
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"net/url"
)

const queryUrl = "api/v1/query"

type Service interface {
	Query(query []byte) ([]byte, error)
}

type metadataService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails *auth.ServiceDetails
}

func NewMetadataService(serviceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) Service {
	return &metadataService{serviceDetails: &serviceDetails, client: client}
}

func (m *metadataService) GetMetadataDetails() auth.ServiceDetails {
	return *m.serviceDetails
}

func (m *metadataService) Query(query []byte) ([]byte, error) {
	graphqlUrl, err := url.Parse(m.GetMetadataDetails().GetUrl() + queryUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	httpClientDetails := m.GetMetadataDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)

	resp, body, err := m.client.SendPost(graphqlUrl.String(), query, &httpClientDetails)
	if err != nil {
		return []byte{}, err
	}
	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

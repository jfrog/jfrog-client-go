package services

import (
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"net/url"
)

type Service interface {
	Query(query QueryDetails) ([]byte, error)
}

type metadataService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails *auth.ServiceDetails
}

func NewMetadataService(serviceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) Service {
	return &metadataService{serviceDetails: &serviceDetails, client: client}
}

func (m *metadataService) GetMetaDataDetails() auth.ServiceDetails {
	return *m.serviceDetails
}

func (m *metadataService) Query(query QueryDetails) ([]byte, error) {
	graphiqlUrl, err := url.Parse(m.GetMetaDataDetails().GetUrl() + "api/v1/query")
	if err != nil {
		return []byte{}, errorutils.CheckError(err)
	}

	httpClientDetails := m.GetMetaDataDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)

	resp, body, err := m.client.SendPost(graphiqlUrl.String(), query.Body, &httpClientDetails)
	if err != nil {
		return []byte{}, err
	}
	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

type QueryDetails struct {
	Body []byte
}

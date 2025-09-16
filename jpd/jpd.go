package jpd

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	JPDsAPI = "mc/api/v1/jpds"
)

type JPDsStatsService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
	ServerUrl  string
}

type GenericError struct {
	Product string
	Err     string
}

func NewJPDsStatsService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *JPDsStatsService {
	return &JPDsStatsService{client: client, ArtDetails: artDetails}
}

func (g GenericError) Error() string {
	return fmt.Sprintf("failed to get stats for '%s': %s", g.Product, g.Err)
}

func (s *JPDsStatsService) SetServerUrl(serverUrl string) {
	s.ServerUrl = serverUrl
}
func (s *JPDsStatsService) GetServerUrl() string {
	return s.ServerUrl
}

func NewGenericError(product string, err string) *GenericError {
	return &GenericError{
		Product: product,
		Err:     err,
	}
}

func (ss *JPDsStatsService) GetAllJPDs(serverUrl string) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, JPDsAPI, nil)
	if err != nil {
		return nil, err
	}
	httpClientsDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, NewGenericError("JPD", err.Error())
	}
	log.Debug("JPDs API response:", resp.Status)
	return body, err
}

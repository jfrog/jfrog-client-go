package jpd

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
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
	Product string `write:"-"`
	Err     error  `write:"Error"`
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

func NewGenericError(product string, err error) *GenericError {
	return &GenericError{
		Product: product,
		Err:     err,
	}
}

func (ss *JPDsStatsService) GetAllJPDs(serverUrl string) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, JPDsAPI, nil)
	if err != nil {
		wrappedError := fmt.Errorf("failed to build JPD API: %w", err)
		return nil, NewGenericError("JPDs", wrappedError)
	}
	httpClientsDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		wrappedError := fmt.Errorf("failed to call JPD API: %w", err)
		return nil, NewGenericError("JPDs", wrappedError)
	}
	if resp.StatusCode != http.StatusOK {
		wrappedError := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return nil, NewGenericError("JPDs", wrappedError)
	}
	log.Debug("JPDs API response:", resp.Status)
	return body, err
}

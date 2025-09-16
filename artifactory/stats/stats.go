package statsservice

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const (
	JPDsAPI           = "mc/api/v1/jpds"
	releaseBundlesAPI = "lifecycle/api/v2/release_bundle/names"
	repositoriesAPI   = "artifactory/api/repositories"
	tokensAPI         = "access/api/v1/tokens"
)

type StatsService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
	ServerUrl  string
}

type APIError struct {
	Product    string
	StatusCode int
	StatusText string
	Suggestion string
}

type GenericError struct {
	Product string
	Err     string
}

func NewStatsService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *StatsService {
	return &StatsService{client: client, ArtDetails: artDetails}
}

func (g GenericError) Error() string {
	return fmt.Sprintf("failed to get stats for '%s': %s", g.Product, g.Err)
}

func (e *APIError) Error() string {
	return fmt.Sprintf("failed to get stats for '%s': %s. %s", e.Product, e.StatusText, e.Suggestion)
}

func (s *StatsService) SetServerUrl(serverUrl string) {
	s.ServerUrl = serverUrl
}
func (s *StatsService) GetServerUrl() string {
	return s.ServerUrl
}

func NewFailedRequestError(statusCode int, statusText string, product string) *APIError {
	var details string
	switch {
	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		details = "Need Valid Token"
	case statusCode == http.StatusNotFound:
		details = "Resource Not Found"
	case statusCode >= 400 && statusCode < 500:
		details = "Client Error"
	case statusCode >= 500 && statusCode < 600:
		details = "Server Error"
	}
	return &APIError{
		Product:    product,
		StatusCode: statusCode,
		StatusText: statusText,
		Suggestion: details,
	}
}

func NewGenericError(product string, err string) *GenericError {
	return &GenericError{
		Product: product,
		Err:     err,
	}
}

func (ss *StatsService) GetRepositoriesStats(serverUrl string) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, repositoriesAPI, nil)
	if err != nil {
		return nil, err
	}
	httpClientsDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, NewGenericError("REPOSITORY", err.Error())
	}
	log.Debug("Repositories API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "REPOSITORY")
		return nil, err
	}
	return body, err
}

func (ss *StatsService) GetJPDsStats(serverUrl string) ([]byte, error) {
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
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "JPD")
		return nil, err
	}
	return body, err
}

func (ss *StatsService) GetReleaseBundlesStats(serverUrl string) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, releaseBundlesAPI, nil)
	if err != nil {
		return nil, err
	}
	httpClientsDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, NewGenericError("RELEASE-BUNDLES", err.Error())
	}
	log.Debug("Release Bundle API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "RELEASE-BUNDLES")
		return nil, err
	}
	return body, err
}

func (ss *StatsService) GetTokenDetails(serverUrl string, tokenId string) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, tokensAPI+"/"+tokenId, nil)
	if err != nil {
		return nil, err
	}
	httpClientsDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "TOKEN")
		return nil, err
	}
	return body, err
}

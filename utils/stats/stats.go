package clientStats

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const (
	artifactoryStatsAPI = "artifactory/api/storageinfo"
	xrayPolicyAPI       = "xray/api/v2/policies"
	xrayWatchesAPI      = "xray/api/v2/watches"
	projectsAPI         = "access/api/v1/projects"
	JPDsAPI             = "mc/api/v1/jpds"
	releaseBundlesAPI   = "lifecycle/api/v2/release_bundle/names"
	repositoriesAPI     = "artifactory/api/repositories"
	tokensAPI           = "access/api/v1/tokens"
)

type APIError struct {
	Product    string
	StatusCode int
	StatusText string
	Suggestion string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("failed to get stats for '%s': %s. %s", e.Product, e.StatusText, e.Suggestion)
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

func GetArtifactoryStats(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, artifactoryStatsAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	log.Debug("Artifactory API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "ARTIFACTORY")
		return nil, err
	}
	return body, err
}

func GetRepositoriesStats(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, repositoriesAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	log.Debug("Repositories API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "REPOSITORY")
		return nil, err
	}
	return body, err
}

func GetXrayPolicies(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, xrayPolicyAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	log.Debug("Xray Policy API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "POLICIES")
		return nil, err
	}
	return body, err
}

func GetXrayWatches(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, xrayWatchesAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	if err != nil {
		return nil, err
	}
	log.Debug("Xray Watch API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "WATCHES")
		return nil, err
	}
	return body, err
}

func GetProjectsStats(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, projectsAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	log.Debug("Projects API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "PROJECTS")
		return nil, err
	}
	return body, err
}

func GetJPDsStats(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, JPDsAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	log.Debug("JPDs API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "JPD")
		return nil, err
	}
	return body, err
}

func GetReleaseBundlesStats(client *httpclient.HttpClient, serverUrl string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(serverUrl, releaseBundlesAPI, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	log.Debug("Release Bundle API response:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "RELEASE-BUNDLES")
		return nil, err
	}
	return body, err
}

func GetTokenDetails(client *httpclient.HttpClient, baseUrl string, tokenId string, httpClientDetails httputils.HttpClientDetails) ([]byte, error) {
	requestFullUrl, err := utils.BuildUrl(baseUrl, tokensAPI+"/"+tokenId, nil)
	if err != nil {
		return nil, err
	}
	resp, body, _, err := client.SendGet(requestFullUrl, true, httpClientDetails, "")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := NewFailedRequestError(resp.StatusCode, resp.Status, "TOKEN")
		return nil, err
	}
	return body, err
}

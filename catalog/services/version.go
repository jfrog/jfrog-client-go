package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	catalogPingApi    = "api/v1/system/ping"
	catalogVersionApi = "api/v1/system/version"

	catalogMinVersionForEnrichApi = "1.0.0"
)

type VersionService struct {
	client         *jfroghttpclient.JfrogHttpClient
	CatalogDetails auth.ServiceDetails
}

// NewVersionService creates a new service to retrieve the version of Catalog
func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

// Catalog currently does not have a version endpoint, so we use the ping endpoint to check if the service is up and running.
// https://jfrog-int.atlassian.net/browse/CTLG-829
func (vs *VersionService) GetVersion() (string, error) {
	versionResponse, err := vs.getVersion()
	if err == nil {
		return versionResponse.Version, nil
	}
	// Since Catalog did not have a version endpoint, at the past, try ping endpoint is used to verify connectivity.
	return catalogMinVersionForEnrichApi, vs.Ping()
}

func (vs *VersionService) getVersion() (VersionResponse, error) {
	var versionResponse VersionResponse
	httpDetails := vs.CatalogDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.CatalogDetails.GetUrl()+catalogVersionApi, true, &httpDetails)
	if err != nil {
		return versionResponse, errors.New("failed while attempting to get JFrog Catalog version: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return versionResponse, errors.New("got unexpected server response while attempting to get JFrog Catalog version:\n" + err.Error())
	}
	if err = errorutils.CheckError(json.Unmarshal(body, &versionResponse)); err != nil {
		return versionResponse, errors.New("failed to parse version response from JFrog Catalog: " + err.Error())
	}
	return versionResponse, nil
}

func (vs *VersionService) Ping() error {
	httpDetails := vs.CatalogDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.CatalogDetails.GetUrl()+catalogPingApi, true, &httpDetails)
	if err != nil {
		return errors.New("failed while attempting to ping JFrog Catalog: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return errors.New("got unexpected server response while attempting to ping JFrog Catalog:\n" + err.Error())
	}
	return nil
}

type VersionResponse struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	BuildDate string `json:"build_date"`
}

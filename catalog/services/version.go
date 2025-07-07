package services

import (
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	catalogPingApi = "api/v1/system/ping"
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
	httpDetails := vs.CatalogDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.CatalogDetails.GetUrl()+catalogPingApi, true, &httpDetails)
	if err != nil {
		return "", errors.New("failed while attempting to ping JFrog Catalog: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", errors.New("got unexpected server response while attempting to ping JFrog Catalog:\n" + err.Error())
	}
	// Since Catalog does not have a version endpoint, we return a hardcoded version.
	return catalogMinVersionForEnrichApi, nil
}

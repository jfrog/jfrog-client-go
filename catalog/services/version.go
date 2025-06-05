package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	catalogVersionApi          = "api/v1/system/version"
)

// VersionService returns the https client and Catalog details
type VersionService struct {
	client      *jfroghttpclient.JfrogHttpClient
	CatalogDetails auth.ServiceDetails
}

// NewVersionService creates a new service to retrieve the version of Catalog
func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

// GetVersion returns the version of Xray
func (vs *VersionService) GetVersion() (string, error) {
	httpDetails := vs.CatalogDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.CatalogDetails.GetUrl() + catalogVersionApi, true, &httpDetails)
	if err != nil {
		return "", errors.New("failed while attempting to get JFrog Catalog version: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", errors.New("got unexpected server response while attempting to get JFrog Catalog version:\n" + err.Error())
	}
	var version xrayVersion
	if err = json.Unmarshal(body, &version); err != nil {
		return "", errorutils.CheckErrorf("couldn't parse JFrog Catalog server version response: %s", err.Error())
	}
	return strings.TrimSpace(version.Version), nil
}

type xrayVersion struct {
	Version  string `json:"xray_version,omitempty"`
	Revision string `json:"xray_revision,omitempty"`
}

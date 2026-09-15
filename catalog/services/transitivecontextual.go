package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const transitiveContextualApi = "api/v1/dependencies/contextual"

type Ecosystem string

const (
	GenericEcosystem Ecosystem = "generic"
	DebianEcosystem  Ecosystem = "debian"
	UbuntuEcosystem  Ecosystem = "ubuntu"
)

// PackageVersionKey identifies a package by type, name, namespace, version and ecosystem.
type PackageVersionKey struct {
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Version   string    `json:"version"`
	Ecosystem Ecosystem `json:"ecosystem"`
}

type transitiveContextualRequest struct {
	Cves     []string            `json:"cves"`
	Packages []PackageVersionKey `json:"packages"`
}

// TransitiveContextualPathEntry is a single node in a dependency path leading to a vulnerable package.
type TransitiveContextualPathEntry struct {
	PackageVersionKey
	Function string `json:"function"`
}

// TransitiveContextualResponse is the per-CVE contextual analysis result: the vulnerable
// package, the function(s) involved, and the dependency path(s) reaching it.
type TransitiveContextualResponse struct {
	PackageVersionKey
	Functions []string                          `json:"functions"`
	Paths     [][]TransitiveContextualPathEntry `json:"paths"`
}

type TransitiveContextualService struct {
	client          *jfroghttpclient.JfrogHttpClient
	CatalogDetails  auth.ServiceDetails
	ScopeProjectKey string
}

func NewTransitiveContextualService(client *jfroghttpclient.JfrogHttpClient) *TransitiveContextualService {
	return &TransitiveContextualService{client: client}
}

func (tc *TransitiveContextualService) getUrl() string {
	return utils.AppendScopedProjectKeyParam(tc.CatalogDetails.GetUrl()+transitiveContextualApi, tc.ScopeProjectKey)
}

// GetContextualPaths requests, per CVE, the transitive dependency path(s) that make it reachable
// given the provided set of package keys.
func (tc *TransitiveContextualService) GetContextualPaths(cves []string, packages []PackageVersionKey) (map[string]TransitiveContextualResponse, error) {
	httpDetails := tc.CatalogDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()

	reqBody, err := json.Marshal(transitiveContextualRequest{Cves: cves, Packages: packages})
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to marshal transitive contextual request: %s", err.Error())
	}

	resp, body, err := tc.client.SendPost(tc.getUrl(), reqBody, &httpDetails)
	if err != nil {
		return nil, fmt.Errorf("failed while attempting to get transitive contextual analysis: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("got unexpected Catalog server response while attempting to get transitive contextual analysis: %w", err)
	}

	var result map[string]TransitiveContextualResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode transitive contextual response: %w", err)
	}
	return result, nil
}

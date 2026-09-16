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

const indirectContextualApi = "api/v1/dependencies/contextual"

// PackageVersionKey identifies a package by type, name, namespace, version and ecosystem.
type PackageVersionKey struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"`
}

type indirectContextualRequest struct {
	Cves     []string            `json:"cves"`
	Packages []PackageVersionKey `json:"packages"`
}

// IndirectContextualPathEntry is a single node in a dependency path leading to a vulnerable package.
type IndirectContextualPathEntry struct {
	PackageVersionKey
	Function string `json:"function"`
}

// IndirectContextualResponse is the per-CVE contextual analysis result: the vulnerable
// package, the function(s) involved, and the dependency path(s) reaching it.
type IndirectContextualResponse struct {
	PackageVersionKey
	Functions []string                        `json:"functions"`
	Paths     [][]IndirectContextualPathEntry `json:"paths"`
}

type IndirectContextualService struct {
	client          *jfroghttpclient.JfrogHttpClient
	CatalogDetails  auth.ServiceDetails
	ScopeProjectKey string
}

func NewIndirectContextualService(client *jfroghttpclient.JfrogHttpClient) *IndirectContextualService {
	return &IndirectContextualService{client: client}
}

func (tc *IndirectContextualService) getUrl() string {
	return utils.AppendScopedProjectKeyParam(tc.CatalogDetails.GetUrl()+indirectContextualApi, tc.ScopeProjectKey)
}

// GetContextualPaths requests, per CVE, the indirect dependency path(s) that make it reachable
// given the provided set of package keys.
func (tc *IndirectContextualService) GetContextualPaths(cves []string, packages []PackageVersionKey) (map[string]IndirectContextualResponse, error) {
	httpDetails := tc.CatalogDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()

	reqBody, err := json.Marshal(indirectContextualRequest{Cves: cves, Packages: packages})
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to marshal indirect contextual request: %s", err.Error())
	}

	resp, body, err := tc.client.SendPost(tc.getUrl(), reqBody, &httpDetails)
	if err != nil {
		return nil, fmt.Errorf("failed while attempting to get indirect contextual analysis: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("got unexpected Catalog server response while attempting to get indirect contextual analysis: %w", err)
	}

	var result map[string]IndirectContextualResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode indirect contextual response: %w", err)
	}
	return result, nil
}

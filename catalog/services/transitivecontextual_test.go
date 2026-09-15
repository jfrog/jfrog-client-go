package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCatalogDetails struct {
	auth.CommonConfigFields
}

func (d *testCatalogDetails) GetVersion() (string, error) {
	return "", nil
}

func newTestTransitiveContextualService(t *testing.T, serverUrl string) *TransitiveContextualService {
	details := &testCatalogDetails{}
	details.SetUrl(serverUrl + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := NewTransitiveContextualService(client)
	service.CatalogDetails = details
	return service
}

func TestGetContextualPaths_SendsExpectedRequestAndParsesResponse(t *testing.T) {
	var gotBody transitiveContextualRequest
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/dependencies/contextual", r.URL.Path)
		gotQuery = r.URL.RawQuery
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]TransitiveContextualResponse{
			"CVE-2024-1234": {
				PackageVersionKey: PackageVersionKey{
					Type:      "npm",
					Name:      "lodash",
					Namespace: "",
					Version:   "4.17.20",
					Ecosystem: GenericEcosystem,
				},
				Functions: []string{"merge"},
				Paths: [][]TransitiveContextualPathEntry{
					{
						{
							PackageVersionKey: PackageVersionKey{Type: "npm", Name: "app", Version: "1.0.0", Ecosystem: GenericEcosystem},
							Function:          "main",
						},
						{
							PackageVersionKey: PackageVersionKey{Type: "npm", Name: "lodash", Version: "4.17.20", Ecosystem: GenericEcosystem},
							Function:          "merge",
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	service := newTestTransitiveContextualService(t, server.URL)
	service.ScopeProjectKey = "myproj"
	packages := []PackageVersionKey{
		{Type: "npm", Name: "app", Version: "1.0.0", Ecosystem: GenericEcosystem},
		{Type: "npm", Name: "lodash", Version: "4.17.20", Ecosystem: GenericEcosystem},
	}
	result, err := service.GetContextualPaths([]string{"CVE-2024-1234"}, packages)

	require.NoError(t, err)
	assert.Equal(t, []string{"CVE-2024-1234"}, gotBody.Cves)
	assert.Equal(t, packages, gotBody.Packages)
	assert.Equal(t, "projectKey=myproj", gotQuery)
	require.Contains(t, result, "CVE-2024-1234")
	entry := result["CVE-2024-1234"]
	assert.Equal(t, "lodash", entry.Name)
	assert.Equal(t, []string{"merge"}, entry.Functions)
	require.Len(t, entry.Paths, 1)
	require.Len(t, entry.Paths[0], 2)
	assert.Equal(t, "app", entry.Paths[0][0].Name)
	assert.Equal(t, "main", entry.Paths[0][0].Function)
}

func TestGetContextualPaths_ServerError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := newTestTransitiveContextualService(t, server.URL)
	_, err := service.GetContextualPaths([]string{"CVE-2024-1234"}, []PackageVersionKey{{Type: "npm", Name: "lodash", Version: "4.17.20"}})

	assert.Error(t, err)
}

func TestGetContextualPaths_NoScopeProjectKey_OmitsQueryParam(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]TransitiveContextualResponse{})
	}))
	defer server.Close()

	service := newTestTransitiveContextualService(t, server.URL)
	_, err := service.GetContextualPaths([]string{"CVE-2024-1234"}, []PackageVersionKey{{Type: "npm", Name: "lodash", Version: "4.17.20"}})

	require.NoError(t, err)
	assert.Empty(t, gotQuery)
}

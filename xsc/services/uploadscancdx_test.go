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

type testXrayDetails struct {
	auth.CommonConfigFields
}

func (d *testXrayDetails) GetVersion() (string, error) {
	return "", nil
}

func newTestUploadScanCdxService(t *testing.T, serverUrl string) *UploadScanCdxService {
	details := &testXrayDetails{}
	details.SetUrl(serverUrl + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := NewUploadScanCdxService(client)
	service.XrayDetails = details
	return service
}

func TestUploadScanCdx_SendsExpectedRequestAndParsesResponse(t *testing.T) {
	var gotBody UploadScanCdxParams
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/xsc/scan-cdx/upload", r.URL.Path)
		gotQuery = r.URL.RawQuery
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(UploadScanCdxResponse{Repository: "myproj-frogbot", Path: "github.com/org/repo/main/commits/source_code_1.cdx.json"})
	}))
	defer server.Close()

	service := newTestUploadScanCdxService(t, server.URL)
	service.ScopeProjectKey = "myproj"
	resp, err := service.Upload(UploadScanCdxParams{
		RepoName: "frogbot",
		RepoPath: "github.com/org/repo/main/commits",
		FileName: "source_code_1.cdx.json",
		Bom:      `{"bomFormat":"CycloneDX"}`,
	})

	require.NoError(t, err)
	assert.Equal(t, "frogbot", gotBody.RepoName)
	assert.Equal(t, "projectKey=myproj", gotQuery)
	assert.Equal(t, "myproj-frogbot", resp.Repository)
	assert.Equal(t, "github.com/org/repo/main/commits/source_code_1.cdx.json", resp.Path)
}

func TestUploadScanCdx_ServerError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := newTestUploadScanCdxService(t, server.URL)
	_, err := service.Upload(UploadScanCdxParams{RepoName: "frogbot", RepoPath: "p", FileName: "f.cdx.json", Bom: `{}`})

	assert.Error(t, err)
}

func TestUploadScanCdx_NoScopeProjectKey_OmitsQueryParam(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(UploadScanCdxResponse{Repository: "frogbot", Path: "p/f.cdx.json"})
	}))
	defer server.Close()

	service := newTestUploadScanCdxService(t, server.URL)
	_, err := service.Upload(UploadScanCdxParams{RepoName: "frogbot", RepoPath: "p", FileName: "f.cdx.json", Bom: `{}`})

	require.NoError(t, err)
	assert.Empty(t, gotQuery)
}

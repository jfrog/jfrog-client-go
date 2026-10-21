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

func newTestApplicabilityService(t *testing.T, serverUrl string) *ApplicabilityService {
	details := &testXrayDetails{}
	details.SetUrl(serverUrl + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := NewApplicabilityService(client)
	service.XrayDetails = details
	return service
}

func TestGetApplicabilityInput_SendsExpectedRequestAndParsesResponse(t *testing.T) {
	var gotBody applicabilityInputRequest
	var gotProjectKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/internal/cve_applicability_input_vulnerabilities", r.URL.Path)
		gotProjectKey = r.URL.Query().Get("projectKey")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"explanation":{"CVE-2024-1234":"explanation text"}}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	service := newTestApplicabilityService(t, server.URL)
	service.ScopeProjectKey = "myproj"
	applicabilityInput, err := service.GetApplicabilityInput([]string{"CVE-2024-1234"})

	require.NoError(t, err)
	assert.Equal(t, []string{"CVE-2024-1234"}, gotBody.Vulnerabilities)
	assert.Equal(t, "myproj", gotProjectKey)
	assert.Equal(t, map[string]string{"CVE-2024-1234": "explanation text"}, applicabilityInput.Explanation)
}

func TestGetApplicabilityInput_ReturnsErrorOnUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`"Cve Applicability is not available"`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	applicabilityInput, err := newTestApplicabilityService(t, server.URL).GetApplicabilityInput([]string{"CVE-2024-1234"})

	assert.ErrorContains(t, err, "Cve Applicability is not available")
	assert.Nil(t, applicabilityInput)
}

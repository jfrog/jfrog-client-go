package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	accessAuth "github.com/jfrog/jfrog-client-go/access/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSystemService(t *testing.T, serverUrl string) *SystemService {
	t.Helper()
	client, err := jfroghttpclient.JfrogClientBuilder().SetRetries(0).Build()
	require.NoError(t, err)
	details := accessAuth.NewAccessDetails()
	details.SetUrl(serverUrl + "/access/")
	details.SetAccessToken("tok")
	service := NewSystemService(client)
	service.ServiceDetails = details
	return service
}

func TestGetRootCertificate_NoAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/access/api/v1/cert/root", r.URL.Path)
		assert.Empty(t, r.Header.Get("Authorization"))
		_, _ = w.Write([]byte("QUJD\n"))
	}))
	defer server.Close()

	cert, err := newTestSystemService(t, server.URL).GetRootCertificate()
	require.NoError(t, err)
	assert.Equal(t, "QUJD", cert)
}

func TestGetRootCertificate_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := newTestSystemService(t, server.URL).GetRootCertificate()
	assert.ErrorContains(t, err, "404")
	assert.ErrorContains(t, err, "api/v1/cert/root")
}

func TestCheckReadiness_SendsBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/access/api/v1/system/readiness", r.URL.Path)
		assert.Equal(t, "Bearer tok", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	status, err := newTestSystemService(t, server.URL).CheckReadiness()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}

func TestCheckReadiness_ReturnsRejectStatus(t *testing.T) {
	for _, code := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		hits := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			hits++
			w.WriteHeader(code)
		}))

		status, err := newTestSystemService(t, server.URL).CheckReadiness()
		server.Close()
		require.NoError(t, err)
		assert.Equal(t, code, status)
		assert.Equal(t, 1, hits)
	}
}

func TestCheckReadiness_TransportError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	service := newTestSystemService(t, server.URL)
	server.Close()

	status, err := service.CheckReadiness()
	assert.Error(t, err)
	assert.Equal(t, 0, status)
}

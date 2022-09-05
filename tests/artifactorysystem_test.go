package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("getVersion", testGetVersion)
	t.Run("getServiceId", testGetServiceId)
	t.Run("getRunningNodes", testGetRunningNodes)
	t.Run("getConfigDescriptor", testGetConfigDescriptor)
	t.Run("activateKeyEncryption", testActivateKeyEncryption)
	t.Run("deactivateKeyEncryption", testDeactivateKeyEncryption)
	t.Run("deactivateKeyEncryptionNotEncrypted", testDeactivateKeyEncryptionNotEncrypted)
}

func testGetVersion(t *testing.T) {
	version, err := testsSystemService.GetVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, version)
}

func testGetServiceId(t *testing.T) {
	serviceId, err := testsSystemService.GetServiceId()
	assert.NoError(t, err)
	assert.NotEmpty(t, serviceId)
}

func testGetRunningNodes(t *testing.T) {
	runningNodes, err := testsSystemService.GetRunningNodes()
	assert.NoError(t, err)
	assert.NotEmpty(t, runningNodes)
}

func testGetConfigDescriptor(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodGet, r.Method)

		// Check URL
		assert.Equal(t, "/api/system/configuration", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<config></config>"))
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createMockSystemService(t, ts.URL)
	results, err := service.GetConfigDescriptor()
	assert.NoError(t, err)
	assert.Equal(t, "<config></config>", results)
}

func testActivateKeyEncryption(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/system/encrypt", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Done"))
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createMockSystemService(t, ts.URL)
	assert.NoError(t, service.ActivateKeyEncryption())
}

func testDeactivateKeyEncryption(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/system/decrypt", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Done"))
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createMockSystemService(t, ts.URL)
	wasEncrypted, err := service.DeactivateKeyEncryption()
	assert.NoError(t, err)
	assert.True(t, wasEncrypted)
}

func testDeactivateKeyEncryptionNotEncrypted(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/system/decrypt", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Cannot decrypt without artifactory key file"))
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createMockSystemService(t, ts.URL)
	wasEncrypted, err := service.DeactivateKeyEncryption()
	assert.NoError(t, err)
	assert.False(t, wasEncrypted)
}

func createMockSystemService(t *testing.T, url string) *services.SystemService {
	// Create artifactory details
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(url + "/")

	// Create http client
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetClientCertPath(rtDetails.GetClientCertPath()).
		SetClientCertKeyPath(rtDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(rtDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err, "Failed to create Artifactory client: %v\n")

	// Create system service
	return services.NewSystemService(rtDetails, client)
}

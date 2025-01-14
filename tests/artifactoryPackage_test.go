package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var leadFileRequest = services.LeadFileParams{
	PackageVersion:  "1.0.0",
	PackageName:     "test-package",
	PackageRepoName: "test-repo",
	PackageType:     "test-type",
}

func TestPackage(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("TestGetLeadFileSuccessfully", TestGetLeadFileSuccessfully)
}

func TestGetLeadFileSuccessfully(t *testing.T) {
	handlerFunc := createDefaultHandlerFunc(t)
	mockServer, packageService := createMockPackageServer(t, handlerFunc)

	expectedLeadFile := "path/to/lead/file"
	defer mockServer.Close()

	leadFilePath, err := packageService.GetPackageLeadFile(leadFileRequest)
	assert.NoError(t, err)
	assert.Equal(t, expectedLeadFile, string(leadFilePath))
}

func createDefaultHandlerFunc(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/packagesSearch/leadFile" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			writeMockLeadFileResponse(t, w, []byte("path/to/lead/file"))
		} else {
			t.Errorf("Unexpected URL: got %s, want %s", r.RequestURI, "/api/packagesSearch/leadFile")
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}

func writeMockLeadFileResponse(t *testing.T, w http.ResponseWriter, payload []byte) {
	_, err := w.Write(payload)
	assert.NoError(t, err)
}

func createMockPackageServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *services.PackageService) {
	testServer := httptest.NewServer(testHandler)

	serviceDetails := auth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/")

	packageService, err := jfroghttpclient.JfrogClientBuilder().Build()

	assert.NoError(t, err)
	return testServer, &services.PackageService{Client: packageService, ArtDetails: serviceDetails}
}

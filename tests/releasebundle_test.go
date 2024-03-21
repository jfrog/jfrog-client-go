package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory"
	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImportReleaseBundle(t *testing.T) {
	mockServer, rbService := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/"+services.ReleaseBundleImportRestApiEndpoint {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`
{
  "errors" : [ {
    "status" : 400,
    "message" : "Bundle already exists"
			  } ]
}
`))
			assert.NoError(t, err)
		}
	})
	defer mockServer.Close()
	err := rbService.ImportReleaseBundle("releasebundle_test.go")
	assert.NoError(t, err)
}

func createMockServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, artifactory.ArtifactoryServicesManager) {
	testServer := httptest.NewServer(testHandler)

	rtDetails := artifactoryAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(testServer.URL + "/")

	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetDryRun(false).
		Build()

	if err != nil {
		t.Error(err)
	}

	artService, err := artifactory.New(serviceConfig)
	assert.NoError(t, err)
	return testServer, artService
}

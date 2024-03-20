package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

const (
	exportPath = "/a/b/c"
)

func TestExportEmptyParams(t *testing.T) {
	initArtifactoryTest(t)
	testExport(t, services.NewExportParams(exportPath))
}

func TestExportTrueValues(t *testing.T) {
	initArtifactoryTest(t)
	testExport(t, createExportTestParams(true))
}

func TestExportDryRun(t *testing.T) {
	initArtifactoryTest(t)
	service := createExportService(t, "127.0.0.1")
	service.DryRun = true
	err := service.Export(createExportTestParams(true))
	assert.NoError(t, err)
}

func TestExportFalseValues(t *testing.T) {
	initArtifactoryTest(t)
	testExport(t, createExportTestParams(false))
}

func testExport(t *testing.T, exportParams services.ExportParams) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/export/system", r.URL.Path)

		// Check body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		response := services.ExportBody{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, exportParams.ExportPath, response.ExportPath)
		assert.Equal(t, exportParams.CreateArchive, response.CreateArchive)
		assert.Equal(t, exportParams.ExcludeContent, response.ExcludeContent)
		assert.Equal(t, exportParams.IncludeMetadata, response.IncludeMetadata)
		assert.Equal(t, exportParams.M2, response.M2)
		assert.Equal(t, exportParams.Verbose, response.Verbose)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Promote docker image
	service := createExportService(t, ts.URL)
	err := service.Export(exportParams)
	assert.NoError(t, err)
}

func createExportService(t *testing.T, url string) *services.ExportService {
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

	// Create export service
	exportService := services.NewExportService(rtDetails, client)
	return exportService
}

func createExportTestParams(valueToTest bool) services.ExportParams {
	params := services.NewExportParams(exportPath)
	params.CreateArchive = &valueToTest
	params.ExcludeContent = &valueToTest
	params.IncludeMetadata = &valueToTest
	params.M2 = &valueToTest
	params.Verbose = &valueToTest
	return params
}

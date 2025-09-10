package tests

import (
	"strconv"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

func TestXrayDownloadIndexer(t *testing.T) {
	initXrayTest(t)
	// Create temp dir for downloading the indexer binary
	outputDir, err := fileutils.CreateTempDir()
	assert.NoError(t, err)
	// Create mock Xray server
	xrayServerPort := xray.StartXrayMockServer(t)
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err)
	// Create indexer service
	indexerService := services.NewIndexerService(client)
	indexerService.XrayDetails = xrayDetails
	indexerService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	// Download the indexer binary
	downloadedFilePath, err := indexerService.Download(outputDir, "test-indexer")
	assert.NoError(t, err)
	assert.Equal(t, outputDir+"/test-indexer", downloadedFilePath)
	// Verify the indexer binary was downloaded successfully
	exists, err := fileutils.IsFileExists(downloadedFilePath, false)
	assert.NoError(t, err)
	assert.True(t, exists)
	// Clean up
	err = fileutils.RemoveTempDir(outputDir)
	assert.NoError(t, err)
}

package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	xrayServices "github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func initXrayEnrichTest(t *testing.T) (xrayServerPort int, xrayDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
	var err error
	initXrayTest(t)
	xrayServerPort = xray.StartXrayMockServer(t)
	xrayDetails = GetXrayDetails()
	client, err = jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err)
	return
}

func TestIsImportSucceeded(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayEnrichTest(t)
	testsEnrichService := xrayServices.NewEnrichService(client)
	testsEnrichService.XrayDetails = xrayDetails
	testsEnrichService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	params := xrayServices.XrayGraphImportParams{SBOMInput: []byte("")}
	result, err := testsEnrichService.ImportGraph(params, "test")
	assert.NoError(t, err)
	assert.Equal(t, result, xray.TestMultiScanId)
}

func TestGetImportResults(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayEnrichTest(t)
	testsEnrichService := xrayServices.NewEnrichService(client)
	testsEnrichService.XrayDetails = xrayDetails
	testsEnrichService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	result, err := testsEnrichService.GetImportGraphResults(xray.TestMultiScanId)
	assert.NoError(t, err)
	assert.Equal(t, result.ScanId, xray.TestMultiScanId)
	assert.Len(t, result.Vulnerabilities, 1)

}

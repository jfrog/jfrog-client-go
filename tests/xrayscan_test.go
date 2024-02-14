package tests

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	xrayServices "github.com/jfrog/jfrog-client-go/xray/services"
)

var testsXrayScanService *services.XrayScanService
var testsScanService *xrayServices.ScanService

func TestNewXrayScanService(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer(t)
	artDetails := GetRtDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(artDetails.GetClientCertPath()).
		SetClientCertKeyPath(artDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(artDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		t.Error(err)
	}
	testsXrayScanService = services.NewXrayScanService(client)
	testsXrayScanService.ArtDetails = artDetails
	testsXrayScanService.ArtDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/")

	// Run tests
	tests := []struct {
		name        string
		buildName   string
		buildNumber string
		expected    string
	}{
		{name: "scanCleanBuild", buildName: xray.CleanScanBuildName, buildNumber: "3", expected: xray.CleanXrayScanResponse},
		{name: "scanVulnerableBuild", buildName: xray.VulnerableBuildName, buildNumber: "3", expected: xray.VulnerableXrayScanResponse},
		{name: "scanFatalBuild", buildName: xray.FatalScanBuildName, buildNumber: "3", expected: xray.FatalErrorXrayScanResponse},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanBuild(t, test.buildName, test.buildNumber, test.expected)
		})
	}
}

func scanBuild(t *testing.T, buildName, buildNumber, expected string) {
	params := services.NewXrayScanParams()
	params.BuildName = buildName
	params.BuildNumber = buildNumber
	result, err := testsXrayScanService.ScanBuild(params)
	if err != nil {
		t.Error(err)
	}

	expected = strings.ReplaceAll(expected, "\n", "")
	if string(result) != expected {
		t.Error("Expected:", expected, "Got: ", string(result))
	}
}

func TestIsXscEnabled(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayScanTest(t)
	testsScanService = xrayServices.NewScanService(client)
	testsScanService.XrayDetails = xrayDetails
	testsScanService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	result, err := testsScanService.IsXscEnabled()
	assert.NoError(t, err)
	assert.Equal(t, xray.TestXscVersion, result)
}

func TestSendScanGitInfoContext(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayScanTest(t)
	testsScanService = xrayServices.NewScanService(client)
	testsScanService.XrayDetails = xrayDetails
	testsScanService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	// Run tests
	tests := []struct {
		name           string
		gitInfoContext *xrayServices.XscGitInfoContext
		expected       string
	}{
		{name: "ValidGitInfoContext", gitInfoContext: &xray.GitInfoContextWithMinimalRequiredFields, expected: xray.TestMultiScanId},
		{name: "InvalidGitInfoContext", gitInfoContext: &xray.GitInfoContextWithMissingFields, expected: xray.XscGitInfoBadResponse},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sendGitInfoContext(t, test.gitInfoContext, test.expected)
		})
	}
}

func sendGitInfoContext(t *testing.T, gitInfoContext *xrayServices.XscGitInfoContext, expected string) {
	result, err := testsScanService.SendScanGitInfoContext(gitInfoContext)
	if err != nil {
		assert.ErrorContains(t, err, expected)
		return
	}
	assert.Equal(t, expected, result)
}

func initXrayScanTest(t *testing.T) (xrayServerPort int, xrayDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
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

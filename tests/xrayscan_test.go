package tests

import (
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
var gitInfoContextWithMinimalRequiredFields = xrayServices.XscGitInfoContext{
	GitRepoUrl: "https://git.jfrog.info/projects/XSC/repos/xsc-service",
	BranchName: "feature/XRAY-123-cool-feature",
	CommitHash: "acc5e24e69a-d3c1-4022-62eb-69e4a1e5",
}
var gitInfoContextWithMissingFields = xrayServices.XscGitInfoContext{
	GitRepoUrl: "https://git.jfrog.info/projects/XSC/repos/xsc-service",
	BranchName: "feature/XRAY-123-cool-feature",
}
var testMultiScanId = "3472b4e2-bddc-11ee-a9c9-acde48001122"
var testXscVersion = "1.0.0"

func TestNewXrayScanService(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer()
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
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer()
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		t.Error(err)
	}
	testsScanService = xrayServices.NewScanService(client)
	testsScanService.XrayDetails = xrayDetails
	testsScanService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	result, err := testsScanService.IsXscEnabled()
	if err != nil {
		t.Error(err)
	}

	if result != testXscVersion {
		t.Error("Expected:", testXscVersion, "Got: ", result)
	}
}

func TestSendScanGitInfoContext(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer()
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		t.Error(err)
	}
	testsScanService = xrayServices.NewScanService(client)
	testsScanService.XrayDetails = xrayDetails
	testsScanService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	// Run tests
	tests := []struct {
		name           string
		gitInfoContext *xrayServices.XscGitInfoContext
		expected       string
	}{
		{name: "ValidGitInfoContext", gitInfoContext: &gitInfoContextWithMinimalRequiredFields, expected: testMultiScanId},
		{name: "InvalidGitInfoContext", gitInfoContext: &gitInfoContextWithMissingFields, expected: xray.XscGitInfoBadResponse},
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

	if result != expected {
		t.Error("Expected:", expected, "Got: ", result)
	}
}

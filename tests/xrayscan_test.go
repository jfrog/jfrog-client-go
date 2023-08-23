package tests

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/xray/manager"
	"github.com/jfrog/jfrog-client-go/xray/scan"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

var testsXrayScanService *services.XrayScanService

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

func TestXrayScanGraph(t *testing.T) {
	initXrayTest(t)
	mockScanId := "9c9dbd61-f544-4e33-4613-34727043d71f"
	xrayServerPort := xray.StartXrayMockServer()
	xrayDetails := newTestXrayDetails(GetXrayDetails())
	xrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	cfp := auth.ServiceDetails(xrayDetails)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cfp).
		Build()
	assert.NoError(t, err)
	securityServiceManager, err = manager.New(serviceConfig)
	assert.NoError(t, err)
	assertSecurityManagerType(t)

	scanId, err := securityServiceManager.ScanGraph(&scan.XrayGraphScanParams{})
	assert.NoError(t, err)
	assert.Equal(t, mockScanId, scanId)
	_, err = securityServiceManager.GetScanGraphResults(scanId, false, false)
	assert.NoError(t, err)
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

package tests

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/xray/manager"
	"github.com/jfrog/jfrog-client-go/xray/scan"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
)

func TestXscVersion(t *testing.T) {
	initXscTest(t)
	version := GetXrayDetails().GetXscVersion()
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func TestXscScanGraph(t *testing.T) {
	initXscTest(t)
	mockScanId := "9c9dbd61-f544-4e33-4613-34727043d71f"
	mockMultiScanId := "f2a8d4fe-40e6-11ee-84e4-02ee10c7f40e"
	xrayServerPort := xray.StartXrayMockServer()
	xrayDetails := newTestXrayDetails(GetXrayDetails())
	xrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	xrayDetails.SetXscUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xsc/")

	cfp := auth.ServiceDetails(xrayDetails)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cfp).
		Build()
	assert.NoError(t, err)
	securityServiceManager, err = manager.New(serviceConfig)
	assert.NoError(t, err)
	// Assert correct security manager
	assertSecurityManagerType(t)

	graphParams := &scan.XrayGraphScanParams{}
	graphParams.XscGitInfoContext = &scan.XscGitInfoContext{}
	scanId, err := securityServiceManager.ScanGraph(graphParams)
	assert.NoError(t, err)
	assert.Equal(t, mockMultiScanId, graphParams.MultiScanId)
	assert.Equal(t, mockScanId, scanId)
	_, err = securityServiceManager.GetScanGraphResults(scanId, false, false)
	assert.NoError(t, err)
}

func assertSecurityManagerType(t *testing.T) {
	switch securityServiceManager.(type) {
	case *manager.XscServicesManger:
		assert.Equal(t, true, *TestXsc)
	case *manager.XrayServicesManager:
		assert.Equal(t, false, *TestXsc)
	}
}

func initXscTest(t *testing.T) {
	if !*TestXsc {
		t.Skip("Skipping xray test. To run xray test add the '-test.xsc=true' option.")
	}
}

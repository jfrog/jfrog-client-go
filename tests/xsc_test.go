package tests

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
)

func TestXscScanGraph(t *testing.T) {
	initXscTest(t)
	expectedScanId := "9c9dbd61-f544-4e33-4613-34727043d71f"
	mockMultiScanId := "f2a8d4fe-40e6-11ee-84e4-02ee10c7f40e"

	tests := []struct {
		name                string
		xrayGraphParams     *services.XrayGraphScanParams
		expectedMultiScanId string
	}{
		{
			name:                "XscScanWithContext",
			xrayGraphParams:     &services.XrayGraphScanParams{XscGitInfoContext: &services.XscGitInfoContext{}},
			expectedMultiScanId: mockMultiScanId,
		}, {
			name:                "XscScanNoContext",
			xrayGraphParams:     &services.XrayGraphScanParams{},
			expectedMultiScanId: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			scanId, err := securityServiceManager.ScanGraph(test.xrayGraphParams)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedMultiScanId, test.xrayGraphParams.MultiScanId)
			assert.Equal(t, expectedScanId, scanId)

			_, err = securityServiceManager.GetScanGraphResults(scanId, false, false)
			assert.NoError(t, err)
		})
	}
}

func TestXscEnabled(t *testing.T) {
	initXscTest(t)
	version, err := securityServiceManager.IsXscEnabled()
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0", version)
}

func initXscTest(t *testing.T) {
	initializeTestSecurityManager(t, initMockXscServer())
}

func initMockXscServer() testXrayDetails {
	xrayServerPort := xray.StartXrayMockServer()
	xrayDetails := newTestXrayDetails(GetXrayDetails())
	// Reroutes URLs to mock server
	xrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	xrayDetails.SetXscUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xsc/")
	return xrayDetails
}

func initializeTestSecurityManager(t *testing.T, xscDetails testXrayDetails) {
	cfp := auth.ServiceDetails(xscDetails)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cfp).
		Build()
	assert.NoError(t, err)
	securityServiceManager, err = services.New(serviceConfig)
	assert.NoError(t, err)
	// Assert correct security manager Xsc/Xray
	assert.IsType(t, securityServiceManager,&services.XscServicesManger{})
}
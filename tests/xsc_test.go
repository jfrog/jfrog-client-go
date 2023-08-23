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

func TestXscScanGraph(t *testing.T) {
	initXscTest(t)
	expectedScanId := "9c9dbd61-f544-4e33-4613-34727043d71f"
	mockMultiScanId := "f2a8d4fe-40e6-11ee-84e4-02ee10c7f40e"

	tests := []struct {
		name                string
		xrayGraphParams     *scan.XrayGraphScanParams
		expectedMultiScanId string
	}{
		{
			name:                "XscScanWithContext",
			xrayGraphParams:     &scan.XrayGraphScanParams{XscGitInfoContext: &scan.XscGitInfoContext{}},
			expectedMultiScanId: mockMultiScanId,
		}, {
			name:                "XscScanNoContext",
			xrayGraphParams:     &scan.XrayGraphScanParams{},
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
	enabled, version, err := securityServiceManager.IsXscEnabled()
	assert.NoError(t, err)
	assert.Equal(t, true, enabled)
	assert.Equal(t, "0.0.0", version)
}

func initXscTest(t *testing.T) {
	if !*TestXsc {
		t.Skip("Skipping xray test. To run xray test add the '-test.xsc=true' option.")
	}
	prepareXscTest(t)
}

func prepareXscTest(t *testing.T) {
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
	securityServiceManager, err = manager.New(serviceConfig)
	assert.NoError(t, err)
	// Assert correct security manager Xsc/Xray
	assertSecurityManagerType(t)
}

func assertSecurityManagerType(t *testing.T) {
	switch securityServiceManager.(type) {
	case *manager.XscServicesManger:
		assert.Equal(t, true, *TestXsc)
	case *manager.XrayServicesManager:
		assert.Equal(t, false, *TestXsc)
	}
}

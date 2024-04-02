package tests

import (
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"testing"
)

func TestXscVersion(t *testing.T) {
	initXscTest(t, "")
	version, err := GetXscDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func initXscTest(t *testing.T, minVersion string) {
	if !*TestXsc {
		t.Skip("Skipping xsc test. To run xsc test add the '-test.xsc=true' option.")
	}
	validateXscVersion(t, minVersion)
}
func validateXscVersion(t *testing.T, minVersion string) {
	// Validate active XSC server.
	version, err := GetXscDetails().GetVersion()
	if err != nil {
		t.Skip(err)
	}
	// Validate minimum XSC version.
	err = clientUtils.ValidateMinimumVersion(clientUtils.Xsc, version, minVersion)
	if err != nil {
		t.Skip(err)
	}
}

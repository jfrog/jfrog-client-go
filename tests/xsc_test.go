package tests

import (
	"testing"
)

func TestXscVersion(t *testing.T) {
	initXscTest(t)
	version, err := GetXscDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func initXscTest(t *testing.T) {
	if !*TestXsc {
		t.Skip("Skipping xsc test. To run xsc test add the '-test.xsc=true' option.")
	}
}

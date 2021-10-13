package tests

import (
	"testing"
)

func TestXrayVersion(t *testing.T) {
	initXrayTest(t)
	version, err := GetXrayDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func initXrayTest(t *testing.T) {
	if !*TestXray {
		t.Skip("Skipping xray test. To run xray test add the '-test.xray=true' option.")
	}
}

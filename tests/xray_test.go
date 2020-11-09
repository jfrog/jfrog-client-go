package tests

import (
	"testing"
)

func TestXrayVersion(t *testing.T) {
	if *XrayUrl == "" {
		t.Skip("Xray is not being tested, skipping...")
	}

	version, err := GetXrayDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

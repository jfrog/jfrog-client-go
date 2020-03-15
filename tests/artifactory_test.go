package tests

import (
	"testing"
)

func TestGetArtifactoryVersion(t *testing.T) {
	version, err := GetRtDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

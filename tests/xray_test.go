package tests

import (
	"testing"
)

func TestGetXrayVersion(t *testing.T) {
	version, err := GetXrayDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

package tests

import (
	"testing"
)

func TestGetArtifactoryVersion(t *testing.T) {
	initArtifactoryTest(t)
	version, err := GetRtDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func initArtifactoryTest(t *testing.T) {
	if !*TestArtifactory {
		t.Skip("Skipping artifactory test. To run artifactory test add the '-test.artifactory=true' option.")
	}
}

func initRepositoryTest(t *testing.T) {
	if !*TestRepository {
		t.Skip("Skipping repository test. To run repository test add the '-test.repository=true' option.")
	}
}

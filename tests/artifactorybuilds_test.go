package tests

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuildRuns(t *testing.T) {
	initArtifactoryTest(t)

	// Create a build
	buildName := fmt.Sprintf("%s-%s", "build-run", getRunId())
	err := createDummyBuild(buildName)
	assert.NoError(t, err)

	runs, found, err := testBuildInfoService.GetBuildRuns(services.BuildInfoParams{BuildName: buildName})
	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotEmpty(t, runs.Uri)
	assert.NotEmpty(t, runs.BuildsNumbers)
	assert.Equal(t, "/"+buildNumber, runs.BuildsNumbers[0].Uri)
	assert.NotEmpty(t, runs.BuildsNumbers[0].Started)
}

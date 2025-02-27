package tests

import (
	"fmt"
	buildinfo "github.com/jfrog/build-info-go/entities"
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

func TestDeleteBuildRuns(t *testing.T) {
	initArtifactoryTest(t)

	// Create a build
	buildName := fmt.Sprintf("%s-%s", "build-run", getRunId())
	err := createDummyBuild(buildName)
	assert.NoError(t, err)

	// Check the number of builds
	runs, found, err := testBuildInfoService.GetBuildRuns(services.BuildInfoParams{BuildName: buildName})
	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotEmpty(t, runs.BuildsNumbers)

	// Delete the build
	buildInfo := &buildinfo.BuildInfo{Name: buildName, Number: buildNumber}
	err = testBuildInfoService.DeleteBuildInfo(buildInfo, "", 1)
	assert.NoError(t, err)

	// Verify the number of builds is 0
	runs, found, err = testBuildInfoService.GetBuildRuns(services.BuildInfoParams{BuildName: buildName})
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, runs)
}

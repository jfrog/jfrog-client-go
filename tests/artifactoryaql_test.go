//go:build itest

package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateAqlQueryForBuildInfoJson(t *testing.T) {
	initArtifactoryTest(t)

	// Create a build
	buildName := fmt.Sprintf("%s-%s", "a / \\ | \t * ? : ; \\ / %b", getRunId())
	err := createDummyBuild(buildName)
	assert.NoError(t, err)

	// Run AQL to get the build from Artifactory
	aqlQuery := utils.CreateAqlQueryForBuildInfoJson("", buildName, buildNumber, buildTimestamp)
	stream, err := testsAqlService.ExecAql(aqlQuery)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, stream.Close())
	}()

	// Parse AQL results
	aqlResults, err := io.ReadAll(stream)
	assert.NoError(t, err)
	parsedResult := new(utils.AqlSearchResult)
	err = json.Unmarshal(aqlResults, parsedResult)
	assert.NoError(t, err)
	assert.Len(t, parsedResult.Results, 1)

	// Verify build checksum exist
	assert.NotEmpty(t, parsedResult.Results[0].Actual_Sha1)
	assert.NotEmpty(t, parsedResult.Results[0].Actual_Md5)

	// Delete build
	encodedBuildName := strings.TrimSuffix(parsedResult.Results[0].Path, "-"+buildTimestamp+".json")
	assert.NoError(t, deleteBuild(encodedBuildName))
}

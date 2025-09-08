//go:build itest

package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXrayBinMgr(t *testing.T) {
	initXrayTest(t)
	t.Run("addBuildsToIndexing", addBuildsToIndexing)
}

func addBuildsToIndexing(t *testing.T) {
	buildName := fmt.Sprintf("%s-%s", "build1", getRunId())
	defer func() {
		assert.NoError(t, deleteBuildIndex(buildName))
		assert.NoError(t, deleteBuild(buildName))
	}()
	// Create a build
	err := createDummyBuild(buildName)
	assert.NoError(t, err)

	// Index build
	err = testXrayBinMgrService.AddBuildsToIndexing([]string{buildName})
	assert.NoError(t, err)

	// Assert build contained in the indexed build list
	indexedBuilds, err := getIndexedBuilds()
	assert.NoError(t, err)
	assert.Contains(t, indexedBuilds, buildName)
}

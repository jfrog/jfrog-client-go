//go:build itest

package tests

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXrayBinMgr(t *testing.T) {
	initXrayTest(t)
	t.Run("addBuildsToIndexing", addBuildsToIndexing)
}

func addBuildsToIndexing(t *testing.T) {
	buildName := fmt.Sprintf("%s-%s", "build1", getRunId())
	t.Cleanup(func() {
		if err := deleteBuildIndex(buildName); err != nil {
			t.Logf("Failed to delete build index: %v", err)
		}
		if err := deleteBuild(buildName); err != nil {
			t.Logf("Failed to delete build: %v", err)
		}
	})
	// Create a build
	err := createDummyBuild(buildName)
	require.NoError(t, err)

	// Index build
	err = testXrayBinMgrService.AddBuildsToIndexing([]string{buildName})
	require.NoError(t, err)

	// Assert build contained in the indexed build list
	assert.Eventuallyf(t, func() bool {
		indexedBuilds, err := getIndexedBuilds()
		if err != nil {
			t.Logf("Failed to get indexed builds: %v", err)
			return false
		}
		return slices.Contains(indexedBuilds, buildName)
	}, time.Second*30, time.Millisecond*500, "Build %s not found in indexed builds", buildName)
}

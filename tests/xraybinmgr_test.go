package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestXrayBinMgr(t *testing.T) {
	if *XrayUrl == "" {
		t.Skip("Xray is not being tested, skipping...")
	}

	t.Run("addBuildsToIndexing", addBuildsToIndexing)
}

func addBuildsToIndexing(t *testing.T) {
	buildName := fmt.Sprintf("%s-%d", "jfrog-build1", time.Now().Unix())
	defer deleteBuild(buildName)

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

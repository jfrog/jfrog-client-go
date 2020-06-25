package utils

import (
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

func TestLoadMissingProperties(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	notSortedWithProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_nosorted_withprops.json"), "results")
	sortedNoProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_sorted_noprops.json"), "results")
	utils.MAX_BUFFER_SIZE = 3
	cr, err := loadMissingProperties(sortedNoProps, notSortedWithProps)
	assert.NoError(t, err)
	isMatch, err := fileutils.FilesMath(cr.GetFilePath(), filepath.Join(testDataPath, "load_missing_props_expected_results.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
}

func TestFilterBuildAqlSearchResults(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	resultsToFilter := content.NewContentReader(filepath.Join(testDataPath, "filter_build_aql_search.json"), "results")
	buildArtifactsSha := map[string]byte{"a": 2, "b": 2, "c": 2}
	resultReader, err := filterBuildAqlSearchResults(resultsToFilter, buildArtifactsSha, "myBuild", "1")
	assert.NoError(t, err)
	isMatch, err := fileutils.FilesMath(resultReader.GetFilePath(), filepath.Join(testDataPath, "filter_build_aql_search_expected.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
}

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
	oldMaxSize := utils.MaxBufferSize
	defer func() { utils.MaxBufferSize = oldMaxSize }()
	for i := 0; i < 2; i++ {
		testDataPath, err := getBaseTestDir()
		assert.NoError(t, err)
		notSortedWithProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_nosorted_withprops.json"), content.DefaultKey)
		sortedNoProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_sorted_noprops.json"), content.DefaultKey)
		reader, err := loadMissingProperties(sortedNoProps, notSortedWithProps)
		defer reader.Close()
		assert.NoError(t, err)
		isMatch, err := fileutils.FilesIdentical(reader.GetFilePath(), filepath.Join(testDataPath, "load_missing_props_expected_results.json"))
		assert.NoError(t, err)
		assert.True(t, isMatch)
		utils.MaxBufferSize = 3
	}
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	notSortedWithProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_nosorted_by_created_withprops.json"), content.DefaultKey)
	sortedNoProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_sorted_by_created_noprops.json"), content.DefaultKey)
	reader, err := loadMissingProperties(sortedNoProps, notSortedWithProps)
	defer reader.Close()
	assert.NoError(t, err)
	isMatch, err := fileutils.FilesIdentical(reader.GetFilePath(), filepath.Join(testDataPath, "load_missing_props_by_created_expected_results.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	utils.MaxBufferSize = 3
}

func TestFilterBuildAqlSearchResults(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	resultsToFilter := content.NewContentReader(filepath.Join(testDataPath, "filter_build_aql_search.json"), content.DefaultKey)
	buildArtifactsSha := map[string]int{"a": 2, "b": 2, "c": 2}
	resultReader, err := filterBuildAqlSearchResults(resultsToFilter, buildArtifactsSha, "myBuild", "1")
	defer resultReader.Close()
	assert.NoError(t, err)
	isMatch, err := fileutils.FilesIdentical(resultReader.GetFilePath(), filepath.Join(testDataPath, "filter_build_aql_search_expected.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
}

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
		testDataPath := getBaseTestDir(t)
		notSortedWithProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_nosorted_withprops.json"), content.DefaultKey)
		sortedNoProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_sorted_noprops.json"), content.DefaultKey)
		reader, err := loadMissingProperties(sortedNoProps, notSortedWithProps)
		defer readerCloseAndAssert(t, reader)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(reader.GetFilesPaths()))
		isMatch, err := fileutils.JsonEqual(reader.GetFilesPaths()[0], filepath.Join(testDataPath, "load_missing_props_expected_results.json"))
		assert.NoError(t, err)
		assert.True(t, isMatch)
		utils.MaxBufferSize = 3
	}
	testDataPath := getBaseTestDir(t)
	notSortedWithProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_nosorted_by_created_withprops.json"), content.DefaultKey)
	sortedNoProps := content.NewContentReader(filepath.Join(testDataPath, "load_missing_props_sorted_by_created_noprops.json"), content.DefaultKey)
	reader, err := loadMissingProperties(sortedNoProps, notSortedWithProps)
	defer readerCloseAndAssert(t, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(reader.GetFilesPaths()))
	isMatch, err := fileutils.JsonEqual(reader.GetFilesPaths()[0], filepath.Join(testDataPath, "load_missing_props_by_created_expected_results.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	utils.MaxBufferSize = 3
}

func TestFilterBuildAqlSearchResults(t *testing.T) {
	testDataPath := getBaseTestDir(t)
	resultsToFilter := content.NewContentReader(filepath.Join(testDataPath, "filter_build_aql_search.json"), content.DefaultKey)
	buildArtifactsSha := map[string]int{"a": 2, "b": 2, "c": 2}
	resultReader, err := filterBuildAqlSearchResults(resultsToFilter, buildArtifactsSha, []Build{{"myBuild", "1"}})
	defer readerCloseAndAssert(t, resultReader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err := fileutils.JsonEqual(resultReader.GetFilesPaths()[0], filepath.Join(testDataPath, "filter_build_aql_search_expected.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
}

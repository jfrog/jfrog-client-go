package utils

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestGetFullUrl(t *testing.T) {
	assertUrl("Repo", "some/path", "name", "Repo/some/path/name", t)
	assertUrl("", "some/path", "name", "some/path/name", t)
	assertUrl("Repo", "", "name", "Repo/name", t)
	assertUrl("Repo", "some/path", "", "Repo/some/path", t)
	assertUrl("", "some/path", "", "some/path", t)
	assertUrl("", "", "", "", t)
}

func assertUrl(repo, path, name, fullUrl string, t *testing.T) {
	testItem := ResultItem{Repo: repo, Path: path, Name: name}
	if fullUrl != testItem.GetItemRelativePath() {
		t.Error("Unexpected URL built. Expected: `" + fullUrl + "` Got `" + testItem.GetItemRelativePath() + "`")
	}
}

func TestReduceTopChainDirResult(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	var reader, resultReader *content.ContentReader
	var isMatch bool

	// Single folder.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step1.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	isMatch, err = fileutils.FilesIdentical(filepath.Join(testDataPath, "reduce_top_chain_results_a.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())

	// Two different folders not sorted.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step2.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	isMatch, err = fileutils.FilesIdentical(filepath.Join(testDataPath, "reduce_top_chain_results_b.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())

	// One folder contains another, should reduce results.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step3.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	isMatch, err = fileutils.FilesIdentical(filepath.Join(testDataPath, "reduce_top_chain_results_b.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())

	oldMaxSize := utils.MaxBufferSize
	defer func() { utils.MaxBufferSize = oldMaxSize }()
	//Test buffer + sort
	utils.MaxBufferSize = 3
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step4.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	isMatch, err = fileutils.FilesIdentical(filepath.Join(testDataPath, "reduce_top_chain_results_c.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())

	//Two files in the same folder and one is a prefix to another.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step5.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	isMatch, err = fileutils.FilesIdentical(filepath.Join(testDataPath, "reduce_top_chain_results_d.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())

	//Two files in the same folder and one is a prefix to another and their folder.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step6.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	isMatch, err = fileutils.FilesIdentical(filepath.Join(testDataPath, "reduce_top_chain_results_e.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())
}

func TestReduceTopChainDirResultNoResults(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	reader := content.NewContentReader(filepath.Join(testDataPath, "no_results.json"), content.DefaultKey)
	resultReader, err := ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.True(t, resultReader.IsEmpty())
}

func TestReduceTopChainDirResultEmptyRepo(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	reader := content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_empty_repo.json"), content.DefaultKey)
	resultReader, err := ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.True(t, resultReader.IsEmpty())
	assert.NoError(t, resultReader.Close())
}

func TestReduceBottomChainDirResult(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	oldMaxSize := utils.MaxBufferSize
	defer func() { utils.MaxBufferSize = oldMaxSize }()
	for i := 0; i < 2; i++ {
		testResult := []int{1, 2, 2, 2, 3}
		for i := 1; i <= 5; i++ {
			reader := content.NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%v.json", i)), content.DefaultKey)
			resultReader, err := ReduceBottomChainDirResult(ResultItem{}, reader)
			assert.NoError(t, err)
			isMatch, err := fileutils.FilesIdentical(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%vresults.json", testResult[i-1])), resultReader.GetFilePath())
			assert.NoError(t, err)
			assert.True(t, isMatch)
			if isMatch == false {
				l, _ := resultReader.Length()
				log.Debug(fmt.Sprintf("reduce_bottom_chain_step%v.json  length: %v name %v", i, l, resultReader.GetFilePath()))
			} else {
				assert.NoError(t, resultReader.Close())
			}
		}
		utils.MaxBufferSize = 2
	}
}

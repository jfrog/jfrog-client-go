package utils

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
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

	// Single folder.
	cr := content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step1.json"), "results")
	resultReader, err := ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err := fileutils.FilesMath(filepath.Join(testDataPath, "reduce_top_chain_results_a.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)

	// Two different folders not sorted.
	cr = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step2.json"), "results")
	resultReader, err = ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err = fileutils.FilesMath(filepath.Join(testDataPath, "reduce_top_chain_results_b.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)

	// One folder contains another, should reduce results.
	cr = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step3.json"), "results")
	resultReader, err = ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err = fileutils.FilesMath(filepath.Join(testDataPath, "reduce_top_chain_results_b.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)

	//Test buffer + sort
	utils.MAX_BUFFER_SIZE = 3
	cr = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step4.json"), "results")
	resultReader, err = ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err = fileutils.FilesMath(filepath.Join(testDataPath, "reduce_top_chain_results_c.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)

	//Two files in the same folder and one is a prefix to another.
	cr = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step5.json"), "results")
	resultReader, err = ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err = fileutils.FilesMath(filepath.Join(testDataPath, "reduce_top_chain_results_d.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)

	//Two files in the same folder and one is a prefix to another and their folder.
	cr = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step6.json"), "results")
	resultReader, err = ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err = fileutils.FilesMath(filepath.Join(testDataPath, "reduce_top_chain_results_e.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestReduceTopChainDirResultNoResults(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	cr := content.NewContentReader(filepath.Join(testDataPath, "no_results.json"), "results")
	resultReader, err := ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err := fileutils.FilesMath(filepath.Join(testDataPath, "no_results.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestReduceTopChainDirResultEmptyRepo(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	cr := content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_empty_repo.json"), "results")
	resultReader, err := ReduceTopChainDirResult(cr)
	assert.NoError(t, err)
	result, err := fileutils.FilesMath(filepath.Join(testDataPath, "no_results.json"), resultReader.GetFilePath())
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestReduceBottomChainDirResult(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	testResult := []int{1, 2, 2, 2, 3}
	for i := 1; i <= 5; i++ {
		cr := content.NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%v.json", i)), "results")
		resultReader, err := ReduceBottomChainDirResult(cr)
		assert.NoError(t, err)
		result, err := fileutils.FilesMath(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%vresults.json", testResult[i-1])), resultReader.GetFilePath())
		assert.NoError(t, err)
		assert.True(t, result)
	}
}

func TestMergeSortedFiles(t *testing.T) {
	testDataPath, err := getBaseTestDir()
	assert.NoError(t, err)
	var sortedFiles []*content.ContentReader
	for i := 1; i <= 3; i++ {
		sortedFiles = append(sortedFiles, content.NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("buffer_file_1_%v.json", i)), "results"))
	}
	cr, err := MergeSortedFiles(sortedFiles)
	assert.NoError(t, err)
	result, err := fileutils.FilesMath(cr.GetFilePath(), filepath.Join(testDataPath, "merged_buffer_1.json"))
	assert.NoError(t, err)
	assert.True(t, result)
}

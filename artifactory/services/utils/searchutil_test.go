package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
	dir, _ := os.Getwd()
	testDataPath := filepath.Join(dir, "reduce_top_chain_tests", "testsdata", "reducedirresult")
	testResult := []int{1, 2, 2, 3, 3, 3, 3, 4}
	for i := 1; i <= 8; i++ {
		cr := content.NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("step%v.json", i)), "results")
		resultReader, err := ReduceTopChainDirResult(cr)
		assert.NoError(t, err)
		assert.True(t, filesMath(t, filepath.Join(testDataPath, fmt.Sprintf("reduce_top_chain_step%vresults.json", testResult[i-1])), resultReader.GetFilePath()))
	}
}

func TestReduceBottomChainDirResult(t *testing.T) {
	dir, _ := os.Getwd()
	testDataPath := filepath.Join(dir, "tests", "testsdata", "reducedirresult")
	testResult := []int{1, 2, 2, 2, 3}
	for i := 1; i <= 5; i++ {
		cr := content.NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%v.json", i)), "results")
		resultReader, err := ReduceBottomChainDirResult(cr)
		assert.NoError(t, err)
		assert.True(t, filesMath(t, filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%vresults.json", testResult[i-1])), resultReader.GetFilePath()))
	}
}

func filesMath(t *testing.T, srcPath string, toComparePath string) bool {
	srcDetails, err := fileutils.GetFileDetails(srcPath)
	assert.NoError(t, err)
	toCompareDetails, err := fileutils.GetFileDetails(toComparePath)
	assert.NoError(t, err)
	return srcDetails.Checksum.Md5 == toCompareDetails.Checksum.Md5
}

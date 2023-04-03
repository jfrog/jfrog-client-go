package utils

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/gofrog/version"

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
	testDataPath := getBaseTestDir(t)
	var reader, resultReader *content.ContentReader
	var isMatch bool

	// Single folder.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step1.json"), content.DefaultKey)
	resultReader, err := ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err = fileutils.JsonEqual(filepath.Join(testDataPath, "reduce_top_chain_results_a.json"), resultReader.GetFilesPaths()[0])
	assert.NoError(t, err)
	assert.True(t, isMatch)
	readerCloseAndAssert(t, resultReader)

	// Two different folders not sorted.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step2.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err = fileutils.JsonEqual(filepath.Join(testDataPath, "reduce_top_chain_results_b.json"), resultReader.GetFilesPaths()[0])
	assert.NoError(t, err)
	assert.True(t, isMatch)
	readerCloseAndAssert(t, resultReader)

	// One folder contains another, should reduce results.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step3.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err = fileutils.JsonEqual(filepath.Join(testDataPath, "reduce_top_chain_results_b.json"), resultReader.GetFilesPaths()[0])
	assert.NoError(t, err)
	assert.True(t, isMatch)
	readerCloseAndAssert(t, resultReader)

	oldMaxSize := utils.MaxBufferSize
	defer func() { utils.MaxBufferSize = oldMaxSize }()
	// Test buffer + sort
	utils.MaxBufferSize = 3
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step4.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err = fileutils.JsonEqual(filepath.Join(testDataPath, "reduce_top_chain_results_c.json"), resultReader.GetFilesPaths()[0])
	assert.NoError(t, err)
	assert.True(t, isMatch)
	readerCloseAndAssert(t, resultReader)

	// Two files in the same folder and one is a prefix to another.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step5.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err = fileutils.JsonEqual(filepath.Join(testDataPath, "reduce_top_chain_results_d.json"), resultReader.GetFilesPaths()[0])
	assert.NoError(t, err)
	assert.True(t, isMatch)
	readerCloseAndAssert(t, resultReader)

	// Two files in the same folder and one is a prefix to another and their folder.
	reader = content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_step6.json"), content.DefaultKey)
	resultReader, err = ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err = fileutils.JsonEqual(filepath.Join(testDataPath, "reduce_top_chain_results_e.json"), resultReader.GetFilesPaths()[0])
	assert.NoError(t, err)
	assert.True(t, isMatch)
	readerCloseAndAssert(t, resultReader)
}

func TestReduceTopChainDirResultNoResults(t *testing.T) {
	testDataPath := getBaseTestDir(t)
	reader := content.NewContentReader(filepath.Join(testDataPath, "no_results.json"), content.DefaultKey)
	resultReader, err := ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.True(t, resultReader.IsEmpty())
}

func TestReduceTopChainDirResultEmptyRepo(t *testing.T) {
	testDataPath := getBaseTestDir(t)
	reader := content.NewContentReader(filepath.Join(testDataPath, "reduce_top_chain_empty_repo.json"), content.DefaultKey)
	resultReader, err := ReduceTopChainDirResult(ResultItem{}, reader)
	assert.NoError(t, err)
	assert.True(t, resultReader.IsEmpty())
	readerCloseAndAssert(t, resultReader)
}

func TestReduceBottomChainDirResult(t *testing.T) {
	testDataPath := getBaseTestDir(t)
	oldMaxSize := utils.MaxBufferSize
	defer func() { utils.MaxBufferSize = oldMaxSize }()
	for i := 0; i < 2; i++ {
		testResult := []int{1, 2, 2, 2, 3}
		for i := 1; i <= 5; i++ {
			reader := content.NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%v.json", i)), content.DefaultKey)
			resultReader, err := ReduceBottomChainDirResult(ResultItem{}, reader)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
			isMatch, err := fileutils.JsonEqual(filepath.Join(testDataPath, fmt.Sprintf("reduce_bottom_chain_step%vresults.json", testResult[i-1])), resultReader.GetFilesPaths()[0])
			assert.NoError(t, err)
			assert.True(t, isMatch)
			if isMatch == false {
				l, _ := resultReader.Length()
				log.Debug(fmt.Sprintf("reduce_bottom_chain_step%v.json  length: %v name %v", i, l, resultReader.GetFilesPaths()))
			} else {
				readerCloseAndAssert(t, resultReader)
			}
		}
		utils.MaxBufferSize = 2
	}
}

func TestToBuildInfoArtifact(t *testing.T) {
	data := []struct {
		artifact ArtifactDetails
		res      *entities.Artifact
	}{
		{ArtifactDetails{
			ArtifactoryPath: "repo/art/text.txt",
			Checksums:       entities.Checksum{Sha1: "1", Md5: "2", Sha256: "3"},
		}, &entities.Artifact{
			Name:     "text.txt",
			Type:     "txt",
			Path:     "art/text.txt",
			Checksum: entities.Checksum{Sha1: "1", Md5: "2", Sha256: "3"},
		}},
		{ArtifactDetails{
			ArtifactoryPath: "text",
		}, nil},
	}

	for _, d := range data {
		got, err := d.artifact.ToBuildInfoArtifact()
		if d.res == nil {
			assert.Error(t, err)
			continue
		}
		assert.Equal(t, d.res.Type, got.Type)
		assert.Equal(t, d.res.Name, got.Name)
		assert.Equal(t, d.res.Path, got.Path)
		assert.Equal(t, d.res.Md5, got.Md5)
		assert.Equal(t, d.res.Sha1, got.Sha1)
		assert.Equal(t, d.res.Sha256, got.Sha256)
	}
}

func TestValidateTransitiveSearchAllowed(t *testing.T) {
	testRuns := []struct {
		params             *CommonParams
		artifactoryVersion *version.Version
		expectedTransitive bool
	}{
		{&CommonParams{Transitive: true}, version.NewVersion("7.0.0"), false},
		{&CommonParams{Transitive: true}, version.NewVersion("7.17.0"), true},
		{&CommonParams{Transitive: true}, version.NewVersion("7.17.0-m029"), true},
		{&CommonParams{Transitive: true}, version.NewVersion("7.19.0"), true},
		{&CommonParams{Transitive: false}, version.NewVersion("7.0.0"), false},
		{&CommonParams{Transitive: false}, version.NewVersion("7.17.0"), false},
		{&CommonParams{Transitive: false}, version.NewVersion("7.17.0-m029"), false},
		{&CommonParams{Transitive: false}, version.NewVersion("7.19.0"), false},
	}
	for _, test := range testRuns {
		t.Run(fmt.Sprintf("transitive:%t,version:%s", test.params.Transitive, test.artifactoryVersion.GetVersion()), func(t *testing.T) {
			DisableTransitiveSearchIfNotAllowed(test.params, test.artifactoryVersion)
			assert.Equal(t, test.expectedTransitive, test.params.Transitive)
		})
	}
}

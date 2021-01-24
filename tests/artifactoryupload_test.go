package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryUpload(t *testing.T) {
	t.Run("flat", flatUpload)
	t.Run("recursive", recursiveUpload)
	t.Run("placeholder", placeholderUpload)
	t.Run("includeDirs", includeDirsUpload)
	t.Run("explode", explodeUpload)
	t.Run("props", propsUpload)
}

func flatUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer os.RemoveAll(workingDir)

	pattern := FixWinPath(filepath.Join(workingDir, "out", "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo}
	up.Flat = true
	uploaded, failed, err := testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	searchParams.Pattern = RtTargetRepo
	reader, err := testsSearchService.Search(searchParams)
	defer reader.Close()
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "." {
			t.Error("Expected path to be root due to using the flat flag.", "Got:", item.Path)
		}
	}
	assert.NoError(t, reader.GetError())
	length, err := reader.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func recursiveUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer os.RemoveAll(workingDir)

	pattern := FixWinPath(filepath.Join(workingDir, "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo}
	up.Flat = true
	uploaded, failed, err := testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	searchParams.Pattern = RtTargetRepo
	reader, err := testsSearchService.Search(searchParams)
	defer reader.Close()
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "." {
			t.Error("Expected path to be root(flat by default).", "Got:", item.Path)
		}
		if item.Name != "a.in" {
			t.Error("Missing File a.in")
		}
	}
	assert.NoError(t, reader.GetError())
	length, err := reader.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func placeholderUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer os.RemoveAll(workingDir)

	pattern := FixWinPath(filepath.Join(workingDir, "(*).in"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "{1}"}
	up.Flat = true
	uploaded, failed, err := testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	searchParams.Pattern = RtTargetRepo
	reader, err := testsSearchService.Search(searchParams)
	defer reader.Close()
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "out" {
			t.Error("Expected path to be out.", "Got:", item.Path)
		}
		if item.Name != "a" {
			t.Error("Missing File a")
		}
	}
	assert.NoError(t, reader.GetError())
	length, err := reader.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func includeDirsUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer os.RemoveAll(workingDir)

	pattern := FixWinPath(filepath.Join(workingDir, "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, IncludeDirs: true, Recursive: false, Target: RtTargetRepo}
	up.Flat = true
	uploaded, failed, err := testsUploadService.UploadFiles(up)
	if uploaded != 0 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	searchParams.Pattern = RtTargetRepo
	searchParams.IncludeDirs = true
	reader, err := testsSearchService.Search(searchParams)
	defer reader.Close()
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Name == "." {
			continue
		}
		if item.Path != "." {
			t.Error("Expected path to be root(flat by default).", "Got:", item.Path)
		}
		if item.Name != "out" {
			t.Error("Missing directory out.")
		}
	}
	assert.NoError(t, reader.GetError())
	length, err := reader.Length()
	assert.NoError(t, err)
	if length < 2 {
		t.Error("Expected to get at least two items, default and the out folder.")
	}
	artifactoryCleanup(t)
}

func explodeUpload(t *testing.T) {
	workingDir, filePath := createWorkingDir(t)
	defer os.RemoveAll(workingDir)

	err := fileutils.ZipFolderFiles(filePath, filepath.Join(workingDir, "zipFile.zip"))
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
	}
	pattern := FixWinPath(filepath.Join(workingDir, "*.zip"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, IncludeDirs: true, Recursive: false, Target: RtTargetRepo}
	up.Flat = true
	up.ExplodeArchive = true
	uploaded, failed, err := testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	searchParams.Pattern = RtTargetRepo
	searchParams.IncludeDirs = true
	reader, err := testsSearchService.Search(searchParams)
	defer reader.Close()
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Name == "." {
			continue
		}
		if item.Name != "a.in" {
			t.Error("Missing file a.in")
		}
	}
	assert.NoError(t, reader.GetError())
	length, err := reader.Length()
	assert.NoError(t, err)
	if length < 2 {
		t.Error("Expected to get at least two items, default and the out folder.")
	}
	artifactoryCleanup(t)
}

func propsUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer os.RemoveAll(workingDir)

	// Upload a.in with property key1=val1
	pattern := FixWinPath(filepath.Join(workingDir, "out", "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Target: RtTargetRepo, TargetProps: "key1=val1"}
	up.Flat = true
	uploaded, failed, err := testsUploadService.UploadFiles(up)
	assert.Equal(t, 1, uploaded)
	assert.Equal(t, 0, failed)
	assert.NoError(t, err)

	// Search a.in with property key1=val1
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	searchParams.Pattern = RtTargetRepo
	searchParams.Props = "key1=val1"
	reader, err := testsSearchService.Search(searchParams)
	defer reader.Close()
	if err != nil {
		t.Error(err)
	}
	length, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, 1, length)

	// Assert property key and value exist in the search results
	item := new(utils.ResultItem)
	err = reader.NextRecord(item)
	assert.NoError(t, err)
	assert.Len(t, item.Properties, 1)
	assert.Equal(t, "key1", item.Properties[0].Key)
	assert.Equal(t, "val1", item.Properties[0].Value)

	artifactoryCleanup(t)
}

func createWorkingDir(t *testing.T) (string, string) {
	workingDir, relativePath, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return workingDir, relativePath
}

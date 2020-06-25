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
	assert.NoError(t, fileutils.CleanupReaderWriterTempFilesAndDirs())
}

func flatUpload(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := FixWinPath(filepath.Join(workingDir, "out", "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo}
	up.Flat = true
	_, uploaded, failed, err := testsUploadService.UploadFiles(up)
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
	cr, err := testsSearchService.Search(searchParams)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); cr.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "." {
			t.Error("Expected path to be root due to using the flat flag.", "Got:", item.Path)
		}
	}
	assert.NoError(t, cr.GetError())
	length, err := cr.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func recursiveUpload(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := FixWinPath(filepath.Join(workingDir, "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo}
	up.Flat = true
	_, uploaded, failed, err := testsUploadService.UploadFiles(up)
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
	cr, err := testsSearchService.Search(searchParams)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); cr.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "." {
			t.Error("Expected path to be root(flat by default).", "Got:", item.Path)
		}
		if item.Name != "a.in" {
			t.Error("Missing File a.in")
		}
	}
	assert.NoError(t, cr.GetError())
	length, err := cr.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func placeholderUpload(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := FixWinPath(filepath.Join(workingDir, "(*).in"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "{1}"}
	up.Flat = true
	_, uploaded, failed, err := testsUploadService.UploadFiles(up)
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
	cr, err := testsSearchService.Search(searchParams)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); cr.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "out" {
			t.Error("Expected path to be out.", "Got:", item.Path)
		}
		if item.Name != "a" {
			t.Error("Missing File a")
		}
	}
	assert.NoError(t, cr.GetError())
	length, err := cr.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func includeDirsUpload(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := FixWinPath(filepath.Join(workingDir, "*"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, IncludeDirs: true, Recursive: false, Target: RtTargetRepo}
	up.Flat = true
	_, uploaded, failed, err := testsUploadService.UploadFiles(up)
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
	cr, err := testsSearchService.Search(searchParams)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); cr.NextRecord(item) == nil; item = new(utils.ResultItem) {
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
	assert.NoError(t, cr.GetError())
	length, err := cr.Length()
	assert.NoError(t, err)
	if length < 2 {
		t.Error("Expected to get at least two items, default and the out folder.")
	}
	artifactoryCleanup(t)
}

func explodeUpload(t *testing.T) {
	workingDir, filePath, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	err = fileutils.ZipFolderFiles(filePath, filepath.Join(workingDir, "zipFile.zip"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = os.Remove(filePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	pattern := FixWinPath(filepath.Join(workingDir, "*.zip"))
	up := services.NewUploadParams()
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, IncludeDirs: true, Recursive: false, Target: RtTargetRepo}
	up.Flat = true
	up.ExplodeArchive = true
	_, uploaded, failed, err := testsUploadService.UploadFiles(up)
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
	cr, err := testsSearchService.Search(searchParams)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); cr.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Name == "." {
			continue
		}
		if item.Name != "a.in" {
			t.Error("Missing file a.in")
		}
	}
	assert.NoError(t, cr.GetError())
	length, err := cr.Length()
	assert.NoError(t, err)
	if length < 2 {
		t.Error("Expected to get at least two items, default and the out folder.")
	}
	artifactoryCleanup(t)
}

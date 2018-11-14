package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArtifactoryDownload(t *testing.T) {
	uploadDummyFile(t)
	t.Run("flat", flatDownload)
	t.Run("recursive", recursiveDownload)
	t.Run("placeholder", placeholderDownload)
	t.Run("includeDirs", includeDirsDownload)
	artifactoryCleanUp(t)
}

func uploadDummyFile(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := filepath.Join(workingDir, "*")
	pattern = strings.Replace(pattern, "\\", "\\\\", -1)
	up := &services.UploadParamsImp{}
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "test/"}
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
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: pattern, Recursive: true, Target: RtTargetRepo + "b.in"}
	up.Flat = true
	_, uploaded, failed, err = testsUploadService.UploadFiles(up)
	if uploaded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if failed != 0 {
		t.Error("Failed to upload", failed, "files.")
	}
	if err != nil {
		t.Error(err)
	}
}

func flatDownload(t *testing.T) {
	var err error
	workingDir, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir)
	downloadPattern := RtTargetRepo + "*"
	downloadTarget := workingDir + string(filepath.Separator)
	_, _, err = testsDownloadService.DownloadFiles(&services.DownloadParamsImpl{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("Missing file a.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("Missing file b.in")
	}

	workingDir2, err := ioutil.TempDir("", "downloadTests")
	downloadTarget = workingDir2 + string(filepath.Separator)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir2)
	_, _, err = testsDownloadService.DownloadFiles(&services.DownloadParamsImpl{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: false})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir2, "test", "a.in"), false) {
		t.Error("Missing file a.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir2, "b.in"), false) {
		t.Error("Missing file b.in")
	}
}

func recursiveDownload(t *testing.T) {
	uploadDummyFile(t)
	var err error
	workingDir, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir)
	downloadPattern := RtTargetRepo + "*"
	downloadTarget := workingDir + string(filepath.Separator)
	_, _, err = testsDownloadService.DownloadFiles(&services.DownloadParamsImpl{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("Missing file a.in")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("Missing file b.in")
	}

	workingDir2, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir2)
	downloadTarget = workingDir2 + string(filepath.Separator)
	_, _, err = testsDownloadService.DownloadFiles(&services.DownloadParamsImpl{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: false, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if fileutils.IsPathExists(filepath.Join(workingDir2, "a.in"), false) {
		t.Error("Should not download a.in")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir2, "b.in"), false) {
		t.Error("Missing file b.in")
	}
}

func placeholderDownload(t *testing.T) {
	uploadDummyFile(t)
	var err error
	workingDir, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir)
	downloadPattern := RtTargetRepo + "(*).in"
	downloadTarget := workingDir + string(filepath.Separator) + "{1}" + string(filepath.Separator)
	_, _, err = testsDownloadService.DownloadFiles(&services.DownloadParamsImpl{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "test", "a", "a.in"), false) {
		t.Error("Missing file a.in")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir, "b", "b.in"), false) {
		t.Error("Missing file b.in")
	}
}

func includeDirsDownload(t *testing.T) {
	var err error
	workingDir, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir)
	downloadPattern := RtTargetRepo + "*"
	downloadTarget := workingDir + string(filepath.Separator)
	_, _, err = testsDownloadService.DownloadFiles(&services.DownloadParamsImpl{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, IncludeDirs: true, Recursive: false, Target: downloadTarget}, Flat: false})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "test"), false) {
		t.Error("Missing test folder")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("Missing file b.in")
	}
}

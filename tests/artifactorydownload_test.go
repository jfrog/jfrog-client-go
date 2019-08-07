package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestArtifactoryDownload(t *testing.T) {
	uploadDummyFile(t)
	t.Run("flat", flatDownload)
	t.Run("recursive", recursiveDownload)
	t.Run("placeholder", placeholderDownload)
	t.Run("includeDirs", includeDirsDownload)
	t.Run("excludePatterns", excludePatternsDownload)
	t.Run("explodeArchive", explodeArchiveDownload)
	artifactoryCleanup(t)
}

func uploadDummyFile(t *testing.T) {
	workingDir, _, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(workingDir)
	pattern := FixWinPath(filepath.Join(workingDir, "*"))
	up := services.NewUploadParams()
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
	archivePath := filepath.Join(workingDir, "c.tar.gz")
	err = archiver.TarGz.Make(archivePath, []string{filepath.Join(workingDir, "out/a.in")})
	if err != nil {
		t.Error(err)
	}
	up.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{Pattern: archivePath, Recursive: true, Target: RtTargetRepo}
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
	// Download all from TargetRepo with flat = true
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("Missing file a.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("Missing file b.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "c.tar.gz"), false) {
		t.Error("Missing file c.tar.gz")
	}

	workingDir2, err := ioutil.TempDir("", "downloadTests")
	downloadTarget = workingDir2 + string(filepath.Separator)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir2)
	// Download all from TargetRepo with flat = false
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: false})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir2, "test", "a.in"), false) {
		t.Error("Missing file a.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir2, "b.in"), false) {
		t.Error("Missing file b.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir2, "c.tar.gz"), false) {
		t.Error("Missing file c.tar.gz")
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
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("Missing file a.in")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("Missing file b.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "c.tar.gz"), false) {
		t.Error("Missing file c.tar.gz")
	}

	workingDir2, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir2)
	downloadTarget = workingDir2 + string(filepath.Separator)
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: false, Target: downloadTarget}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if fileutils.IsPathExists(filepath.Join(workingDir2, "a.in"), false) {
		t.Error("Should not download a.in")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir2, "b.in"), false) {
		t.Error("Missing file b.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir2, "c.tar.gz"), false) {
		t.Error("Missing file c.tar.gz")
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
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true})
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
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, IncludeDirs: true, Recursive: false, Target: downloadTarget}, Flat: false})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "test"), false) {
		t.Error("Missing test folder")
	}

	if !fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("Missing file b.in")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "c.tar.gz"), false) {
		t.Error("Missing file c.tsr.gz")
	}
}

func excludePatternsDownload(t *testing.T) {
	workingDir, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir)
	downloadPattern := RtTargetRepo + "*"
	downloadTarget := workingDir + string(filepath.Separator)
	excludePatterns := []string{"b.in", "*.tar.gz"}
	_, _, err = testsDownloadService.DownloadFiles(services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget, ExcludePatterns: excludePatterns}, Flat: true})
	if err != nil {
		t.Error(err)
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("Missing file a.in")
	}

	if fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("File b.in should have been excluded")
	}
	if fileutils.IsPathExists(filepath.Join(workingDir, "c.tar.gz"), false) {
		t.Error("File c.tar.gz should have been excluded")
	}
}

func explodeArchiveDownload(t *testing.T) {
	workingDir, err := ioutil.TempDir("", "downloadTests")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(workingDir)
	downloadPattern := RtTargetRepo + "*.tar.gz"
	downloadTarget := workingDir + string(filepath.Separator)
	downloadParams := services.DownloadParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{Pattern: downloadPattern, Recursive: true, Target: downloadTarget}, Flat: true, Explode: false}
	// First we'll download c.tar.gz without extracting it (explode = false by default).
	_, _, err = testsDownloadService.DownloadFiles(downloadParams)
	if err != nil {
		t.Error(err)
	}
	if fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("File a.in should not have been downloaded")
	}
	if fileutils.IsPathExists(filepath.Join(workingDir, "b.in"), false) {
		t.Error("File b.in should not have been downloaded")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "c.tar.gz"), false) {
		t.Error("Missing file c.tar.gz")
	}

	// Scenario 1:  Download c.tar.gz with explode = true, when it already exists in the target dir.
	// Artifactory should perform "checksum download" and not actually downloading it, but still need to extract it.
	downloadParams.Explode = true
	explodeDownloadAndVerify(t, &downloadParams, workingDir)

	// Remove the download target dir.
	err = os.RemoveAll(workingDir)
	if err != nil {
		t.Error(err)
	}
	// Scenario 2: Download c.tar.gz with explode = true, when it does not exist in the target dir.
	// Artifactory should download the file and extract it.
	explodeDownloadAndVerify(t, &downloadParams, workingDir)
}

func explodeDownloadAndVerify(t *testing.T, downloadParams *services.DownloadParams, workingDir string) {
	_, _, err := testsDownloadService.DownloadFiles(*downloadParams)
	if err != nil {
		t.Error(err)
	}
	if fileutils.IsPathExists(filepath.Join(workingDir, "c.tar.gz"), false) {
		t.Error("File c.tar.gz should have been extracted")
	}
	if !fileutils.IsPathExists(filepath.Join(workingDir, "a.in"), false) {
		t.Error("Missing file a.in")
	}
}

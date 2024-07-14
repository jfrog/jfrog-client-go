package fileutils

import (
	biutils "github.com/jfrog/build-info-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSsh(t *testing.T) {
	testRuns := []struct {
		url      string
		expected bool
	}{
		{"http://some.url", false},
		{"https://some.url", false},
		{"sshd://wrong.url", false},
		{"assh://wrong.url", false},
		{"ssh://some.url", true},
		{"sSh://some.url/some/api", true},
		{"SSH://some.url/some/api", true},
	}
	for _, test := range testRuns {
		t.Run(test.url, func(t *testing.T) {
			assert.Equal(t, test.expected, IsSshUrl(test.url), "Wrong ssh for URL: "+test.url)
		})
	}
}

func TestFindUpstreamFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		assert.Error(t, err)
		return
	}
	defer func() {
		assert.NoError(t, os.Chdir(wd))
	}()

	// CD into a directory with a goDotMod.test file.
	projectRoot := filepath.Join("testdata", "project")
	err = os.Chdir(projectRoot)
	if err != nil {
		assert.Error(t, err)
		return
	}

	// Make projectRoot an absolute path.
	projectRoot, err = os.Getwd()
	if err != nil {
		assert.Error(t, err)
		return
	}

	// Get the project root.
	if err = assertFindUpstreamExistsAndEqual(t, "goDotMod.test", projectRoot, File); err != nil {
		return
	}

	// Assert with Any too.
	if err = assertFindUpstreamExistsAndEqual(t, "goDotMod.test", projectRoot, Any); err != nil {
		return
	}

	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		assert.Error(t, err)
		return
	}

	// CD into a subdirectory in the same project, and expect to get the same project root.
	assert.NoError(t, os.Chdir(wd))
	projectSubDirectory := filepath.Join("testdata", "project", "dir")
	err = os.Chdir(projectSubDirectory)
	if err != nil {
		assert.Error(t, err)
		return
	}

	if err = assertFindUpstreamExistsAndEqual(t, "goDotMod.test", projectRoot, File); err != nil {
		return
	}

	root, exists, err := FindUpstream("go-missing.mod", File)
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.False(t, exists, "File go-missing.mod found but shouldn't.")
	assert.Empty(t, root, "File go-missing.mod shouldn't be found")

	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		assert.Error(t, err)
		return
	}

	// Now CD into a directory outside the project, and expect to get a different project root.
	noProjectRoot := filepath.Join("testdata", "noproject")
	err = os.Chdir(noProjectRoot)
	if err != nil {
		assert.Error(t, err)
		return
	}
	root, exists, err = FindUpstream("goDotMod.test", File)
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.True(t, exists, "File goDotMod.test is missing.")
	assert.NotEqual(t, projectRoot, root)
}

func TestFindUpstreamFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		assert.Error(t, err)
		return
	}
	defer func() {
		assert.NoError(t, os.Chdir(wd))
	}()

	// Create path to directory to find.
	dirPath := "testdata"
	err = os.Chdir(dirPath)
	if err != nil {
		assert.Error(t, err)
		return
	}
	// Get absolute path.
	dirPath, err = os.Getwd()
	if err != nil {
		assert.Error(t, err)
		return
	}
	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		assert.Error(t, err)
		return
	}

	// Go to starting dir to search from.
	searchFromDir := filepath.Join("testdata", "project", "dir")
	err = os.Chdir(searchFromDir)
	if err != nil {
		assert.Error(t, err)
		return
	}

	// Get the directory path.
	if err = assertFindUpstreamExistsAndEqual(t, "noproject", dirPath, Dir); err != nil {
		return
	}

	// Assert with Any too.
	if err = assertFindUpstreamExistsAndEqual(t, "noproject", dirPath, Any); err != nil {
		return
	}
}

func assertFindUpstreamExistsAndEqual(t *testing.T, path, expectedPath string, itemType ItemType) error {
	foundPath, exists, err := FindUpstream(path, itemType)
	if err != nil {
		assert.Error(t, err)
		return err
	}
	assert.True(t, exists)
	assert.Equal(t, expectedPath, foundPath)
	return nil
}

func TestIsEqualToLocalFile(t *testing.T) {
	localFilePath := filepath.Join("testdata", "files", "comparisonFile")

	// Get file actual details.
	localFileDetails, err := GetFileDetails(localFilePath, true)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	actualMd5 := localFileDetails.Checksum.Md5
	actualSha1 := localFileDetails.Checksum.Sha1
	tests := []struct {
		name           string
		localPath      string
		remoteMd5      string
		remoteSha1     string
		expectedResult bool
	}{
		{"realEquality", localFilePath, actualMd5, actualSha1, true},
		{"unequalPath", "non/existing/path", actualMd5, actualSha1, false},
		{"unequalChecksum", localFilePath, "wrongMd5", "wrongSha1", false},
		{"unequalMd5", localFilePath, "wrongMd5", actualSha1, false},
		{"unequalSha1", localFilePath, actualMd5, "wrongSha1", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isEqual, err := IsEqualToLocalFile(test.localPath, test.remoteMd5, test.remoteSha1)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			assert.Equal(t, test.expectedResult, isEqual)
		})
	}
}

func TestListFilesByFilterFunc(t *testing.T) {
	testDir := filepath.Join("testdata", "listextension")
	expected := []string{filepath.Join(testDir, "a.proj"),
		filepath.Join(testDir, "b.csproj"),
		filepath.Join(testDir, "someproj.csproj")}

	// List files with extension that satisfy the filter function.
	filterFunc := func(filePath string) (bool, error) {
		ext := strings.TrimLeft(filepath.Ext(filePath), ".")
		return regexp.MatchString(`.*proj$`, ext)
	}
	files, err := listFilesRecursiveWithFilterFunc(testDir, false, filterFunc)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.ElementsMatch(t, expected, files)
}

func TestGetFileAndDirFromPath(t *testing.T) {
	testRuns := []struct {
		path         string
		expectedFile string
		expectedDir  string
	}{
		{"a\\\\b\\\\c.in", "c.in", "a\\\\b"},
		{"a\\b\\c.in", "c.in", "a\\b"},
		{"a/b/c.in", "c.in", "a/b"},
		{"a\\\\b\\\\", "", "a\\\\b"},
		{"", "", ""},
		{"a\\\\b\\c.in", "c.in", "a\\\\b"},
		{"a\\b\\\\c.in", "c.in", "a\\b"},
		{"\\c.in", "c.in", ""},
		{"\\\\c.in", "c.in", ""},
	}
	for _, test := range testRuns {
		File, Dir := GetFileAndDirFromPath(test.path)
		assert.Equal(t, test.expectedFile, File, "Wrong file name for path: "+test.path)
		assert.Equal(t, test.expectedDir, Dir, "Wrong dir for path: "+test.path)
	}
}

func TestRemoveDirContents(t *testing.T) {
	// Prepare the test environment in a temporary directory
	tmpDirPath, err := CreateTempDir()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, RemoveTempDir(tmpDirPath))
	}()
	err = biutils.CopyDir(filepath.Join("testdata", "removedircontents"), tmpDirPath, true, nil)
	assert.NoError(t, err)

	// Run the function
	dirToEmptyPath := filepath.Join(tmpDirPath, "dirtoempty")
	err = RemoveDirContents(dirToEmptyPath)
	assert.NoError(t, err)

	// Assert the directories contents: dirtoempty should be empty and dirtoremain should contain one file.
	emptyDirFiles, err := os.ReadDir(dirToEmptyPath)
	assert.NoError(t, err)
	assert.Empty(t, emptyDirFiles)
	remainedDirPath := filepath.Join(tmpDirPath, "dirtoremain")
	remainedDirFiles, err := os.ReadDir(remainedDirPath)
	assert.NoError(t, err)
	assert.Len(t, remainedDirFiles, 1)
}

func TestListFilesRecursiveWalkIntoDirSymlink(t *testing.T) {
	if io.IsWindows() {
		t.Skip("Running on windows, skipping...")
	}
	expectedFileList := []string{
		"testdata/dirsymlinks",
		"testdata/dirsymlinks/d1",
		"testdata/dirsymlinks/d1/File_F1",
		"testdata/dirsymlinks/d1/linktoparent",
		"testdata/dirsymlinks/d1/linktoparent/d1",
		"testdata/dirsymlinks/d1/linktoparent/d1/File_F1",
		"testdata/dirsymlinks/d1/linktoparent/d2",
		"testdata/dirsymlinks/d1/linktoparent/d2/d1link",
		"testdata/dirsymlinks/d1/linktoparent/d2/d1link/File_F1",
		"testdata/dirsymlinks/d2",
	}

	// This directory and its subdirectories contain a symlink to a parent directory and a symlink to a sibling directory.
	testDirPath := filepath.Join("testdata", "dirsymlinks")
	filesList, err := ListFilesRecursiveWalkIntoDirSymlink(testDirPath, true)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(expectedFileList, filesList))
}

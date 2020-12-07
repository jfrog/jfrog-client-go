package fileutils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSsh(t *testing.T) {
	tests := []struct {
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
	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			assert.Equal(t, test.expected, IsSshUrl(test.url), "Wrong ssh for URL: "+test.url)
		})
	}
}

func TestGetFileOrDirPathFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		assert.Error(t, err)
		return
	}
	defer os.Chdir(wd)

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
	root, exists, err := FindUpstream("goDotMod.test", File)
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.True(t, exists, "File goDotMod.test is missing.")
	assert.Equal(t, projectRoot, root)

	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		assert.Error(t, err)
		return
	}

	// CD into a sub directory in the same project, and expect to get the same project root.
	os.Chdir(wd)
	projectSubDirectory := filepath.Join("testdata", "project", "dir")
	err = os.Chdir(projectSubDirectory)
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
	assert.Equal(t, projectRoot, root)

	root, exists, err = FindUpstream("go-missing.mod", File)
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

func TestGetFileOrDirPathFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		assert.Error(t, err)
		return
	}
	defer os.Chdir(wd)

	// Create path to directory to find.
	dirPath := filepath.Join("testdata")
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
	root, exists, err := FindUpstream("noproject", Dir)
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.True(t, exists, "Dir noproject is missing.")
	assert.Equal(t, dirPath, root)
}

func TestIsEqualToLocalFile(t *testing.T) {
	localFilePath := filepath.Join("testdata", "files", "comparisonFile")

	// Get file actual details.
	localFileDetails, err := GetFileDetails(localFilePath)
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
	files, err := ListFilesByFilterFunc(testDir, filterFunc)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.ElementsMatch(t, expected, files)
}

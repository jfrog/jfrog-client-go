package utils

import (
	testsutils "github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

func TestVcsDetails(t *testing.T) {
	// Test the following .git types, on their corresponding paths in testdata.
	tests := []string{"vcs", "packedvcs", "submodule"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			var projectPath, tmpDir string
			if test == "submodule" {
				projectPath, tmpDir = testsutils.InitVcsSubmoduleTestDir(t, filepath.Join("testdata", test))
			} else {
				projectPath, tmpDir = initVcsTestDir(t, filepath.Join("testdata", test))
			}
			defer fileutils.RemoveTempDir(tmpDir)
			vcsDetails := NewVcsDetals()
			revision, url, branch, err := vcsDetails.GetVcsDetails(filepath.Join(projectPath))
			assert.NoError(t, err)
			assert.Equal(t, "https://github.com/jfrog/jfrog-cli.git", url)
			assert.Equal(t, "6198a6294722fdc75a570aac505784d2ec0d1818", revision)
			assert.Equal(t, "master", branch)
		})
	}
}

func initVcsTestDir(t *testing.T, srcPath string) (projectPath, tmpDir string) {
	var err error
	tmpDir, err = fileutils.CreateTempDir()
	assert.NoError(t, err)

	err = fileutils.CopyDir(srcPath, tmpDir, true, nil)
	assert.NoError(t, err)
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "gitdata"), false); found {
		assert.NoError(t, err)
		err := fileutils.RenamePath(filepath.Join(tmpDir, "gitdata"), filepath.Join(tmpDir, ".git"))
		assert.NoError(t, err)
	}
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "othergit", "gitdata"), false); found {
		assert.NoError(t, err)
		err := fileutils.RenamePath(filepath.Join(tmpDir, "othergit", "gitdata"), filepath.Join(tmpDir, "othergit", ".git"))
		assert.NoError(t, err)
	}
	projectPath, err = filepath.Abs(tmpDir)
	assert.NoError(t, err)
	return projectPath, tmpDir
}

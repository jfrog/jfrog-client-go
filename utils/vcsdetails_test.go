package utils

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

func TestVcsDetails(t *testing.T) {
	projectPath, tmpDir := initVcsTestDir(t, filepath.Join("testsdata", "vcs"))
	defer fileutils.RemoveTempDir(tmpDir)
	vcsDetals := NewVcsDetals()
	revision, url, err := vcsDetals.GetVcsDetails(filepath.Join(projectPath))
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/jfrog/jfrog-cli.git", url)
	assert.Equal(t, "d63c5957ad6819f4c02a817abe757f210d35ff92", revision)
}

func initVcsTestDir(t *testing.T, srcPath string) (projectPath, tmpDir string) {
	var err error
	tmpDir, err = fileutils.CreateTempDir()
	assert.NoError(t, err)

	err = fileutils.CopyDir(srcPath, tmpDir, true)
	assert.NoError(t, err)
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "gitdata"), false); found {
		assert.NoError(t, err)
		err := fileutils.RenamePath(filepath.Join(tmpDir, "gitdata"), filepath.Join(tmpDir, ".git"))
		assert.NoError(t, err)
	}
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "OtherGit", "gitdata"), false); found {
		assert.NoError(t, err)
		err := fileutils.RenamePath(filepath.Join(tmpDir, "OtherGit", "gitdata"), filepath.Join(tmpDir, "OtherGit", ".git"))
		assert.NoError(t, err)
	}
	projectPath, err = filepath.Abs(tmpDir)
	assert.NoError(t, err)
	return projectPath, tmpDir
}

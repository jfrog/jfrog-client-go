package utils

import (
	biutils "github.com/jfrog/build-info-go/utils"
	testsutils "github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

func TestVcsDetails(t *testing.T) {
	// Test the following .git types, on their corresponding paths in testdata.
	testRuns := []string{"vcs", "packedvcs", "submodule", "worktree"}
	for _, test := range testRuns {
		t.Run(test, func(t *testing.T) {
			var projectPath, tmpDir string
			// Create temp folder.
			tmpDir, err := fileutils.CreateTempDir()
			assert.NoError(t, err, "Couldn't create temp dir")
			defer func() {
				assert.NoError(t, fileutils.RemoveTempDir(tmpDir), "Couldn't remove temp dir")
			}()

			switch test {
			case "submodule":
				projectPath = testsutils.InitVcsSubmoduleTestDir(t, filepath.Join("testdata", test), tmpDir)
			case "worktree":
				projectPath = testsutils.InitVcsWorktreeTestDir(t, filepath.Join("testdata", test), tmpDir)
			default:
				projectPath = initVcsTestDir(t, filepath.Join("testdata", test), tmpDir)
			}
			vcsDetails := NewVcsDetails()
			revision, url, branch, err := vcsDetails.GetVcsDetails(projectPath)
			assert.NoError(t, err)
			assert.Equal(t, "https://github.com/jfrog/jfrog-cli.git", url)
			assert.Equal(t, "6198a6294722fdc75a570aac505784d2ec0d1818", revision)
			assert.Equal(t, "master", branch)
		})
	}
}

func initVcsTestDir(t *testing.T, srcPath, tmpDir string) (projectPath string) {
	var err error
	assert.NoError(t, biutils.CopyDir(srcPath, tmpDir, true, nil))
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "gitdata"), false); found {
		assert.NoError(t, err)
		assert.NoError(t, fileutils.RenamePath(filepath.Join(tmpDir, "gitdata"), filepath.Join(tmpDir, ".git")))
	}
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "othergit", "gitdata"), false); found {
		assert.NoError(t, err)
		assert.NoError(t, fileutils.RenamePath(filepath.Join(tmpDir, "othergit", "gitdata"), filepath.Join(tmpDir, "othergit", ".git")))
	}
	projectPath, err = filepath.Abs(tmpDir)
	assert.NoError(t, err)
	return projectPath
}

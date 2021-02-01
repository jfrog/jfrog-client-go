package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGitManager(t *testing.T) {
	// Test the following .git types, on their corresponding paths in testdata.
	tests := []string{"vcs", "packedVcs"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			projectPath, tmpDir := initVcsTestDir(t, filepath.Join("testdata", test))
			defer fileutils.RemoveTempDir(tmpDir)
			gitManager := NewGitManager(projectPath)
			err := gitManager.ReadConfig()
			assert.NoError(t, err)
			assert.Equal(t, "https://github.com/jfrog/jfrog-cli.git", gitManager.GetUrl())
			assert.Equal(t, "d63c5957ad6819f4c02a817abe757f210d35ff92", gitManager.GetRevision())
		})
	}
}

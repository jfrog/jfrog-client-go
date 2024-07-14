package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBreakFileDownloadPathToParts(t *testing.T) {
	testCases := []struct {
		name         string
		downloadPath string
		expectedRepo string
		expectedPath string
		expectedName string
		expectError  bool
	}{
		{"Single level path", "repo/file.txt", "repo", "", "file.txt", false},
		{"Multi-level path", "repo/folder/subfolder/file.txt", "repo", "folder/subfolder", "file.txt", false},
		{"Root level file", "repo/", "", "", "", true},
		{"Empty path", "", "", "", "", true},
		{"Invalid path", "file.txt", "", "", "", true},
		{"Wildcard path", "repo/*.txt", "", "", "", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repo, path, name, err := breakFileDownloadPathToParts(tt.downloadPath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedRepo, repo)
			assert.Equal(t, tt.expectedPath, path)
			assert.Equal(t, tt.expectedName, name)
		})
	}
}

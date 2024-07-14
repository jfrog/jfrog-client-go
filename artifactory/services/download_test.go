package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBreakFileDownloadPathToParts(t *testing.T) {
	tests := []struct {
		name         string
		downloadPath string
		expectedRepo string
		expectedPath string
		expectedName string
		expectError  bool
	}{
		{
			name:         "Single level path",
			downloadPath: "repo/file.txt",
			expectedRepo: "repo",
			expectedPath: "",
			expectedName: "file.txt",
			expectError:  false,
		},
		{
			name:         "Multi-level path",
			downloadPath: "repo/folder/subfolder/file.txt",
			expectedRepo: "repo",
			expectedPath: "folder/subfolder",
			expectedName: "file.txt",
			expectError:  false,
		},
		{
			name:         "Root level file",
			downloadPath: "repo/",
			expectedRepo: "repo",
			expectedPath: "",
			expectedName: "",
			expectError:  false,
		},
		{
			name:         "Empty path",
			downloadPath: "",
			expectedRepo: "",
			expectedPath: "",
			expectedName: "",
			expectError:  true,
		},
		{
			name:         "Invalid path",
			downloadPath: "file.txt",
			expectedRepo: "",
			expectedPath: "",
			expectedName: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
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

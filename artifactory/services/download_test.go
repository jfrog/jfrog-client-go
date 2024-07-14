package services

import "testing"

func BreakFileDownloadPathToParts_WithSingleLevelPath_ReturnsCorrectParts(t *testing.T) {
	downloadPath := "repo/file.txt"
	expectedRepo, expectedPath, expectedName := "repo", "", "file.txt"
	repo, path, name, err := breakFileDownloadPathToParts(downloadPath)
	if err != nil || repo != expectedRepo || path != expectedPath || name != expectedName {
		t.Errorf("Expected (%s, %s, %s), got (%s, %s, %s, %v)", expectedRepo, expectedPath, expectedName, repo, path, name, err)
	}
}

func BreakFileDownloadPathToParts_WithMultiLevelPath_ReturnsCorrectParts(t *testing.T) {
	downloadPath := "repo/folder/subfolder/file.txt"
	expectedRepo, expectedPath, expectedName := "repo", "folder/subfolder", "file.txt"
	repo, path, name, err := breakFileDownloadPathToParts(downloadPath)
	if err != nil || repo != expectedRepo || path != expectedPath || name != expectedName {
		t.Errorf("Expected (%s, %s, %s), got (%s, %s, %s, %v)", expectedRepo, expectedPath, expectedName, repo, path, name, err)
	}
}

func BreakFileDownloadPathToParts_WithRootLevelFile_ReturnsEmptyPathAndName(t *testing.T) {
	downloadPath := "repo/"
	expectedRepo, expectedPath, expectedName := "repo", "", ""
	repo, path, name, err := breakFileDownloadPathToParts(downloadPath)
	if err != nil || repo != expectedRepo || path != expectedPath || name != expectedName {
		t.Errorf("Expected (%s, %s, %s), got (%s, %s, %s, %v)", expectedRepo, expectedPath, expectedName, repo, path, name, err)
	}
}

func BreakFileDownloadPathToParts_WithEmptyPath_ReturnsError(t *testing.T) {
	downloadPath := ""
	_, _, _, err := breakFileDownloadPathToParts(downloadPath)
	if err == nil {
		t.Errorf("Expected error for empty download path, got nil")
	}
}

func BreakFileDownloadPathToParts_WithInvalidPath_ReturnsError(t *testing.T) {
	downloadPath := "file.txt"
	_, _, _, err := breakFileDownloadPathToParts(downloadPath)
	if err == nil {
		t.Errorf("Expected error for invalid download path, got nil")
	}
}

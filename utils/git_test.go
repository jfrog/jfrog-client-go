package utils

import (
	"testing"
)

func TestGitManager(t *testing.T) {
	projectPath := initVcsTestDir(t)
	gitManager := NewGitManager(projectPath)
	err := gitManager.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	want := "https://github.com/jfrog/jfrog-cli.git"
	if gitManager.GetUrl() != want {
		t.Errorf("TestGitManager() error, want %s, got %s", want, gitManager.GetUrl())
	}
	want = "d63c5957ad6819f4c02a817abe757f210d35ff92"
	if gitManager.GetRevision() != want {
		t.Errorf("TestGitManager() error, want %s, got %s", want, gitManager.GetRevision())
	}
}

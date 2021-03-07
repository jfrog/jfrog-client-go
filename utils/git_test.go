package utils

import (
	"bytes"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type gitManager struct {
	dotGitPath string
}

func GitExecutor(dotGitPath string) *gitManager {
	return &gitManager{dotGitPath: dotGitPath}
}

func (m *gitManager) GetUrl() (string, string, error) {
	return m.execGit("config", "--get", "remote.origin.url")
}

func (m *gitManager) GetRevision() (string, string, error) {
	return m.execGit("show", "-s", "--format=%H", "HEAD")
}

func (m *gitManager) GetBranch() (string, string, error) {
	return m.execGit("branch", "--show-current")
}

func (m *gitManager) GetMessage(revision string) (string, string, error) {
	return m.execGit("show", "-s", "--format=%B", revision)
}

func (m *gitManager) execGit(args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Dir = m.dotGitPath
	cmd.Stdin = nil
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	errorutils.CheckError(err)
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func TestReadConfig(t *testing.T) {
	dotGitPath := getDotGitPath(t)
	gitManager := NewGitManager(dotGitPath)
	err := gitManager.ReadConfig()

	gitExecutor := GitExecutor(dotGitPath)
	revision, _, err := gitExecutor.GetRevision()
	assert.NoError(t, err)
	assert.Equal(t, revision, gitManager.GetRevision(), "Wrong revision")
	branch, _, err := gitExecutor.GetBranch()
	assert.NoError(t, err)
	assert.Equal(t, branch, gitManager.GetBranch(), "Wrong branch")
	message, _, err := gitExecutor.GetMessage(revision)
	assert.NoError(t, err)
	assert.Equal(t, message, gitManager.GetMessage(), "Wrong message")

}

func getDotGitPath(t *testing.T) string {
	dotGitPath, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current dir.")
	dotGitPath = filepath.Dir(dotGitPath)
	dotGitExists, err := fileutils.IsDirExists(filepath.Join(dotGitPath, ".git"), false)
	assert.NoError(t, err)
	assert.True(t, dotGitExists, "Can't find .git")
	return dotGitPath
}

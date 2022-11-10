package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type gitExecutor struct {
	gitManager GitManager
}

func NewGitExecutor(dotGitPath string) *gitExecutor {
	return &gitExecutor{gitManager: *NewGitManager(dotGitPath)}
}

func (m *gitExecutor) execGit(args ...string) (string, string, error) {
	return m.gitManager.ExecGit(args...)
}
func (m *gitExecutor) GetUrl() (string, string, error) {
	return m.execGit("config", "--get", "remote.origin.url")
}

func (m *gitExecutor) GetRevision() (string, string, error) {
	return m.execGit("show", "-s", "--format=%H", "HEAD")
}

func (m *gitExecutor) GetBranch() (string, string, error) {
	return m.execGit("branch", "--show-current")
}

func (m *gitExecutor) GetMessage(revision string) (string, string, error) {
	return m.execGit("show", "-s", "--format=%B", revision)
}

func TestReadConfig(t *testing.T) {
	testReadConfig(t)
}

// Open a git repo using 'go-git' package fails when:
//  1. OS is Windows.
//  2. using go-git v4.7.0.
//  3. .git/config file contains path with backslashes.
func TestReadConfigWithBackslashes(t *testing.T) {
	dotGitPath := getDotGitPath(t)
	gitExec := NewGitExecutor(dotGitPath)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	_, _, err := gitExec.execGit("config", "--local", "--add", "http.https://github.com.sslCAInfo"+timestamp, dotGitPath)
	assert.NoError(t, err)
	defer func() {
		_, _, err = gitExec.execGit("config", "--local", "--unset", "http.https://github.com.sslCAInfo"+timestamp)
		assert.NoError(t, err)
	}()
	testReadConfig(t)
}

func testReadConfig(t *testing.T) {
	dotGitPath := getDotGitPath(t)
	gitManager := NewGitManager(dotGitPath)
	err := gitManager.ReadConfig()
	assert.NoError(t, err)

	gitExecutor := NewGitExecutor(dotGitPath)
	url, _, err := gitExecutor.GetUrl()
	assert.NoError(t, err)
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}
	assert.Equal(t, url, gitManager.GetUrl(), "Wrong url")
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
	assert.NoError(t, err, "Failed to get current dir")
	dotGitPath = filepath.Dir(dotGitPath)
	dotGitExists, err := fileutils.IsDirExists(filepath.Join(dotGitPath, ".git"), false)
	assert.NoError(t, err)
	assert.True(t, dotGitExists, "Can't find .git")
	return dotGitPath
}

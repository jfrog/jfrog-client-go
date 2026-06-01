package utils

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
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

// TestReadConfigEmptyGitDir tests that ReadConfig handles an empty .git directory gracefully
// instead of crashing with "no such file or directory" error when .git/HEAD is missing.
// This simulates a corrupt or uninitialized git repository scenario common in CI environments.
func TestReadConfigEmptyGitDir(t *testing.T) {
	// Create a temp directory with an empty .git folder (no HEAD file)
	tempDir, err := os.MkdirTemp("", "test-empty-git")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	// Create empty .git directory (simulating corrupt/uninitialized repo)
	emptyGitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(emptyGitDir, 0755)
	assert.NoError(t, err)

	// This should NOT crash - it should handle the missing HEAD file gracefully
	gitManager := NewGitManager(tempDir)
	err = gitManager.ReadConfig()

	// The function should not return an error for missing HEAD/config files
	// Instead, it should gracefully handle this case and return empty values
	assert.NoError(t, err, "ReadConfig should handle empty .git directory gracefully")

	// Values should be empty but not cause a crash
	assert.Empty(t, gitManager.GetRevision())
	assert.Empty(t, gitManager.GetUrl())
	assert.Empty(t, gitManager.GetBranch())
}

// TestReadConfigMissingHeadFile tests that ReadConfig handles missing HEAD file gracefully
func TestReadConfigMissingHeadFile(t *testing.T) {
	// Create a temp directory with .git folder containing only config file
	tempDir, err := os.MkdirTemp("", "test-missing-head")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	// Create .git directory with config but no HEAD
	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	assert.NoError(t, err)

	// Create a minimal config file
	configContent := `[remote "origin"]
	url = https://github.com/test/repo.git
`
	err = os.WriteFile(filepath.Join(gitDir, "config"), []byte(configContent), 0644)
	assert.NoError(t, err)

	gitManager := NewGitManager(tempDir)
	err = gitManager.ReadConfig()

	// Should not error - gracefully handle missing HEAD
	assert.NoError(t, err)
	// URL should be read from config
	assert.Equal(t, "https://github.com/test/repo.git", gitManager.GetUrl())
	// Revision and branch should be empty (no HEAD file)
	assert.Empty(t, gitManager.GetRevision())
	assert.Empty(t, gitManager.GetBranch())
}

// TestReadConfigMissingConfigFile tests that ReadConfig handles missing config file gracefully
func TestReadConfigMissingConfigFile(t *testing.T) {
	// Create a temp directory with .git folder containing only HEAD file
	tempDir, err := os.MkdirTemp("", "test-missing-config")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()

	// Create .git directory with HEAD but no config
	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	assert.NoError(t, err)

	// Create HEAD file pointing to a branch
	headContent := "ref: refs/heads/main\n"
	err = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(headContent), 0644)
	assert.NoError(t, err)

	gitManager := NewGitManager(tempDir)
	err = gitManager.ReadConfig()

	// Should not error - gracefully handle missing config
	assert.NoError(t, err)
	// URL should be empty (no config file)
	assert.Empty(t, gitManager.GetUrl())
	// Branch should be read from HEAD
	assert.Equal(t, "main", gitManager.GetBranch())
}

package services

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"github.com/jfrog/jfrog-client-go/utils/tests"
)

func TestExtractRepo(t *testing.T) {
	pwd, err := os.Getwd()
	testPath := filepath.Join(pwd, "testsdata", "gitlfs")
	repo, err := extractRepo(testPath, "lfsConfigExample", "https://localhost:8080/artifactory", lfsConfigUrlExtractor)
	if err != nil {
		t.Error("Got err: ", err)
	}
	if repo != "lfs-local" {
		t.Error("Failed to extract repo from .lfsconfig file format. Expected: \"lfs-local\" Got: ", repo)

	}
	repo, err = extractRepo(testPath, "configExample", "http://localhost:8081/artifactory", configLfsUrlExtractor)
	if err != nil {
		t.Error("Got err: ", err)
	}
	if repo != "lfs-local" {
		t.Error("Failed to extract repo from .git/config file format. Expected: \"lfs-local\" Got: ", repo)
	}
}

func TestGetLfsFilesFromGit(t *testing.T) {
	fileId := "4bf4c8c0fef3f5c8cf6f255d1c784377138588c0a9abe57e440bce3ccb350c2e"
	gitPath := getCliDotGitPath(t)
	refs := strings.Join([]string{"refs", "heads", "*"}, "/")
	if runtime.GOOS == "windows" {
		refs = strings.Join([]string{"refs", "heads", "*"}, "\\\\")
	}
	results, err := getLfsFilesFromGit(gitPath, refs)
	if err != nil {
		t.Error("Got err: ", err)
	}
	_, ok := results[fileId]
	if !ok {
		t.Error("couldn't find test.bin test file")
	}
}

func getCliDotGitPath(t *testing.T) string {
	dotGitPath, err := tests.FindGitRoot()
	if err != nil {
		t.Error("Failed to get current dir.")
	}
	return dotGitPath
}

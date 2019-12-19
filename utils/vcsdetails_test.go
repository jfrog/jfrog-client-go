package utils

import (
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

func TestVcsDetails(t *testing.T) {
	path := initVcsTestDir(t)
	vcsDetals := NewVcsDetals()
	revision, url, err := vcsDetals.GetVcsDetails(filepath.Join(path))
	if err != nil {
		t.Error(err)
	}
	if url != "https://github.com/jfrog/jfrog-cli.git" {
		t.Errorf("TestGitManager() error, want %s, got %s", url, "https://github.com/jfrog/jfrog-cli.git")
	}
	if revision != "d63c5957ad6819f4c02a817abe757f210d35ff92" {
		t.Errorf("TestGitManager() error, want %s, got %s", url, "d63c5957ad6819f4c02a817abe757f210d35ff92")
	}
}

func initVcsTestDir(t *testing.T) string {
	testsdataSrc := filepath.Join("testsdata", "vcs")
	testsdataTarget := filepath.Join("testsdata", "tmp")
	err := fileutils.CopyDir(testsdataSrc, testsdataTarget, true)
	if err != nil {
		t.Error(err)
	}
	if found, err := fileutils.IsDirExists(filepath.Join(testsdataTarget, "gitdata"), false); found {
		if err != nil {
			t.Error(err)
		}
		err := fileutils.RenamePath(filepath.Join(testsdataTarget, "gitdata"), filepath.Join(testsdataTarget, ".git"))
		if err != nil {
			t.Error(err)
		}
	}
	if found, err := fileutils.IsDirExists(filepath.Join(testsdataTarget, "OtherGit", "gitdata"), false); found {
		if err != nil {
			t.Error(err)
		}
		err := fileutils.RenamePath(filepath.Join(testsdataTarget, "OtherGit", "gitdata"), filepath.Join(testsdataTarget, "OtherGit", ".git"))
		if err != nil {
			t.Error(err)
		}
	}
	path, err := filepath.Abs(testsdataTarget)
	if err != nil {
		t.Error(err)
	}
	return path
}

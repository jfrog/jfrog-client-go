package fileutils

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestIsSsh(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://some.url", false},
		{"https://some.url", false},
		{"sshd://wrong.url", false},
		{"assh://wrong.url", false},
		{"ssh://some.url", true},
		{"sSh://some.url/some/api", true},
		{"SSH://some.url/some/api", true},
	}
	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			if IsSshUrl(test.url) != test.expected {
				t.Error("Expected '"+strconv.FormatBool(test.expected)+"' Got: '"+strconv.FormatBool(!test.expected)+"' For URL:", test.url)
			}
		})
	}
}

func TestGetFileOrDirPathFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	defer os.Chdir(wd)

	// CD into a directory with a go.mod file.
	projectRoot := filepath.Join("testdata", "project")
	err = os.Chdir(projectRoot)
	if err != nil {
		t.Error(err)
	}

	// Make projectRoot an absolute path.
	projectRoot, err = os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// Get the project root.
	root, exists, err := FindUpstream("go.mod", File)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("File go.mod is missing.")
	}

	if root != projectRoot {
		t.Error("Expecting", projectRoot, "got:", root)
	}

	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		t.Error(err)
	}

	// CD into a sub directory in the same project, and expect to get the same project root.
	os.Chdir(wd)
	projectSubDirectory := filepath.Join("testdata", "project", "dir")
	err = os.Chdir(projectSubDirectory)
	if err != nil {
		t.Error(err)
	}
	root, exists, err = FindUpstream("go.mod", File)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("File go.mod is missing.")
	}
	if root != projectRoot {
		t.Error("Expecting", projectRoot, "got:", root)
	}

	root, exists, err = FindUpstream("go-missing.mod", File)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error("File go-missing.mod found but shouldn't.")
	}

	if root != "" {
		t.Error("File go-missing.mod shouldn't be found, however, got:", root)
	}

	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		t.Error(err)
	}

	// Now CD into a directory outside the project, and expect to get a different project root.
	noProjectRoot := filepath.Join("testdata", "noproject")
	err = os.Chdir(noProjectRoot)
	if err != nil {
		t.Error(err)
	}
	root, exists, err = FindUpstream("go.mod", File)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("File go.mod is missing.")
	}
	if root == projectRoot {
		t.Error("Expecting a different value than", root)
	}
}

func TestGetFileOrDirPathFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	defer os.Chdir(wd)

	// Create path to directory to find.
	dirPath := filepath.Join("testdata")
	err = os.Chdir(dirPath)
	if err != nil {
		t.Error(err)
	}
	// Get absolute path.
	dirPath, err = os.Getwd()
	if err != nil {
		t.Error(err)
	}
	// CD back to the original directory.
	if err := os.Chdir(wd); err != nil {
		t.Error(err)
	}

	// Go to starting dir to search from.
	searchFromDir := filepath.Join("testdata", "project", "dir")
	err = os.Chdir(searchFromDir)
	if err != nil {
		t.Error(err)
	}

	// Get the directory path.
	root, exists, err := FindUpstream("noproject", Dir)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("Dir noproject is missing.")
	}
	if root != dirPath {
		t.Error("Expecting", dirPath, "got:", root)
	}
}

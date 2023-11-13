package fspatterns

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterFiles(t *testing.T) {
	data := []struct {
		files          []string
		ExcludePattern string
		root           string
		result         []string
	}{
		{[]string{"file1", filepath.Join("dir", "file1"), "file2.zip"}, "^*.zip$", "", []string{"file1", filepath.Join("dir", "file1")}},
		{[]string{
			"file1",
			"test.zip",
			filepath.Join("test", "file1"),
			filepath.Join("dir", "test", "should-be-filter"),
		}, "(^.*test.*$)", "test", []string{"file1", "test.zip", filepath.Join("test", "file1")}},
	}
	for _, d := range data {
		got, err := filterFiles(d.root, d.files, d.ExcludePattern)
		assert.NoError(t, err)
		assert.Len(t, got, len(d.result))
		assert.Contains(t, got, d.files[0])
		assert.Contains(t, got, d.files[1])
	}
}

func TestSearchPatterns(t *testing.T) {
	data := []struct {
		path    string
		pattern string
		result  []string
	}{
		{filepath.Join("testdata", "a", "a3.zip"), "^*.zip$", []string{filepath.Join("testdata", "a", "a3")}},
		{filepath.Join("testdata", "a", "a3"), "^*.zip$", []string{}},
	}
	for _, d := range data {
		patternRegex, err := regexp.Compile(d.pattern)
		assert.NoError(t, err)

		matches, isDir, err := SearchPatterns(d.path, true, true, patternRegex)
		assert.NoError(t, err)
		assert.False(t, isDir)
		assert.Len(t, matches, len(d.result))
	}
}

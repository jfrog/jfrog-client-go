package fspatterns

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterFilesFunc(t *testing.T) {
	testCases := []struct {
		file           string
		ExcludePattern string
		root           string
		included       bool
	}{
		// Patterns with regex
		{"file1", "^*.zip$", "", true},
		{"file2.zip", "^*.zip$", "", false},
		{"dir/file1", "^*.zip$", "", true},
		{"dir/dir2/file1.zip", "^*.zip$", "", false},

		{"test/file1", "(^.*test.*$)", "test", true},
		{"dir/test/should-be-filter", "(^.*test.*$)", "test", false},
		{"file1", "(^.*test.*$)", "", true},
		{"file2.zip", "(^.*test.*$)", "", true},

		// Patterns without regex (exact match)
		{"file1", "file1", "", false},
		{"file2.zip", "file1", "", true},
		// No pattern
		{"file1", "", "", true},
		{"file2.zip", "", "", true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("File: %s, Pattern: %s, Root: %s", tc.file, tc.ExcludePattern, tc.root), func(t *testing.T) {
			// Create the filter function with the mocked isPathExcluded
			filterFunc := filterFilesFunc(tc.root, tc.ExcludePattern, nil)
			excluded, err := filterFunc(tc.file)
			assert.NoError(t, err)
			assert.True(t, excluded == tc.included, "Expected included = %v, but got %v", tc.included, excluded)
		})
	}
}

func TestSearchPatterns(t *testing.T) {
	testCases := []struct {
		path    string
		pattern string
		result  []string
	}{
		{filepath.Join("testdata", "a", "a3.zip"), "^*.zip$", []string{filepath.Join("testdata", "a", "a3")}},
		{filepath.Join("testdata", "a", "a3"), "^*.zip$", []string{}},
	}
	for _, d := range testCases {
		patternRegex, err := regexp.Compile(d.pattern)
		assert.NoError(t, err)

		matches, isDir, err := SearchPatterns(d.path, true, true, patternRegex)
		assert.NoError(t, err)
		assert.False(t, isDir)
		assert.Len(t, matches, len(d.result))
	}
}

func TestFilterFilesFuncWithSizeThreshold(t *testing.T) {
	rootPath := t.TempDir()

	// Create test files and directories
	files := []struct {
		path string
		size int64
	}{
		{filepath.Join(rootPath, "file.txt"), 100},
		{filepath.Join(rootPath, "largefile.txt"), 2048},
		{filepath.Join(rootPath, "dir", "subfile.txt"), 50},
		{filepath.Join(rootPath, "equalfile.txt"), 1024},
	}

	for _, file := range files {
		dir := filepath.Dir(file.path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			assert.NoError(t, os.MkdirAll(dir, 0755))
		}
		f, err := os.Create(file.path)
		assert.NoError(t, err)
		assert.NoError(t, f.Truncate(file.size))
		assert.NoError(t, f.Close())
	}

	// Test cases
	testCases := []struct {
		name          string
		path          string
		sizeThreshold *SizeThreshold
		expectInclude bool
	}{
		{"Include file within size threshold", "file.txt", &SizeThreshold{Size: 1024, Condition: LessThan}, true},
		{"Exclude file exceeding size threshold", "largefile.txt", &SizeThreshold{Size: 1024, Condition: LessThan}, false},
		{"Include directory", "dir", nil, true},
		{"Include file in subdirectory within size threshold", filepath.Join("dir", "subfile.txt"), &SizeThreshold{Size: 1024, Condition: LessThan}, true},
		{"Include file with size equal to threshold", "equalfile.txt", &SizeThreshold{Size: 1024, Condition: GreaterEqualThan}, true},
		{"Exclude file below threshold with GreaterEqualThan", "file.txt", &SizeThreshold{Size: 150, Condition: GreaterEqualThan}, false},
		{"Include file above threshold with GreaterEqualThan", "largefile.txt", &SizeThreshold{Size: 150, Condition: GreaterEqualThan}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filterFunc := filterFilesFunc(rootPath, "", tc.sizeThreshold)
			included, err := filterFunc(filepath.Join(rootPath, tc.path))
			assert.Equal(t, tc.expectInclude, included)
			assert.NoError(t, err)
		})
	}
}

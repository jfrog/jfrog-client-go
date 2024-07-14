package fspatterns

import (
	"fmt"
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

package fileutils

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCleanOldDirs(t *testing.T) {
	defer func(originTempPrefix string) {
		tempPrefix = originTempPrefix
	}(tempPrefix)
	tempPrefix = "test." + tempPrefix
	tempDir, err := CreateTempDir()
	assert.NoError(t, err)
	tempFile, err := CreateTempFile()
	assert.NoError(t, tempFile.Close())
	assert.NoError(t, err)

	// Check file exists.
	AssertFileExists(t, tempDir)
	AssertFileExists(t, tempFile.Name())

	// Don't delete valid files.
	assert.NoError(t, CleanOldDirs())
	AssertFileExists(t, tempDir)
	AssertFileExists(t, tempFile.Name())

	// Delete expired files.
	oldMaxFileAge := maxFileAge
	maxFileAge = 0
	defer func() { maxFileAge = oldMaxFileAge }()
	assert.NoError(t, CleanOldDirs())

	// Check if the file got deleted.
	_, err1 := os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err1))
	_, err2 := os.Stat(tempFile.Name()) // #nosec G703 -- test file; path from temp file
	assert.True(t, os.IsNotExist(err2))
}

func TestExtractTimestamp(t *testing.T) {
	testCases := []struct {
		item         string
		expectedTime time.Time
		expectError  bool
	}{
		// Valid cases
		{"jfrog.cli.temp.prefix-1625097600-suffix", time.Unix(1625097600, 0), false},
		{"jfrog.cli.temp.some-1234567890-other", time.Unix(1234567890, 0), false},

		// Invalid cases
		{"jfrog.cli.temp.no-dash", time.Time{}, true},
		{"jfrog.cli.temp.one-dash-1234567890", time.Time{}, true},
		{"jfrog.cli.temp.two-dashes--", time.Time{}, true},
		{"jfrog.cli.temp.prefix--suffix", time.Time{}, true},
		{"jfrog.cli.temp.prefix-abc-suffix", time.Time{}, true},
		{"jfrog.cli.temp.prefix-1625097600suffix", time.Time{}, true},
	}

	for _, test := range testCases {
		t.Run(test.item, func(t *testing.T) {
			result, err := extractTimestamp(test.item)
			if (err != nil) != test.expectError {
				t.Errorf("expected error: %v, got: %v", test.expectError, err)
			}
			if !result.Equal(test.expectedTime) {
				t.Errorf("expected time: %v, got: %v", test.expectedTime, result)
			}
		})
	}
}

func AssertFileExists(t *testing.T, name string) {
	_, err := os.Stat(name) // #nosec G703 -- test helper; path from test
	assert.NoError(t, err)
}

// TestCleanOldDirsContinuesOnError tests that cleanup continues even when encountering errors.
// This test verifies that CleanOldDirs() processes all files and collects errors instead of stopping at first error.
func TestCleanOldDirsContinuesOnError(t *testing.T) {
	// Save original values
	defer func(originTempPrefix string) {
		tempPrefix = originTempPrefix
	}(tempPrefix)
	tempPrefix = "test.continue." + tempPrefix

	oldTempDirBase := tempDirBase
	defer func() { tempDirBase = oldTempDirBase }()

	// Create a temporary test directory
	testDir, err := os.MkdirTemp("", "test-cleanup-*")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(testDir))
	}()

	tempDirBase = testDir

	// Save original maxFileAge
	oldMaxFileAge := maxFileAge
	defer func() { maxFileAge = oldMaxFileAge }()
	maxFileAge = 0 // Make all files appear old

	// Create valid temp files that should be deleted
	validFile1, err := CreateTempFile()
	assert.NoError(t, err)
	validFile1Name := validFile1.Name()
	assert.NoError(t, validFile1.Close())

	// Create a file with invalid timestamp format (will cause extractTimestamp error)
	invalidFile := path.Join(testDir, tempPrefix+"invalid-no-timestamp")
	err = os.WriteFile(invalidFile, []byte("test"), 0644)
	assert.NoError(t, err)

	// Create another valid file that should be deleted
	validFile2, err := CreateTempFile()
	assert.NoError(t, err)
	validFile2Name := validFile2.Name()
	assert.NoError(t, validFile2.Close())

	// Verify all files exist before cleanup
	AssertFileExists(t, validFile1Name)
	AssertFileExists(t, invalidFile)
	AssertFileExists(t, validFile2Name)

	// Run cleanup
	err = CleanOldDirs()

	// After fix: should return error mentioning the invalid file
	// but continue processing other files
	assert.Error(t, err, "Should return error for invalid file")
	assert.Contains(t, err.Error(), "failed to cleanup")
	assert.Contains(t, err.Error(), "invalid-no-timestamp")

	// Verify valid files were deleted despite error with invalid file
	_, err1 := os.Stat(validFile1Name) // #nosec G703 -- test file; path from test temp dir
	assert.True(t, os.IsNotExist(err1), "validFile1 should be deleted")

	_, err2 := os.Stat(validFile2Name) // #nosec G703 -- test file; path from test temp dir
	assert.True(t, os.IsNotExist(err2), "validFile2 should be deleted")

	// Invalid file should still exist (couldn't be processed)
	AssertFileExists(t, invalidFile)

	// Cleanup
	assert.NoError(t, os.Remove(invalidFile))
}

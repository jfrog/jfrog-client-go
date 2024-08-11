package fileutils

import (
	"os"
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
	_, err2 := os.Stat(tempFile.Name())
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
	_, err := os.Stat(name)
	assert.NoError(t, err)
}

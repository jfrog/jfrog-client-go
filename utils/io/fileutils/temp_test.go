package fileutils

import (
	"os"
	"testing"

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
	// Extract time from a file.
	fileName := "jfrog.cli.temp.-008652489-1595147819.json"
	timeStamp, err := extractTimestamp(fileName)
	assert.NoError(t, err)
	assert.Equal(t, int64(8652489), timeStamp.Unix())

	// Extract time from a dir.
	fileName = "asd-asjfrog.cli.temp.-008652489-1595147444"
	timeStamp, err = extractTimestamp(fileName)
	assert.NoError(t, err)
	assert.Equal(t, int64(8652489), timeStamp.Unix())
}

func AssertFileExists(t *testing.T, name string) {
	_, err := os.Stat(name)
	assert.NoError(t, err)
}

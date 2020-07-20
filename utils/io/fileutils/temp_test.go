package fileutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanOldDirs(t *testing.T) {
	tempDir, err := CreateTempDir()
	assert.NoError(t, err)
	tempFile, err := CreateTempFile()
	assert.NoError(t, err)

	// Check file exists.
	_, err = os.Stat(tempDir)
	assert.NoError(t, err)
	_, err = os.Stat(tempFile.Name())
	assert.NoError(t, err)

	// Don't delete valid files.
	assert.NoError(t, CleanOldDirs())
	_, err = os.Stat(tempDir)
	assert.NoError(t, err)
	_, err = os.Stat(tempFile.Name())
	assert.NoError(t, err)

	// Delete expired files.
	maxFileAge = 0
	assert.NoError(t, CleanOldDirs())

	// Check if the file got deleted.
	_, err1 := os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err1))
	_, err2 := os.Stat(tempFile.Name())
	assert.True(t, os.IsNotExist(err2))
}

func TestExtractTimestamp(t *testing.T) {
	// Extract time from a file.
	fileName := "jfrog.cli.temp.008652489-1595147819.json"
	timeStamp, err := extractTimestamp(fileName)
	assert.NoError(t, err)
	assert.Equal(t, timeStamp.Unix(), int64(1595147819))

	// Extract time from a dir.
	fileName = "jfrog.cli.temp.008652489-1595147444"
	timeStamp, err = extractTimestamp(fileName)
	assert.NoError(t, err)
	assert.Equal(t, timeStamp.Unix(), int64(1595147444))
}

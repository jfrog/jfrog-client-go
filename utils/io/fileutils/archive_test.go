package fileutils

import (
	clientTestUtils "github.com/jfrog/jfrog-client-go/utils/tests"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnarchive(t *testing.T) {
	tests := []string{"zip", "tar", "tar.gz"}
	for _, extension := range tests {
		t.Run(extension, func(t *testing.T) {
			// Create temp directory
			tmpDir, createTempDirCallback := clientTestUtils.CreateTempDirWithCallbackAndAssert(t)
			defer createTempDirCallback()
			// Run unarchive on archive created on Unix
			err := runUnarchive("unix."+extension, "archives", filepath.Join(tmpDir, "unix"))
			assert.NoError(t, err)
			assert.FileExists(t, filepath.Join(tmpDir, "unix", "link"))
			assert.FileExists(t, filepath.Join(tmpDir, "unix", "dir", "file"))

			// Run unarchive on archive created on Windows
			err = runUnarchive("win."+extension, "archives", filepath.Join(tmpDir, "win"))
			assert.NoError(t, err)
			assert.FileExists(t, filepath.Join(tmpDir, "win", "link.lnk"))
			assert.FileExists(t, filepath.Join(tmpDir, "win", "dir", "file.txt"))
		})
	}
}

func TestUnarchiveZipSlip(t *testing.T) {
	tests := []struct {
		testType    string
		archives    []string
		errorSuffix string
	}{
		{"rel", []string{"zip", "tar", "tar.gz"}, "illegal path in archive: '../file'"},
		{"abs", []string{"tar", "tar.gz"}, "illegal path in archive: '/tmp/bla/file'"},
		{"softlink-abs", []string{"zip", "tar", "tar.gz"}, "illegal link path in archive: '/tmp/bla/file'"},
		{"softlink-rel", []string{"zip", "tar", "tar.gz"}, "illegal link path in archive: '../file'"},
	}
	for _, test := range tests {
		t.Run(test.testType, func(t *testing.T) {
			// Create temp directory
			tmpDir, createTempDirCallback := clientTestUtils.CreateTempDirWithCallbackAndAssert(t)
			defer createTempDirCallback()
			for _, archive := range test.archives {
				// Unarchive and make sure an error returns
				err := runUnarchive(test.testType+"."+archive, "zipslip", tmpDir)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.errorSuffix)
			}
		})
	}
}

func runUnarchive(archiveFileName, sourceDir, targetDir string) error {
	return Unarchive(filepath.Join("testdata", sourceDir, archiveFileName), archiveFileName, targetDir)
}

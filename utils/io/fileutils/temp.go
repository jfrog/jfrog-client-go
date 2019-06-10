package fileutils

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io/ioutil"
	"os"
)

var tempDirBase string

func init() {
	tempDirBase = os.TempDir()
}

// Creates the temp dir at tempDirBase.
// Set tempDirPath to the created directory path.
func CreateTempDir() (string, error) {
	if tempDirBase == "" {
		return "", errorutils.CheckError(errors.New("Temp dir cannot be created in an empty base dir."))
	}

	path, err := ioutil.TempDir(tempDirBase, "jfrog.cli.")
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	return path, nil
}

// Change the containing directory of temp dir.
func SetTempDirBase(dirPath string) {
	tempDirBase = dirPath
}

func RemoveTempDir(dirPath string) error {
	exists, err := IsDirExists(dirPath, false)
	if err != nil {
		return err
	}
	if exists {
		return os.RemoveAll(dirPath)
	}
	return nil
}

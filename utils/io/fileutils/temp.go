package fileutils

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io/ioutil"
	"os"
)

var tempDirPath string
var tempDirBase string

func init() {
	tempDirBase = os.TempDir()
}

// Return the path of the existing temp dir.
// If not exist, return an error.
func GetTempDirPath() (dirPath string, err error) {
	if tempDirPath == "" {
		return "", errorutils.CheckError(errors.New("Function cannot be used before 'tempDirPath' is created."))
	}

	return tempDirPath, err
}

// Creates the temp dir at tempDirBase.
// Set tempDirPath to the created directory path.
func CreateTempDirPath() error {
	if tempDirBase == "" {
		return errorutils.CheckError(errors.New("Temp dir cannot be created in an empty base dir."))
	}
	if tempDirPath != "" {
		return errorutils.CheckError(errors.New("'tempDirPath' has already been initialized."))
	}

	path, err := ioutil.TempDir(tempDirBase, "jfrog.cli.")
	if err != nil {
		return errorutils.CheckError(err)
	}

	tempDirPath = path
	return nil
}

// Change the containing directory of temp dir.
func SetTempDirBase(dirPath string) error {
	if tempDirPath != "" {
		return errorutils.CheckError(errors.New("Cannot set temp base path after the temp dir has already been initialized."))
	}
	tempDirBase = dirPath
	return nil
}

func RemoveTempDir() error {
	defer func() {
		tempDirPath = ""
	}()

	exists, err := IsDirExists(tempDirPath, false)
	if err != nil {
		return err
	}
	if exists {
		return os.RemoveAll(tempDirPath)
	}
	return nil
}

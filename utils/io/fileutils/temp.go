package fileutils

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	tempPrefix = "jfrog.temp."
	deadline   = 24.0
)

//Path to the root temp dir
var tempDirBase string

//Path to the current flow temp dir
var tempDirReaderWriter string

func init() {
	tempDirBase, tempDirReaderWriter = os.TempDir(), os.TempDir()
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

// Create a temp dir named "tempPrefix+timeStamp".
func CreateReaderWriterTempDir() (err error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	tempDirReaderWriter, err = ioutil.TempDir(tempDirBase, tempPrefix+timeStamp)
	err = errorutils.CheckError(err)
	return
}

// Create a new temp file named "tempPrefix+timeStamp".
func CreateReaderWriterTempFile() (*os.File, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	if tempDirReaderWriter == "" {
		return nil, errorutils.CheckError(errors.New("Temp folder was not created"))
	}
	fd, err := ioutil.TempFile(tempDirReaderWriter, tempPrefix+timeStamp)
	return fd, err
}

func CleanupReaderWriterTempFilesAndDirs() error {
	// Cleanup dirs
	exists, err := IsDirExists(tempDirReaderWriter, false)
	if err != nil {
		return err
	}
	if exists {
		err = os.RemoveAll(tempDirReaderWriter)
		if err != nil {
			return err
		}
	}
	// Cleanup files
	exists, err = IsFileExists(tempDirReaderWriter, false)
	if exists {
		err = os.RemoveAll(tempDirReaderWriter)
		if err != nil {
			return err
		}
	}
	return cleanOldDirs()
}

// Old runs/tests may left junk at temp dir.
// Each temp file/Dir is named with prefix+timestamp, search for all temp files/dirs that match the common prefix and validate their timestamp.
func cleanOldDirs() error {
	// Get all files at temp dir
	files, err := ioutil.ReadDir(tempDirBase)
	if err != nil {
		log.Fatal(err)
	}

	// Search for files/dirs that match the template.
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), tempPrefix) {
			timeStamp, err := extractTimestamp(file.Name())
			if err != nil {
				return err
			}
			now := time.Now()
			// Delete old file/dirs.
			if now.Sub(timeStamp).Hours() > deadline {
				if err := os.RemoveAll(path.Join(tempDirBase, file.Name())); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func extractTimestamp(item string) (time.Time, error) {
	// Get timestamp from file/dir.
	timestampStr := strings.Replace(item, tempPrefix, "", 1)
	// Convert to int.
	timeStampint, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}, errorutils.CheckError(err)
	}
	// Convert to time type.
	return time.Unix(timeStampint, 0), nil
}

package fileutils

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

var (
	tempPrefix = "jfrog.cli.temp."

	// Max temp file age in hours
	maxFileAge = 24.0

	// Path to the root temp dir
	tempDirBase string
)

func init() {
	tempDirBase = os.TempDir()
}

// Creates the temp dir at tempDirBase.
// Set tempDirPath to the created directory path.
func CreateTempDir() (string, error) {
	if tempDirBase == "" {
		return "", errorutils.CheckErrorf("temp dir cannot be created in an empty base dir.")
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	dirPath, err := os.MkdirTemp(tempDirBase, tempPrefix+"-"+timestamp+"-")
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return dirPath, nil
}

// Change the containing directory of temp dir.
func SetTempDirBase(dirPath string) {
	tempDirBase = dirPath
}

func GetTempDirBase() string {
	return tempDirBase
}

func RemoveTempDir(dirPath string) error {
	exists, err := IsDirExists(dirPath, false)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	if err = os.RemoveAll(dirPath); err == nil {
		return nil
	}
	// Sometimes removing the directory fails (in Windows) because it's locked by another process.
	// That's a known issue, but its cause is unknown (golang.org/issue/30789).
	// In this case, we'll only remove the contents of the directory, and let CleanOldDirs() remove the directory itself at a later time.
	return RemoveDirContents(dirPath)
}

// Create a new temp file named "tempPrefix+timeStamp".
func CreateTempFile() (*os.File, error) {
	if tempDirBase == "" {
		return nil, errorutils.CheckErrorf("temp File cannot be created in an empty base dir.")
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	fd, err := os.CreateTemp(tempDirBase, tempPrefix+"-"+timestamp+"-")
	return fd, err
}

// Old runs/tests may leave junk at temp dir.
// Each temp file/Dir is named with prefix+timestamp, search for all temp files/dirs that match the common prefix and validate their timestamp.
func CleanOldDirs() error {
	// Get all files at temp dir
	files, err := os.ReadDir(tempDirBase)
	if err != nil {
		return errorutils.CheckError(err)
	}
	now := time.Now()
	// Search for files/dirs that match the template.
	for _, file := range files {
		fileName := file.Name()
		if strings.HasPrefix(fileName, tempPrefix) {
			var timestamp time.Time
			timestamp, err = extractTimestamp(fileName)
			if err != nil {
				return errorutils.CheckErrorf("could not extract timestamp from file %s: %q", fileName, err)
			}
			// Delete old file/dirs.
			if now.Sub(timestamp).Hours() > maxFileAge {
				if err = RemovePath(path.Join(tempDirBase, fileName)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func extractTimestamp(item string) (time.Time, error) {
	// Get the index of the last dash
	endTimestampIdx := strings.LastIndex(item, "-")
	if endTimestampIdx == -1 {
		return time.Time{}, fmt.Errorf("invalid format: no dash found")
	}

	// Get the index of the dash before the last dash
	beginningTimestampIdx := strings.LastIndex(item[:endTimestampIdx], "-")
	if beginningTimestampIdx == -1 {
		return time.Time{}, fmt.Errorf("invalid format: only one dash found")
	}

	// Extract the timestamp string
	timestampStr := item[beginningTimestampIdx+1 : endTimestampIdx]
	if len(timestampStr) == 0 {
		return time.Time{}, fmt.Errorf("invalid format: empty timestamp")
	}

	// Convert to int
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing timestamp: %v", err)
	}
	// Convert to time type
	return time.Unix(timestampInt, 0), nil
}

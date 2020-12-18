package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// localPath - The path of the downloaded archive file.
// localFileName - name of the archive file.
// originFileName - name of the archive file in Artifactory.
// logMsgPrefix - prefix log message.
// Extract an archive file to the 'localPath'.
func ExtractArchive(localPath, localFileName, originFileName, logMsgPrefix string) error {
	if !fileutils.IsSupportedArchive(originFileName) {
		return nil
	}
	extractionPath, err := getExtractionPath(localPath)
	if err != nil {
		return err
	}
	// In order to extract an archive file, the file extension is required, therefore,
	// we replace the local downloaded file name, with its origin in Artifactory.
	archivePath := filepath.Join(localPath, originFileName)
	if !strings.HasSuffix(localFileName, originFileName) {
		relativeLocalFilePath := localFileName
		if !strings.HasPrefix(relativeLocalFilePath, localPath) {
			relativeLocalFilePath = filepath.Join(localPath, relativeLocalFilePath)
		}
		err = fileutils.MoveFile(relativeLocalFilePath, archivePath)
		if err != nil {
			return err
		}
	}
	archivePath, err = filepath.Abs(archivePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(extractionPath, 0777)
	if errorutils.CheckError(err) != nil {
		return err
	}
	log.Info(logMsgPrefix+"Extracting archive:", archivePath, "to", extractionPath)
	return extract(archivePath, extractionPath)
}

func extract(localFilePath, extractionPath string) error {
	err := fileutils.Unarchive(localFilePath, extractionPath)
	if err != nil {
		return err
	}
	// If the file was extracted successfully, remove it from the file system
	return errorutils.CheckError(os.Remove(localFilePath))
}

func getExtractionPath(localPath string) (string, error) {
	// The local path to which the file is going to be extracted,
	// needs to be absolute.
	absolutePath, err := filepath.Abs(localPath)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	// Add a trailing slash to the local path, since it has to be a directory.
	return absolutePath + string(os.PathSeparator), nil
}

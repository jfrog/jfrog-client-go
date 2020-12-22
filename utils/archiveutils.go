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
	// TODO: localFileName may already include the local path.
	// We should investigate and make sure 'localFileName' contains file name only.
	var archivePath string
	if !strings.HasPrefix(localFileName, localPath) {
		archivePath = filepath.Join(localPath, localFileName)
	} else {
		archivePath = localFileName
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
	if err = extract(archivePath, originFileName, extractionPath); err != nil {
		log.Info(logMsgPrefix+"failed to extract", archivePath, ". error:", err.Error())
	}
	return err
}

func extract(localFilePath, originArchiveName, extractionPath string) error {
	err := fileutils.Unarchive(localFilePath, originArchiveName, extractionPath)
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

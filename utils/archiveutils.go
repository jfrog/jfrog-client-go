package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jfrog/gofrog/unarchive"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// localPath - The path of the downloaded archive file.
// localFileName - The name of the archive file.
// originFileName - The name of the archive file in Artifactory.
// logMsgPrefix - A prefix to the log message.
// bypassInspection - Set to true to bypass archive inspection against ZipSlip
// Extract an archive file to the 'localPath'.
func ExtractArchive(localPath, localFileName, originFileName, logMsgPrefix string, bypassInspection bool) error {
	unarchiver := &unarchive.Unarchiver{
		BypassInspection: bypassInspection,
	}
	if !unarchiver.IsSupportedArchive(originFileName) {
		return nil
	}
	extractionPath, err := getExtractionPath(localPath)
	if err != nil {
		return err
	}
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
	err = os.MkdirAll(extractionPath, 0755)
	if errorutils.CheckError(err) != nil {
		return err
	}
	log.Info(logMsgPrefix+"Extracting archive:", archivePath, "to", extractionPath)
	return errorutils.CheckError(extract(archivePath, originFileName, extractionPath, unarchiver))
}

func extract(localFilePath, originArchiveName, extractionPath string, unarchiver *unarchive.Unarchiver) error {
	if err := unarchiver.Unarchive(localFilePath, originArchiveName, extractionPath); err != nil {
		return errorutils.CheckError(err)
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

package fileutils

import (
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/mholt/archiver/v3"
)

func IsSupportedArchive(filePath string) bool {
	iArchiver, err := archiver.ByExtension(filePath)
	if err != nil {
		return false
	}
	_, ok := iArchiver.(archiver.Unarchiver)
	return ok
}

func Unarchive(archivePath, destinationPath string) error {
	tempDirPath, err := CreateTempDir()
	if err != nil {
		return err
	}
	defer RemoveTempDir(tempDirPath)

	err = archiver.Unarchive(archivePath, tempDirPath)
	if err != nil {
		return errorutils.CheckError(errors.New(fmt.Sprintf("Failed unarchiving: %s", archivePath) + err.Error()))
	}

	return MoveDir(tempDirPath, destinationPath)
}

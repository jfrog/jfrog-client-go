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

// The 'archiver' dependency includes an API called 'Unarchive' to extract archive files. This API uses the archive file
// extension to determine the archive type.// the local file path to extract the archive.
// We therefore need to use the file name as it was in Artifactory, and not the file name which was downloaded. To achieve this,
// we added a new implementation of the 'Unarchive' func and use it instead of the default one.
func Unarchive(localArchivePath, originArchiveName, destinationPath string) error {
	uaIface, err := byExtension(originArchiveName)
	if err != nil {
		return err
	}
	u, ok := uaIface.(archiver.Unarchiver)
	if !ok {
		return errorutils.CheckError(errors.New("format specified by source filename is not an archive format: " + originArchiveName))
	}
	return u.Unarchive(localArchivePath, destinationPath)
}

// Instead of using 'archiver.byExtension' that by default sets OverwriteExisting to false, we implement our own.
func byExtension(filename string) (interface{}, error) {
	var ec interface{}
	for _, c := range extCheckers {
		if err := c.CheckExt(filename); err == nil {
			ec = c
			break
		}
	}
	switch ec.(type) {
	case *archiver.Rar:
		archiveInstance := archiver.NewRar()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.Tar:
		archiveInstance := archiver.NewTar()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarBrotli:
		archiveInstance := archiver.NewTarBrotli()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarBz2:
		archiveInstance := archiver.NewTarBz2()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarGz:
		archiveInstance := archiver.NewTarGz()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarLz4:
		archiveInstance := archiver.NewTarLz4()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarSz:
		archiveInstance := archiver.NewTarSz()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarXz:
		archiveInstance := archiver.NewTarXz()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarZstd:
		archiveInstance := archiver.NewTarZstd()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.Zip:
		archiveInstance := archiver.NewZip()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.Gz:
		return archiver.NewGz(), nil
	case *archiver.Bz2:
		return archiver.NewBz2(), nil
	case *archiver.Lz4:
		return archiver.NewLz4(), nil
	case *archiver.Snappy:
		return archiver.NewSnappy(), nil
	case *archiver.Xz:
		return archiver.NewXz(), nil
	case *archiver.Zstd:
		return archiver.NewZstd(), nil
	}
	return nil, fmt.Errorf("format unrecognized by filename: %s", filename)
}

var extCheckers = []archiver.ExtensionChecker{
	&archiver.TarBrotli{},
	&archiver.TarBz2{},
	&archiver.TarGz{},
	&archiver.TarLz4{},
	&archiver.TarSz{},
	&archiver.TarXz{},
	&archiver.TarZstd{},
	&archiver.Rar{},
	&archiver.Tar{},
	&archiver.Zip{},
	&archiver.Brotli{},
	&archiver.Gz{},
	&archiver.Bz2{},
	&archiver.Lz4{},
	&archiver.Snappy{},
	&archiver.Xz{},
	&archiver.Zstd{},
}

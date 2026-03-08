package fileutils

import (
	"archive/zip"
	"errors"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io"
	"os"
	"path/filepath"
)

func ZipFolderFiles(source, target string) (err error) {
	zipFile, err := os.Create(target)
	if err != nil {
		return errorutils.CheckError(err)
	}
	defer func() {
		if zipFile != nil {
			err = errors.Join(err, errorutils.CheckError(zipFile.Close()))
		}
	}()

	archive := zip.NewWriter(zipFile)
	defer func() {
		err = errors.Join(err, errorutils.CheckError(archive.Close()))
	}()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) (currentErr error) {
		if info.IsDir() {
			return
		}

		if currentErr = errors.Join(currentErr, err); currentErr != nil {
			return
		}

		header, currentErr := zip.FileInfoHeader(info)
		if currentErr != nil {
			return errorutils.CheckError(currentErr)
		}

		header.Method = zip.Deflate
		writer, currentErr := archive.CreateHeader(header)
		if currentErr != nil {
			return errorutils.CheckError(currentErr)
		}

		file, currentErr := os.Open(path)
		if currentErr != nil {
			return errorutils.CheckError(currentErr)
		}
		defer func() {
			if file != nil {
				currentErr = errors.Join(currentErr, errorutils.CheckError(file.Close()))
			}
		}()
		_, currentErr = io.Copy(writer, file)
		return
	})
}

package fileutils

import (
	"archive/zip"
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
		if cerr := zipFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	archive := zip.NewWriter(zipFile)
	defer func() {
		if cerr := archive.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) (currentErr error) {
		if info.IsDir() {
			return
		}

		if err != nil {
			currentErr = err
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
			if cerr := file.Close(); cerr != nil && currentErr == nil {
				currentErr = cerr
			}
		}()
		_, currentErr = io.Copy(writer, file)
		return
	})
}

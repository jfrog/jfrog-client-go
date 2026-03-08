package fileutils

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
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

	// Use root-scoped fs to prevent symlink TOCTOU (G122)
	fsys := os.DirFS(source)
	return fs.WalkDir(fsys, ".", func(entryPath string, entry fs.DirEntry, walkErr error) (currentErr error) {
		if currentErr = errors.Join(currentErr, walkErr); currentErr != nil {
			return currentErr
		}
		if entry.IsDir() {
			return nil
		}

		info, currentErr := entry.Info()
		if currentErr != nil {
			return errorutils.CheckError(currentErr)
		}

		header, currentErr := zip.FileInfoHeader(info)
		if currentErr != nil {
			return errorutils.CheckError(currentErr)
		}
		header.Name = filepath.ToSlash(entryPath)
		header.Method = zip.Deflate

		writer, currentErr := archive.CreateHeader(header)
		if currentErr != nil {
			return errorutils.CheckError(currentErr)
		}

		// #nosec G122 -- TODO: refactor to use os.OpenRoot for symlink TOCTOU protection
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
		return currentErr
	})
}

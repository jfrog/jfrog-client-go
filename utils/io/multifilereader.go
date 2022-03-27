package io

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io"
	"os"
	"sort"
)

// Create new multi file ReaderAt
func NewMultiFileReaderAt(filePaths []string) (readerAt *multiFileReaderAt, err error) {
	readerAt = &multiFileReaderAt{}
	for _, v := range filePaths {
		f, curErr := os.Open(v)
		if curErr != nil {
			return nil, curErr
		}
		defer func() {
			e := f.Close()
			if err == nil {
				err = errorutils.CheckError(e)
			}
		}()
		stat, curErr := f.Stat()
		if curErr != nil {
			return nil, curErr
		}

		readerAt.filesPaths = append(readerAt.filesPaths, v)
		readerAt.sizeIndex = append(readerAt.sizeIndex, readerAt.size)
		readerAt.size += stat.Size()
	}

	return readerAt, nil
}

type multiFileReaderAt struct {
	filesPaths []string
	size       int64
	sizeIndex  []int64
}

// Get overall size of all the files.
func (multiFileReader *multiFileReaderAt) Size() int64 {
	return multiFileReader.size
}

// ReadAt implementation for multi files
func (multiFileReader *multiFileReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	// Search for the correct index to find the correct file offset
	i := sort.Search(len(multiFileReader.sizeIndex), func(i int) bool { return multiFileReader.sizeIndex[i] > off }) - 1

	readBytes := 0
	for {
		var f *os.File
		f, err = os.Open(multiFileReader.filesPaths[i])
		if err != nil {
			return
		}
		defer func() {
			e := f.Close()
			if err == nil {
				err = errorutils.CheckError(e)
			}
		}()
		relativeOff := off + int64(n) - multiFileReader.sizeIndex[i]
		readBytes, err = f.ReadAt(p[n:], relativeOff)
		n += readBytes
		if len(p) == n {
			// Finished reading enough bytes
			return
		}
		if err != nil && err != io.EOF {
			// Error
			return
		}
		if i+1 == len(multiFileReader.filesPaths) {
			// No more files to read from
			return
		}
		// Read from the next file
		i++
	}
	// not suppose to get here
}

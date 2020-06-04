package content

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

const (
	searchResult      = "SearchResult.json"
	emptySearchResult = "EmptySearchResult.json"
	arrayKey          = "results"
)

type inputRecord struct {
	IntKey   int          `json:"intKey"`
	StrKey   string       `json:"strKey"`
	BoolKey  bool         `json:"boolKey"`
	ArrayKey []ArrayValue `json:"arrayKey"`
}

type ArrayValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func TestContentReaderPath(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	assert.Equal(t, cr.GetFilePath(), searchResultPath)
}

func TestContentReaderNextRecord(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	// Read the same file two times
	for i := 0; i < 2; i++ {
		var rSlice []inputRecord
		for item := new(inputRecord); cr.NextRecord(item) == nil; item = new(inputRecord) {
			rSlice = append(rSlice, *item)
		}
		assert.NoError(t, cr.GetError())
		// First element
		assert.Equal(t, rSlice[0].IntKey, 1)
		assert.Equal(t, rSlice[0].StrKey, "A")
		assert.Equal(t, rSlice[0].BoolKey, true)
		assert.ElementsMatch(t, rSlice[0].ArrayKey, []ArrayValue{{Key: "build.number", Value: "6"}})
		// Second element
		assert.Equal(t, rSlice[1].IntKey, 2)
		assert.Equal(t, rSlice[1].StrKey, "B")
		assert.Equal(t, rSlice[1].BoolKey, false)
		assert.Empty(t, rSlice[1].ArrayKey)
		// Length validation
		assert.Equal(t, cr.Length(), 2)
		cr.Reset()
	}
}

func TestContentReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", emptySearchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	for item := new(inputRecord); cr.NextRecord(item) == nil; item = new(inputRecord) {
		t.Error("Can't loop over empty file")
	}
	assert.NoError(t, cr.GetError())
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "tests", "testsdata")
}

func TestCloseReader(t *testing.T) {
	// Create a file.
	fd, err := fileutils.CreateReaderWriterTempFile()
	assert.NoError(t, err)
	fd.Close()
	filePathToBeDeleted := fd.Name()

	// Load file to reader
	cr := NewContentReader(filePathToBeDeleted, arrayKey)

	// Check file exists
	_, err = os.Stat(filePathToBeDeleted)
	assert.NoError(t, err)

	// Check if the file got deleted
	err = cr.Close()
	assert.NoError(t, err)
	_, err = os.Stat(filePathToBeDeleted)
	assert.True(t, os.IsNotExist(err))
}

func TestDuplicate(t *testing.T) {
	// Create files
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	dupCr, err := cr.Duplicate()
	// Don't delete the origin testdata file, only the duplicate.
	defer dupCr.Close()
	assert.NoError(t, err)

	// Create md5
	originMd5, err := getFileMd5(cr.filePath)
	assert.NoError(t, err)
	expectedMd5, err := getFileMd5(dupCr.filePath)
	assert.NoError(t, err)
	assert.Equal(t, originMd5, expectedMd5)
}

func TestLengthCount(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	assert.Equal(t, cr.Length(), 2)
	// Check cache works with no Reset() being called.
	assert.Equal(t, cr.Length(), 2)
}

func getFileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}

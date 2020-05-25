package content

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

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
	rr := NewContentReader(searchResultPath, arrayKey)
	assert.Equal(t, rr.GetFilePath(), searchResultPath)
}

func TestContentReader(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", searchResult)
	rr := NewContentReader(searchResultPath, arrayKey)
	// Read the same file two times
	for i := 0; i < 2; i++ {
		var rSlice []inputRecord
		var r inputRecord
		var err error
		for err = rr.NextRecord(&r); err == nil; err = rr.NextRecord(&r) {
			rSlice = append(rSlice, r)
		}
		assert.Equal(t, err, io.EOF)
		assert.NoError(t, rr.GetError())
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
		rr.Reset()
	}
}

func TestContentReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "contentreaderwriter", emptySearchResult)
	rr := NewContentReader(searchResultPath, arrayKey)
	var r inputRecord
	for e := rr.NextRecord(&r); e == nil; e = rr.NextRecord(&r) {
		t.Error("Can't loop over empty file")
	}
	assert.NoError(t, rr.GetError())
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "tests", "testsdata")
}

func TestCloseReader(t *testing.T) {
	// Create a file.
	fd, err := ioutil.TempFile("", strconv.FormatInt(time.Now().Unix(), 10))
	assert.NoError(t, err)
	fd.Close()
	filePathToBeDeleted := fd.Name()

	// Load file to reader
	rr := NewContentReader(filePathToBeDeleted, arrayKey)

	// Check file exists
	_, err = os.Stat(filePathToBeDeleted)
	assert.NoError(t, err)

	// Check if the file got deleted
	err = rr.Close()
	assert.NoError(t, err)
	_, err = os.Stat(filePathToBeDeleted)
	assert.True(t, os.IsNotExist(err))
}

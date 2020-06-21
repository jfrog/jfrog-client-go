package content

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
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

func init() {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
}

func TestContentReaderPath(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	assert.Equal(t, cr.GetFilePath(), searchResultPath)
}

func TestContentReaderNextRecord(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	// Read the same file two times
	for i := 0; i < 2; i++ {
		var rSlice []inputRecord
		for item := new(inputRecord); cr.NextRecord(item) == nil; item = new(inputRecord) {
			rSlice = append(rSlice, *item)
		}
		assert.NoError(t, cr.GetError())
		// First element
		assert.Equal(t, 1, rSlice[0].IntKey)
		assert.Equal(t, "A", rSlice[0].StrKey)
		assert.Equal(t, true, rSlice[0].BoolKey)
		assert.ElementsMatch(t, rSlice[0].ArrayKey, []ArrayValue{{Key: "build.number", Value: "6"}})
		// Second element
		assert.Equal(t, 2, rSlice[1].IntKey)
		assert.Equal(t, "B", rSlice[1].StrKey)
		assert.Equal(t, false, rSlice[1].BoolKey)
		assert.Empty(t, rSlice[1].ArrayKey)
		// Length validation
		len, err := cr.Length()
		assert.NoError(t, err)
		assert.Equal(t, 2, len)
		cr.Reset()
	}
}

func TestContentReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), emptySearchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	for item := new(inputRecord); cr.NextRecord(item) == nil; item = new(inputRecord) {
		t.Error("Can't loop over empty file")
	}
	assert.NoError(t, cr.GetError())
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "tests", "testsdata", "contentreaderwriter")
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

func TestLengthCount(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), searchResult)
	cr := NewContentReader(searchResultPath, arrayKey)
	len, err := cr.Length()
	assert.NoError(t, err)
	assert.Equal(t, len, 2)
	// Check cache works with no Reset() being called.
	len, err = cr.Length()
	assert.NoError(t, err)
	assert.Equal(t, len, 2)
}

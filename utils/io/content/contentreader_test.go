package content

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
)

const (
	searchResult      = "SearchResult.json"
	emptySearchResult = "EmptySearchResult.json"
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
	reader := NewContentReader(searchResultPath, DefaultKey)
	assert.Equal(t, reader.GetFilePath(), searchResultPath)
}

func TestContentReaderNextRecord(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), searchResult)
	reader := NewContentReader(searchResultPath, DefaultKey)
	// Read the same file two times
	for i := 0; i < 2; i++ {
		var rSlice []inputRecord
		for item := new(inputRecord); reader.NextRecord(item) == nil; item = new(inputRecord) {
			rSlice = append(rSlice, *item)
		}
		assert.NoError(t, reader.GetError())
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
		len, err := reader.Length()
		assert.NoError(t, err)
		assert.Equal(t, 2, len)
		reader.Reset()
	}
}

func TestContentReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), emptySearchResult)
	reader := NewContentReader(searchResultPath, DefaultKey)
	for item := new(inputRecord); reader.NextRecord(item) == nil; item = new(inputRecord) {
		t.Error("Can't loop over empty file")
	}
	assert.NoError(t, reader.GetError())
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "tests", "testdata", "contentreaderwriter")
}

func TestCloseReader(t *testing.T) {
	// Create a file.
	fd, err := fileutils.CreateTempFile()
	assert.NoError(t, err)
	fd.Close()
	filePathToBeDeleted := fd.Name()

	// Load file to reader
	reader := NewContentReader(filePathToBeDeleted, DefaultKey)

	// Check file exists
	_, err = os.Stat(filePathToBeDeleted)
	assert.NoError(t, err)

	// Check if the file got deleted
	assert.NoError(t, reader.Close())
	_, err = os.Stat(filePathToBeDeleted)
	assert.True(t, os.IsNotExist(err))
}

func TestLengthCount(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), searchResult)
	reader := NewContentReader(searchResultPath, DefaultKey)
	len, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, len, 2)
	// Check cache works with no Reset() being called.
	len, err = reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, len, 2)
}

func TestMergeIncreasingSortedFiles(t *testing.T) {
	testDataPath := getTestDataPath()
	var sortedFiles []*ContentReader
	for i := 1; i <= 3; i++ {
		sortedFiles = append(sortedFiles, NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("buffer_file_ascending_order_%v.json", i)), DefaultKey))
	}
	resultReader, err := MergeSortedReaders(ReaderTestItem{}, sortedFiles, true)
	assert.NoError(t, err)
	isMatch, err := fileutils.FilesIdentical(resultReader.GetFilePath(), filepath.Join(testDataPath, "merged_buffer_ascending_order.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())
}

func TestMergeDecreasingSortedFiles(t *testing.T) {
	testDataPath := getTestDataPath()
	var sortedFiles []*ContentReader
	for i := 1; i <= 3; i++ {
		sortedFiles = append(sortedFiles, NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("buffer_file_descending_order_%v.json", i)), DefaultKey))
	}
	resultReader, err := MergeSortedReaders(ReaderTestItem{}, sortedFiles, false)
	assert.NoError(t, err)
	isMatch, err := fileutils.FilesIdentical(resultReader.GetFilePath(), filepath.Join(testDataPath, "merged_buffer_descending_order.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	assert.NoError(t, resultReader.Close())
}

type ReaderTestItem struct {
	Repo string `json:"repo,omitempty"`
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func (rti ReaderTestItem) GetSortKey() string {
	return path.Join(rti.Repo, rti.Path, rti.Name)
}

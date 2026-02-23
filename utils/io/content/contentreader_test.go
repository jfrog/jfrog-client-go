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
	unsortedFile      = "UnsortedFile.json"
	sortedFile        = "SortedFile.json"
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
	assert.Equal(t, 1, len(reader.GetFilesPaths()))
	assert.Equal(t, searchResultPath, reader.GetFilesPaths()[0])
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
		getErrorAndAssert(t, reader)
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
		length, err := reader.Length()
		assert.NoError(t, err)
		assert.Equal(t, 2, length)
		reader.Reset()
	}
}

func TestContentReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), emptySearchResult)
	reader := NewContentReader(searchResultPath, DefaultKey)
	for item := new(inputRecord); reader.NextRecord(item) == nil; item = new(inputRecord) {
		t.Error("Can't loop over empty file")
	}
	getErrorAndAssert(t, reader)
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "tests", "testdata", "contentreaderwriter")
}

func TestCloseReader(t *testing.T) {
	// Create a file.
	fd, err := fileutils.CreateTempFile()
	assert.NoError(t, err)
	assert.NoError(t, fd.Close())
	filePathToBeDeleted := fd.Name()

	// Load file to reader
	reader := NewContentReader(filePathToBeDeleted, DefaultKey)

	// Check file exists
	_, err = os.Stat(filePathToBeDeleted) // #nosec G703 -- test file; path from test temp
	assert.NoError(t, err)

	// Check if the file got deleted
	closeAndAssert(t, reader)
	_, err = os.Stat(filePathToBeDeleted) // #nosec G703 -- test file; path from test temp
	assert.True(t, os.IsNotExist(err))
}

func TestLengthCount(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), searchResult)
	reader := NewContentReader(searchResultPath, DefaultKey)
	length, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, length, 2)
	// Check cache works with no Reset() being called.
	length, err = reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, length, 2)
}

func TestMergeIncreasingSortedFiles(t *testing.T) {
	testDataPath := getTestDataPath()
	var sortedFiles []*ContentReader
	for i := 1; i <= 3; i++ {
		sortedFiles = append(sortedFiles, NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("buffer_file_ascending_order_%v.json", i)), DefaultKey))
	}
	resultReader, err := MergeSortedReaders(ReaderTestItem{}, sortedFiles, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err := fileutils.JsonEqual(resultReader.GetFilesPaths()[0], filepath.Join(testDataPath, "merged_buffer_ascending_order.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	closeAndAssert(t, resultReader)
}

func TestMergeDecreasingSortedFiles(t *testing.T) {
	testDataPath := getTestDataPath()
	var sortedFiles []*ContentReader
	for i := 1; i <= 3; i++ {
		sortedFiles = append(sortedFiles, NewContentReader(filepath.Join(testDataPath, fmt.Sprintf("buffer_file_descending_order_%v.json", i)), DefaultKey))
	}
	resultReader, err := MergeSortedReaders(ReaderTestItem{}, sortedFiles, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resultReader.GetFilesPaths()))
	isMatch, err := fileutils.JsonEqual(resultReader.GetFilesPaths()[0], filepath.Join(testDataPath, "merged_buffer_descending_order.json"))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	closeAndAssert(t, resultReader)
}

func TestSortContentReaderByCalculatedKey(t *testing.T) {
	testDataPath := getTestDataPath()
	unsortedFilePath := filepath.Join(testDataPath, unsortedFile)
	reader := NewContentReader(unsortedFilePath, DefaultKey)

	getSortKeyFunc := func(result interface{}) (string, error) {
		resultItem := new(ReaderTestItem)
		err := ConvertToStruct(result, &resultItem)
		if err != nil {
			return "", err
		}
		return resultItem.Name, nil
	}

	sortedReader, err := SortContentReaderByCalculatedKey(reader, getSortKeyFunc, true)
	assert.NoError(t, err)
	isMatch, err := fileutils.JsonEqual(sortedReader.GetFilesPaths()[0], filepath.Join(testDataPath, sortedFile))
	assert.NoError(t, err)
	assert.True(t, isMatch)
	closeAndAssert(t, sortedReader)
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

func closeAndAssert(t *testing.T, reader *ContentReader) {
	assert.NoError(t, reader.Close(), "Couldn't close reader")
}

func getErrorAndAssert(t *testing.T, reader *ContentReader) {
	assert.NoError(t, reader.GetError(), "Couldn't get reader error")
}

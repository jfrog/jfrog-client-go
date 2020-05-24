package jsonreaderwriter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	searchResult      = "SearchResult.json"
	emptySearchResult = "EmptySearchResult.json"
)

func TestJsonReaderPath(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "jsonreaderwriter", searchResult)
	rr := NewJsonReader(searchResultPath, arrayKey)
	assert.Equal(t, rr.GetFilePath(), searchResultPath)
}

func TestJsonReader(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "jsonreaderwriter", searchResult)
	rr := NewJsonReader(searchResultPath, arrayKey)
	assert.Equal(t, rr.GetFilePath(), searchResultPath)

	channel, channelErr := rr.Run()
	for data := range channel {
		rawJson, err := json.Marshal(data)
		assert.NoError(t, err)
		x := string(rawJson)
		assert.Equal(t, x, `{"properties":[{"key":"build.number","value":"6"}],"repo":"MyRepo"}`)
	}
	assert.NoError(t, channelErr.GetError())

}

func TestJsonReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "jsonreaderwriter", emptySearchResult)
	rr := NewJsonReader(searchResultPath, arrayKey)
	channel, channelErr := rr.Run()
	for range channel {
		t.Error("Can't loop over empty file")
	}
	assert.NoError(t, channelErr.GetError())

}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "..", "tests", "testsdata")
}

package responsereaderwriter

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

func TestResponseReaderPath(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "responsereaderwriter", searchResult)
	rr := NewResponseReader(searchResultPath)
	assert.Equal(t, rr.GetFilePath(), searchResultPath)
}

func TestResponseReader(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "responsereaderwriter", searchResult)
	rr := NewResponseReader(searchResultPath)
	assert.Equal(t, rr.GetFilePath(), searchResultPath)

	channel, channelErr := rr.Run()
	for data := range channel {
		rawJson, err := json.Marshal(data)
		assert.NoError(t, err)
		x := string(rawJson)
		assert.Equal(t, x, `{"properties":[{"key":"build.number","value":"6"}],"repo":"MyRepo"}`)
	}
	assert.NoError(t, channelErr)

}

func TestResponseReaderEmptyResult(t *testing.T) {
	searchResultPath := filepath.Join(getTestDataPath(), "responsereaderwriter", emptySearchResult)
	rr := NewResponseReader(searchResultPath)
	channel, channelErr := rr.Run()
	for range channel {
		t.Error("Can't loop over empty file")
	}
	assert.NoError(t, channelErr)

}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "..", "..", "tests", "testsdata")
}

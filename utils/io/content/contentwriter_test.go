package content

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

type outputRecord struct {
	IntKey  int    `json:"intKey"`
	StrKey  string `json:"strKey"`
	BoolKey bool   `json:"boolKey"`
}

var records = []outputRecord{
	{1, "1", true},
	{2, "2", false},
	{3, "3", true},
	{4, "4", false},
	{5, "5", false},
	{6, "6", true},
	{7, "7", true},
	{8, "8", true},
	{9, "9", false},
	{10, "10", false},
	{11, "11", false},
	{12, "12", true},
	{13, "13", false},
	{14, "14", true},
	{15, "15", true},
	{16, "16", true},
	{17, "17", false},
	{18, "18", true},
	{19, "19", false},
	{20, "20", false},
	{21, "21", true},
	{22, "22", true},
	{23, "23", true},
	{24, "24", false},
	{25, "25", false},
	{26, "26", false},
	{27, "27", true},
	{28, "28", false},
	{29, "29", true},
	{30, "30", true},
}

type Response struct {
	Arr []outputRecord `json:"arr"`
}

func writeTestRecords(t *testing.T, cw *ContentWriter) {
	var sendersWaiter sync.WaitGroup
	for i := 0; i < len(records); i += 3 {
		sendersWaiter.Add(1)
		go func(start, end int) {
			defer sendersWaiter.Done()
			for j := start; j < end; j++ {
				cw.Write(records[j])
			}
		}(i, i+3)
	}
	sendersWaiter.Wait()
	assert.NoError(t, cw.Close())
}

func TestContentWriter(t *testing.T) {
	rw, err := NewContentWriter("arr", true, false)
	assert.NoError(t, err)
	writeTestRecords(t, rw)
	of, err := os.Open(rw.GetFilePath())
	assert.NoError(t, err)
	byteValue, _ := ioutil.ReadAll(of)
	var response Response
	assert.NoError(t, json.Unmarshal(byteValue, &response))
	assert.NoError(t, of.Close())
	assert.NoError(t, rw.RemoveOutputFilePath())
	for i := range records {
		assert.Contains(t, response.Arr, records[i], "record %s missing", records[i].StrKey)
	}
}

func TestContentReaderAfterWriter(t *testing.T) {
	cw, err := NewContentWriter("results", true, false)
	assert.NoError(t, err)
	writeTestRecords(t, cw)
	cr := NewContentReader(cw.GetFilePath(), "results")
	assert.NoError(t, err)
	defer cr.Close()
	recordCount := 0
	for item := new(outputRecord); cr.NextRecord(item) == nil; item = new(outputRecord) {
		assert.Contains(t, records, *item, "record %s missing", item.StrKey)
		recordCount++
	}
	assert.NoError(t, cr.GetError())
	assert.Equal(t, len(records), recordCount, "The amount of records were read (%d) is different then expected", recordCount)
}

func TestRemoveOutputFilePath(t *testing.T) {
	// Create a file.
	cw, err := NewContentWriter("results", true, false)
	assert.NoError(t, err)
	assert.NoError(t, cw.Close())
	filePathToBeDeleted := cw.GetFilePath()

	// Check file exists
	_, err = os.Stat(filePathToBeDeleted)
	assert.NoError(t, err)

	// Check if the file got deleted
	cw.RemoveOutputFilePath()
	_, err = os.Stat(filePathToBeDeleted)
	assert.True(t, os.IsNotExist(err))
}

func TestEmptyContentWriter(t *testing.T) {
	cw, err := NewEmptyContentWriter("results", true, false)
	assert.NoError(t, err)
	searchResultPath := filepath.Join(getTestDataPath(), emptySearchResult)
	result, err := fileutils.FilesMath(cw.GetFilePath(), searchResultPath)
	assert.NoError(t, err)
	assert.True(t, result)

	cw, err = NewContentWriter("results", true, false)
	assert.NoError(t, err)
	assert.NoError(t, cw.Close())
	searchResultPath = filepath.Join(getTestDataPath(), emptySearchResult)
	assert.NoError(t, err)
	result, err = fileutils.FilesMath(cw.GetFilePath(), searchResultPath)
	assert.NoError(t, err)
	assert.True(t, result)
}

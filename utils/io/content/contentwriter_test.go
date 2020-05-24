package content

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type outputRecord struct {
	IntKey  int    `json:"intKey"`
	StrKey  string `json:"strKey"`
	BoolKey bool   `json:boolKey`
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

func writeTestRecords(t *testing.T, rw *ContentWriter) {
	var sendersWaiter sync.WaitGroup
	rw.Run()
	for i := 0; i < len(records); i += 3 {
		sendersWaiter.Add(1)
		go func(start, end int) {
			defer sendersWaiter.Done()
			for j := start; j < end; j++ {
				rw.Write(records[j])
			}
		}(i, i+3)
	}
	sendersWaiter.Wait()
	err := rw.Done()
	assert.NoError(t, err)
}

func TestContentWriter(t *testing.T) {
	rw, err := NewContentWriter(5, "arr", true, false)
	assert.NoError(t, err)
	writeTestRecords(t, rw)
	of, err := os.Open(rw.GetOutputFilePath())
	assert.NoError(t, err)
	byteValue, _ := ioutil.ReadAll(of)
	var response Response
	err = json.Unmarshal(byteValue, &response)
	assert.NoError(t, err)
	err = of.Close()
	assert.NoError(t, err)
	err = rw.RemoveOutputFilePath()
	assert.NoError(t, err)
	for i, _ := range records {
		assert.Contains(t, response.Arr, records[i], "record %s missing", records[i].StrKey)
	}
}

func TestResponseReadeAfterWriter(t *testing.T) {
	rw, err := NewContentWriter(5, "results", true, false)
	assert.NoError(t, err)
	writeTestRecords(t, rw)
	rr := NewContentReader(rw.GetOutputFilePath(), "results")
	assert.NoError(t, err)
	_, errQueue := rr.Run()
	recordCount := 0
	var r outputRecord
	for e := rr.GetRecord(&r); e == nil; e = rr.GetRecord(&r) {
		assert.Contains(t, records, r, "record %s missing", r.StrKey)
		recordCount++
	}
	assert.NoError(t, errQueue.GetError())
	assert.Equal(t, len(records), recordCount, "The amount of records were read (%d) is different then expected", recordCount)
}

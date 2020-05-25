package content

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	channelSize = 100
)

// Open and read JSON file, find the array key inside it and load its value to the memory in small chunks.
// Only array objects are valides to be fetch.
// Each chunk can be fetch using 'GetRecord' (thread-safe).
//
// This technick solve the limit of memory size which may be too small to fit.
// That way, we can handle big queries result from Artifactory that may not fit in size to the host total memory.
type ContentReader struct {
	// filePath - Source data file path.
	// arrayKey = Read the value of the specific object in JSON.
	filePath, arrayKey string
	// The objects from the source data file are being pushed into the data channel.
	dataChannel chan map[string]interface{}
	errorsQueue *utils.ErrorsQueue
	once        *sync.Once
}

func NewContentReader(filePath string, arrayKey string) *ContentReader {
	self := ContentReader{}
	self.filePath = filePath
	self.arrayKey = arrayKey
	self.dataChannel = make(chan map[string]interface{}, channelSize)
	self.errorsQueue = utils.NewErrorsQueue(channelSize)
	self.once = new(sync.Once)
	return &self
}

func (rc *ContentReader) ArrayKey(arrayKey string) *ContentReader {
	rc.arrayKey = arrayKey
	return rc
}

// fetch the next chunk into 'recordOutput' param.
// 'io.EOF' will be returned if no data is left.
func (rc *ContentReader) NextRecord(recordOutput interface{}) error {
	rc.once.Do(func() {
		go func() {
			rc.run()
		}()
	})
	record, ok := <-rc.dataChannel
	if !ok {
		return errorutils.CheckError(io.EOF)
	}
	data, _ := json.Marshal(record)
	return errorutils.CheckError(json.Unmarshal(data, recordOutput))
}

// Initialize the reader to read a file that has already been read (not thread-safe).
func (rc *ContentReader) Reset() {
	rc.dataChannel = make(chan map[string]interface{}, channelSize)
	rc.once = new(sync.Once)
}

// Cleanup the reader data.
func (rc *ContentReader) Close() error {
	if rc.filePath != "" {
		if err := errorutils.CheckError(os.Remove(rc.filePath)); err != nil {
			return err
		}
		rc.filePath = ""
	}
	return nil
}

func (rc *ContentReader) GetFilePath() string {
	return rc.filePath
}

func (rc *ContentReader) SetFilePath(newPath string) error {
	if rc.filePath != "" {
		if err := rc.Close(); err != nil {
			return err
		}
	}
	rc.filePath = newPath
	rc.dataChannel = make(chan map[string]interface{}, channelSize)
	return nil
}

// Open and read the file. Push each array element to the channel.
// The channel may block the thread, therefore should run async.
func (rc *ContentReader) run() {
	fd, err := os.Open(rc.filePath)
	if err != nil {
		log.Error(err.Error())
		rc.errorsQueue.AddError(errorutils.CheckError(err))
		return
	}
	br := bufio.NewReaderSize(fd, 65536)
	defer fd.Close()
	defer close(rc.dataChannel)
	dec := json.NewDecoder(br)
	err = findDecoderTargetPosition(dec, rc.arrayKey, true)
	if err != nil {
		if err == io.EOF {
			rc.errorsQueue.AddError(errors.New("results not found"))
			return
		}
		rc.errorsQueue.AddError(err)
		log.Error(err.Error())
		return
	}
	for dec.More() {
		var ResultItem map[string]interface{}
		err := dec.Decode(&ResultItem)
		if err != nil {
			log.Error(err)
			rc.errorsQueue.AddError(errorutils.CheckError(err))
			return
		}
		rc.dataChannel <- ResultItem
	}
}

// Return true if the file has more than one element in array.
func (rc *ContentReader) IsEmpty() (bool, error) {
	fd, err := os.Open(rc.filePath)
	if err != nil {
		log.Error(err.Error())
		rc.errorsQueue.AddError(errorutils.CheckError(err))
		return false, err
	}
	br := bufio.NewReaderSize(fd, 65536)
	defer fd.Close()
	defer close(rc.dataChannel)
	dec := json.NewDecoder(br)
	err = findDecoderTargetPosition(dec, rc.arrayKey, true)
	return isEmptyArray(dec, rc.arrayKey, true)
}

func (rc *ContentReader) GetError() error {
	return rc.errorsQueue.GetError()
}

// Search and set the decoder's position at the desired key in the JSON file.
func findDecoderTargetPosition(dec *json.Decoder, target string, isArray bool) error {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return errorutils.CheckError(err)
		}
		if t == target {
			if isArray {
				// Skip '['
				_, err = dec.Token()
			}
			return errorutils.CheckError(err)
		}
	}
	return nil
}

// Scan the JSON file and check if the array contains at least one element.
func isEmptyArray(dec *json.Decoder, target string, isArray bool) (bool, error) {
	if err := findDecoderTargetPosition(dec, target, isArray); err != nil {
		return false, err
	}
	t, err := dec.Token()
	if err != nil {
		return false, errorutils.CheckError(err)
	}
	return t == json.Delim('{'), nil
}

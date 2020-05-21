package responsereaderwriter

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type ResponseReader struct {
	// Response data file path.
	filePath string
	// The objects from the response data file are being pushed to the data channel.
	dataChannel chan map[string]interface{}
}

func NewResponseReader(filePath string) *ResponseReader {
	self := ResponseReader{}
	self.filePath = filePath
	self.dataChannel = make(chan map[string]interface{}, 2)
	return &self
}

// Fire up a goroutine in order to fill the data channel.
func (rr *ResponseReader) Run() (chan map[string]interface{}, error) {
	var err error
	go func() {
		err = rr.run()
	}()
	return rr.dataChannel, err
}

// Iterator to get next record from the file.
// The file be deleted and io.EOF error will be returned when there are no more records in the channel and the channel is closed.
func (rr *ResponseReader) GetRecord(recordOutput interface{}) error {
	record, ok := <-rr.dataChannel
	if !ok {
		rr.DeleteFile()
		return io.EOF
	}
	data, _ := json.Marshal(record)
	return json.Unmarshal(data, recordOutput)
}

func (rr *ResponseReader) DeleteFile() error {
	return os.Remove(rr.filePath)
}

func (rr *ResponseReader) GetFilePath() string {
	return rr.filePath
}

func (rr *ResponseReader) SetFilePath(newPath string) {
	rr.filePath = newPath
	rr.dataChannel = make(chan map[string]interface{}, 2)
}

func (rr *ResponseReader) IsZeroResults() (bool, error) {
	fd, err := os.Open(rr.filePath)
	br := bufio.NewReaderSize(fd, 65536)
	defer fd.Close()
	if err != nil {
		return false, err
	}
	dec := json.NewDecoder(br)
	err = findDecoderTargetPosition(dec, "results", true)
	if err != nil {
		return false, err
	}
	t, err := dec.Token()
	if err != nil {
		return false, err
	}
	return t != json.Delim('{'), nil
}

func (rr *ResponseReader) run() error {
	fd, err := os.Open(rr.filePath)
	if err != nil {
		x := err.Error()
		log.Fatal(x)
		return err
	}
	br := bufio.NewReaderSize(fd, 65536)
	defer fd.Close()
	defer close(rr.dataChannel)
	dec := json.NewDecoder(br)
	err = findDecoderTargetPosition(dec, "results", true)
	if err != nil {
		if err == io.EOF {
			return errors.New("results not found")
		}
		return err
	}
	for dec.More() {
		var ResultItem map[string]interface{}
		err := dec.Decode(&ResultItem)
		if err != nil {
			log.Fatal(err)
			return err
		}
		rr.dataChannel <- ResultItem
	}
	return err
}

func findDecoderTargetPosition(dec *json.Decoder, target string, isArray bool) error {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		if t == target {
			if isArray {
				// Skip '['
				_, err = dec.Token()
			}
			return err
		}
	}
	return nil
}

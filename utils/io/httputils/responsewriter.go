package httputils

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

type ResponseReader struct {
	// Response data file path.
	filePath string
	// The objects from the response data file are being pushed to the data channel.
	dataChannel chan map[string]interface{}
}

func NewResponseWriter(filePath string) (*ResponseReader, error) {
	self := ResponseReader{}
	self.filePath = filePath
	self.dataChannel = make(chan map[string]interface{}, 2)
	return &self, nil
}

// Fire up a goroutine in order to fill the data channel.
func (rw *ResponseReader) Run() (chan map[string]interface{}, error) {
	var err error
	go func() {
		err = rw.run()
	}()
	return rw.dataChannel, err
}

func (rw *ResponseReader) DeleteFiles() error {
	return os.RemoveAll(filepath.Dir(rw.filePath))
}

func (rw *ResponseReader) run() error {
	fd, err := os.Open(rw.filePath)
	br := bufio.NewReaderSize(fd, 65536)
	defer fd.Close()
	defer close(rw.dataChannel)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(br)
	err = findDecoderTargetPosition(dec, "results", true)
	if err != nil {
		if err == io.EOF {
			return errors.New("results not found")
		}
		return err
	}
	var ResultItem map[string]interface{}
	for dec.More() {
		err := dec.Decode(&ResultItem)
		if err != nil {
			log.Fatal(err)
			return err
		}
		rw.dataChannel <- ResultItem
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

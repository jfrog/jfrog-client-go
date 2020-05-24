package content

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	arrayKey = "results"
)

type DataDrive struct {
	Reader      *ContentReader
	Writer      *ContentWriter
	buffer      []map[string]interface{}
	bufferSize  int
	channel     chan map[string]interface{}
	errorsQueue *utils.ErrorsQueue
}

func NewDataDrive(bufferSize int) *DataDrive {
	self := DataDrive{}
	if bufferSize <= 0 {
		bufferSize = 50000
	}
	self.bufferSize = bufferSize
	self.channel = make(chan map[string]interface{}, bufferSize)
	self.errorsQueue = utils.NewErrorsQueue(10)
	self.buffer = make([]map[string]interface{}, 0)
	return &self
}

func NewDataDriveWithStream(src io.Reader, bufferSize int) (*DataDrive, error) {
	self := *NewDataDrive(bufferSize)
	dec := json.NewDecoder(src)
	err := findDecoderTargetPosition(dec, arrayKey, true)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New(arrayKey + " not found")
		}
		return nil, err
	}
	for dec.More() && bufferSize > 0 {
		var ResultItem map[string]interface{}
		err := dec.Decode(&ResultItem)
		if err != nil {
			log.Fatal(err)
			return nil, errorutils.CheckError(err)
		}
		self.buffer = append(self.buffer, ResultItem)
		bufferSize--
	}
	if bufferSize == 0 {
		rw, err := NewContentWriter(self.bufferSize, arrayKey, true, false)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		rw.Run()
		for dec.More() {
			var ResultItem map[string]interface{}
			err := dec.Decode(&ResultItem)
			if err != nil {
				log.Fatal(err)
				return nil, errorutils.CheckError(err)
			}
			rw.Write(ResultItem)
		}
		rw.Done()
		self.Reader = NewContentReader(rw.GetOutputFilePath(), arrayKey)
	}
	return &self, nil
}

func NewDataDriveWithFile(filePath string, bufferSize int) (*DataDrive, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
		return nil, errorutils.CheckError(err)
	}
	bufio := bufio.NewReaderSize(fd, 65536)
	defer fd.Close()
	return NewDataDriveWithStream(bufio, bufferSize)
}

func (dd *DataDrive) run() {
	// In case of running more than once, the channel is closed
	dd.channel = make(chan map[string]interface{}, dd.bufferSize)
	go func() {
		var runWaiter sync.WaitGroup
		runWaiter.Add(1)
		go func() {
			for _, bufferChunk := range dd.buffer {
				dd.channel <- bufferChunk
			}
			runWaiter.Done()
		}()
		if dd.Reader != nil {
			runWaiter.Add(1)
			go func() {
				channel, errQueue := dd.Reader.Run()
				for readerItem := range channel {
					dd.channel <- readerItem
				}
				dd.errorsQueue = errQueue
				runWaiter.Done()
			}()
		}
		runWaiter.Wait()
		defer close(dd.channel)
	}()
}

func (dd *DataDrive) RunReader() (chan map[string]interface{}, *utils.ErrorsQueue) {
	dd.run()
	return dd.channel, dd.errorsQueue
}

func (dd *DataDrive) AddRecord(dataToWrite map[string]interface{}) error {
	var err error
	if len(dd.buffer) < dd.bufferSize {
		dd.buffer = append(dd.buffer, dataToWrite)
	} else {
		if dd.Writer == nil {
			dd.Writer, err = NewContentWriter(dd.bufferSize, "results", true, false)
			if err != nil {
				return err
			}
			dd.Writer.Run()
		}
		dd.Writer.Write(dataToWrite)
	}
	return nil
}

func (dd *DataDrive) CloseWriter() error {
	if len(dd.buffer) == dd.bufferSize {
		if err := dd.Writer.Done(); err != nil {
			return err
		}
		dd.Reader.SetFilePath(dd.Writer.GetOutputFilePath())
	}
	return nil
}

func (dd *DataDrive) Empty() bool {
	return len(dd.buffer) == 0
}

func (dd *DataDrive) Cleanup() error {
	if dd.Reader != nil {
		if err := dd.Reader.Close(); err != nil {
			return err
		}
	}
	if dd.Writer != nil {
		return dd.Writer.RemoveOutputFilePath()
	}
	return nil
}

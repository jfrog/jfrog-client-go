package responsereaderwriter

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io/ioutil"
	"os"
	"sync"
)

const (
	outputFilePattern      = "%s.*.json"
	jsonArrayPrefixPattern = "{\n\t\"%s\" : [\n"
	jsonArraySuffix        = "\t]\n}"
)

type ResponseWriter struct {
	arrayKey string
	// The output data file path.
	outputFile *os.File
	// The chanel which from the output records will be pulled.
	recordsChannel chan interface{}
	recordCount    int
	errorsQueue    *utils.ErrorsQueue
	runWaiter      sync.WaitGroup
}

func NewResponseWriter(chanCapacity int, arrayKey string) (*ResponseWriter, error) {
	fd, err := ioutil.TempFile("", fmt.Sprintf(outputFilePattern, arrayKey))
	if err != nil {
		return nil, err
	}
	self := ResponseWriter{}
	self.arrayKey = arrayKey
	self.outputFile = fd
	self.recordsChannel = make(chan interface{}, chanCapacity)
	self.errorsQueue = utils.NewErrorsQueue(chanCapacity)
	self.recordCount = 0
	return &self, nil
}

func (rw *ResponseWriter) GetOutputFilePath() string {
	return rw.outputFile.Name()
}

func (rw *ResponseWriter) GetRecordCount() int {
	return rw.recordCount
}

func (rw *ResponseWriter) RemoveOutputFilePath() error {
	return errorutils.CheckError(os.Remove(rw.outputFile.Name()))
}

func (rw *ResponseWriter) AddRecord(record interface{}) {
	rw.recordsChannel <- record
	rw.recordCount++
}

func (rw *ResponseWriter) Run() {
	rw.runWaiter.Add(1)
	go func() {
		defer rw.runWaiter.Done()
		rw.run()
	}()
	return
}

func (rw *ResponseWriter) run() {
	defer rw.outputFile.Close()
	_, err := rw.outputFile.WriteString(fmt.Sprintf(jsonArrayPrefixPattern, rw.arrayKey))
	if err != nil {
		rw.errorsQueue.AddError(err)
		return
	}
	enc := json.NewEncoder(rw.outputFile)
	recordPrefix := "\t\t"
	for record := range rw.recordsChannel {
		rw.outputFile.WriteString(recordPrefix)
		err = enc.Encode(record)
		if err != nil {
			rw.errorsQueue.AddError(err)
		}
		recordPrefix = "\t\t,"
	}
	_, err = rw.outputFile.WriteString(jsonArraySuffix)
	if err != nil {
		rw.errorsQueue.AddError(err)
	}
	return
}

func (rw *ResponseWriter) Stop() error {
	close(rw.recordsChannel)
	rw.runWaiter.Wait()
	return rw.errorsQueue.GetError()
}

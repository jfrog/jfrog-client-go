package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	outputFilePattern      = "%s.*.json"
	jsonArrayPrefixPattern = "  \"%s\": ["
	jsonArraySuffix        = "]\n"
)

type ContentWriter struct {
	arrayKey string
	// The output data file path.
	outputFile *os.File
	// The chanel which from the output records will be pulled.
	recordsChannel chan interface{}
	isCompleteFile bool
	errorsQueue    *utils.ErrorsQueue
	runWaiter      sync.WaitGroup
}

func NewContentWriter(chanCapacity int, arrayKey string, isCompleteFile, useStdout bool) (*ContentWriter, error) {
	var fd *os.File
	var err error
	if useStdout {
		fd = os.Stdout
	} else {
		fd, err = ioutil.TempFile("", fmt.Sprintf(outputFilePattern, arrayKey))
		if err != nil {
			return nil, errorutils.CheckError(err)
		}
	}
	self := ContentWriter{}
	self.arrayKey = arrayKey
	self.outputFile = fd
	self.recordsChannel = make(chan interface{}, chanCapacity)
	self.errorsQueue = utils.NewErrorsQueue(chanCapacity)
	self.isCompleteFile = isCompleteFile
	return &self, nil
}

func (rw *ContentWriter) SetArrayKey(arrKey string) *ContentWriter {
	rw.arrayKey = arrKey
	return rw
}

func (rw *ContentWriter) GetOutputFilePath() string {
	return rw.outputFile.Name()
}

func (rw *ContentWriter) RemoveOutputFilePath() error {
	return errorutils.CheckError(os.Remove(rw.outputFile.Name()))
}

func (rw *ContentWriter) Write(record interface{}) {
	rw.recordsChannel <- record
}

func (rw *ContentWriter) Run() {
	rw.runWaiter.Add(1)
	go func() {
		defer rw.runWaiter.Done()
		rw.run()
	}()
	return
}

func (rw *ContentWriter) run() {
	if rw.outputFile != os.Stdout {
		defer rw.outputFile.Close()
	}
	openString := jsonArrayPrefixPattern
	closeString := ""
	if rw.isCompleteFile {
		openString = "{\n" + openString
	}
	_, err := rw.outputFile.WriteString(fmt.Sprintf(openString, rw.arrayKey))
	if err != nil {
		rw.errorsQueue.AddError(errorutils.CheckError(err))
		return
	}
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetIndent("    ", "  ")
	recordPrefix := "\n    "
	firstRecord := true
	for record := range rw.recordsChannel {
		rw.outputFile.WriteString(recordPrefix)
		err = enc.Encode(record)
		if err != nil {
			rw.errorsQueue.AddError(errorutils.CheckError(err))
		}
		b := bytes.TrimRight(buf.Bytes(), "\n")
		_, err = rw.outputFile.Write(b)
		if err != nil {
			rw.errorsQueue.AddError(errorutils.CheckError(err))
		}
		buf.Reset()
		if firstRecord {
			// If a record was printed, we want to print a comma and ne before each and every future record.
			recordPrefix = "," + recordPrefix
			// We will close the array in a new-indent line.
			closeString = "\n  "
			firstRecord = false
		}
	}
	closeString = closeString + jsonArraySuffix
	if rw.isCompleteFile {
		closeString += "}\n"
	}
	_, err = rw.outputFile.WriteString(closeString)
	if err != nil {
		rw.errorsQueue.AddError(errorutils.CheckError(err))
	}
	return
}

func (rw *ContentWriter) Done() error {
	close(rw.recordsChannel)
	rw.runWaiter.Wait()
	return rw.errorsQueue.GetError()
}

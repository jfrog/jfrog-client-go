package httputils

import (
	"time"

	"github.com/jfrog/jfrog-client-go/utils"
)

type PollingAction func() (shouldStop bool, responseBody []byte, err error)

type PollingExecutor struct {
	// Maximum wait time in nanoseconds.
	Timeout time.Duration
	// Number of nanoseconds to sleep between polling attempts.
	PollingInterval time.Duration
	// Prefix to add at the beginning of each info/error message.
	MsgPrefix string
	// pollingAction is the operation to run until the condition fulfilled.
	PollingAction PollingAction
}

func (runner *PollingExecutor) Execute() ([]byte, error) {
	var finalResponse []byte
	retryExecutor := utils.RetryExecutor{
		MaxRetries:               int(runner.Timeout.Seconds() / (runner.PollingInterval.Seconds())),
		RetriesIntervalMilliSecs: int(runner.PollingInterval.Milliseconds()),
		ErrorMessage:             "",
		LogMsgPrefix:             runner.MsgPrefix,
		ExecutionHandler: func() (bool, error) {
			shouldStop, response, err := runner.PollingAction()
			finalResponse = response
			return !shouldStop, err
		},
	}
	return finalResponse, retryExecutor.Execute()
}

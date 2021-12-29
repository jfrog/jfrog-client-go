package utils

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strconv"
	"time"
)

type ExecutionHandlerFunc func() (bool, error)

type RetryExecutor struct {
	// The amount of retries to perform.
	MaxRetries int

	// Number of milliseconds to sleep between retries.
	RetriesIntervalMilliSecs int

	// Message to display when retrying.
	ErrorMessage string

	// Prefix to print at the beginning of each log.
	LogMsgPrefix string

	// ExecutionHandler is the operation to run with retries.
	ExecutionHandler ExecutionHandlerFunc
}

func (runner *RetryExecutor) Execute() error {
	var err error
	var shouldRetry bool
	for i := 0; i <= runner.MaxRetries; i++ {
		// Run ExecutionHandler
		shouldRetry, err = runner.ExecutionHandler()

		// If should not retry, return
		if !shouldRetry {
			return err
		}

		log.Warn(runner.getLogRetryMessage(i, err))
		// Going to sleep for RetryInterval milliseconds
		if runner.RetriesIntervalMilliSecs > 0 && i < runner.MaxRetries {
			log.Info("Waiting ", strconv.Itoa(runner.RetriesIntervalMilliSecs), "ms before trying again")
			time.Sleep(time.Millisecond * time.Duration(runner.RetriesIntervalMilliSecs))
		}
	}

	return err
}

func (runner *RetryExecutor) getLogRetryMessage(attemptNumber int, err error) (message string) {
	message = fmt.Sprintf("%sAttempt %v - %s", runner.LogMsgPrefix, attemptNumber, runner.ErrorMessage)
	if err != nil {
		message = fmt.Sprintf("%s - %s", message, err.Error())
	}
	return
}

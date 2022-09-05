package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ExecutionHandlerFunc func() (bool, error)

type RetryExecutor struct {
	// The context
	Context context.Context

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
		if cancelledErr := runner.checkCancelled(); cancelledErr != nil {
			return cancelledErr
		}

		// Print retry log message
		runner.LogRetry(i, err)

		// Going to sleep for RetryInterval milliseconds
		if runner.RetriesIntervalMilliSecs > 0 && i < runner.MaxRetries {
			time.Sleep(time.Millisecond * time.Duration(runner.RetriesIntervalMilliSecs))
		}
	}
	log.Info(fmt.Sprintf("%s executor timeout after %v attempts with %v milliseconds wait intervals", runner.LogMsgPrefix, runner.MaxRetries, runner.RetriesIntervalMilliSecs))
	return err
}

func (runner *RetryExecutor) LogRetry(attemptNumber int, err error) {
	message := fmt.Sprintf("%s(Attempt %v)", runner.LogMsgPrefix, attemptNumber+1)
	if runner.ErrorMessage != "" {
		message = fmt.Sprintf("%s - %s", message, runner.ErrorMessage)
	}
	if err != nil {
		message = fmt.Sprintf("%s: %s", message, err.Error())
	}

	if err != nil || runner.ErrorMessage != "" {
		log.Warn(message)
	} else {
		log.Debug(message)
	}
}

func (runner *RetryExecutor) checkCancelled() error {
	if runner.Context == nil {
		return nil
	}
	contextErr := runner.Context.Err()
	if errors.Is(contextErr, context.Canceled) {
		log.Info("Retry executor was cancelled")
		return contextErr
	}
	return nil
}

package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ExecutionHandlerFunc func() (shouldRetry bool, err error)

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

		// If we should not retry, return.
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
	// If the error is not nil, return it and log the timeout message. Otherwise, generate new error.
	if err != nil {
		log.Info(runner.getTimeoutErrorMsg())
		return err
	}
	return errorutils.CheckError(RetryExecutorTimeoutError{runner.getTimeoutErrorMsg()})
}

// Error of this type will be returned if the executor reaches timeout and no other error is returned by the execution handler.
type RetryExecutorTimeoutError struct {
	errMsg string
}

func (retryErr RetryExecutorTimeoutError) Error() string {
	return retryErr.errMsg
}

func (runner *RetryExecutor) getTimeoutErrorMsg() string {
	prefix := ""
	if runner.LogMsgPrefix != "" {
		prefix = runner.LogMsgPrefix + " "
	}
	return fmt.Sprintf("%sexecutor timeout after %v attempts with %v milliseconds wait intervals", prefix, runner.MaxRetries, runner.RetriesIntervalMilliSecs)
}

func (runner *RetryExecutor) LogRetry(attemptNumber int, err error) {
	var message string
	if attemptNumber == 0 {
		message = fmt.Sprintf("%s[Initial attempt]", runner.LogMsgPrefix)
	} else {
		message = fmt.Sprintf("%s[Retry %d/%d]", runner.LogMsgPrefix, attemptNumber, runner.MaxRetries)
	}

	if runner.ErrorMessage != "" {
		message = fmt.Sprintf("%s %s", message, runner.ErrorMessage)
	}

	if err != nil {
		message = fmt.Sprintf("%s - Error: %s", message, err.Error())
	}

	if err != nil {
		if attemptNumber == runner.MaxRetries {
			log.Error(fmt.Sprintf("%s - Final attempt failed", message))
		} else {
			log.Warn(message)
		}
	} else {
		if attemptNumber == 0 {
			log.Debug(fmt.Sprintf("%s - Succeeded on first try", message))
		} else {
			log.Info(fmt.Sprintf("%s - Succeeded after %d retries", message, attemptNumber))
		}
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

package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestRetryExecutorSuccess(t *testing.T) {
	retriesToPerform := 10
	breakRetriesAt := 4
	runCount := 0
	executor := RetryExecutor{
		MaxRetries:               retriesToPerform,
		RetriesIntervalMilliSecs: 0,
		ErrorMessage:             "Testing RetryExecutor",
		ExecutionHandler: func() (bool, error) {
			runCount++
			if runCount == breakRetriesAt {
				log.Warn("Breaking after", runCount-1, "retries")
				return false, nil
			}
			return true, nil
		},
	}

	assert.NoError(t, executor.Execute())
	assert.Equal(t, breakRetriesAt, runCount)
}

func TestRetryExecutorTimeoutWithDefaultError(t *testing.T) {
	retriesToPerform := 5
	runCount := 0

	executor := RetryExecutor{
		MaxRetries:               retriesToPerform,
		RetriesIntervalMilliSecs: 0,
		ErrorMessage:             "Testing RetryExecutor",
		ExecutionHandler: func() (bool, error) {
			runCount++
			return true, nil
		},
	}

	assert.Equal(t, executor.Execute(), RetryExecutorTimeoutError{executor.getTimeoutErrorMsg()})
	assert.Equal(t, retriesToPerform+1, runCount)
}

func TestRetryExecutorTimeoutWithCustomError(t *testing.T) {
	retriesToPerform := 5
	runCount := 0

	executionHandler := errors.New("retry failed due to reason")

	executor := RetryExecutor{
		MaxRetries:               retriesToPerform,
		RetriesIntervalMilliSecs: 0,
		ErrorMessage:             "Testing RetryExecutor",
		ExecutionHandler: func() (bool, error) {
			runCount++
			return true, executionHandler
		},
	}

	assert.Equal(t, executor.Execute(), executionHandler)
	assert.Equal(t, retriesToPerform+1, runCount)
}

func TestRetryExecutorCancel(t *testing.T) {
	retriesToPerform := 5
	runCount := 0

	retryContext, cancelFunc := context.WithCancel(context.Background())
	executor := RetryExecutor{
		Context:                  retryContext,
		MaxRetries:               retriesToPerform,
		RetriesIntervalMilliSecs: 0,
		ErrorMessage:             "Testing RetryExecutor",
		ExecutionHandler: func() (bool, error) {
			runCount++
			return true, nil
		},
	}

	cancelFunc()
	assert.EqualError(t, executor.Execute(), context.Canceled.Error())
	assert.Equal(t, 1, runCount)
}

package utils

import (
	"context"
	"fmt"
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
	if runCount != breakRetriesAt {
		t.Error(fmt.Errorf("expected, %d, got: %d", breakRetriesAt, runCount))
	}
}

func TestRetryExecutorFail(t *testing.T) {
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

	assert.NoError(t, executor.Execute())
	if runCount != retriesToPerform+1 {
		t.Error(fmt.Errorf("expected, %d, got: %d", retriesToPerform, runCount))
	}
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
	if runCount != 1 {
		t.Error(fmt.Errorf("expected, %d, got: %d", retriesToPerform, runCount))
	}
}

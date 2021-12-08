package httputils

import (
	"time"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type PollingAction func() (shouldStop bool, responseBody []byte, err error)

type PollingExecutor struct {
	// Maximum wait time in seconds.
	Timeout time.Duration
	// Number of seconds to sleep between polling attempts.
	PollingInterval time.Duration
	// Prefix to add at the beginning of each info/error message.
	MsgPrefix string
	// pollingAction is the operation to run until the condition fullfiled.
	PollingAction PollingAction
}

func (runner *PollingExecutor) Execute() ([]byte, error) {
	ticker := time.NewTicker(runner.PollingInterval)
	timeout := make(chan bool)
	errChan := make(chan error)
	resultChan := make(chan []byte)
	go func() {
		for {
			select {
			case <-timeout:
				errChan <- errorutils.CheckErrorf("%s Polling executor timeout after %v secondes", runner.MsgPrefix, runner.Timeout.Seconds())
				resultChan <- nil
				return
			case _ = <-ticker.C:
				shouldStop, responseBody, err := runner.PollingAction()
				if err != nil {
					errChan <- err
					resultChan <- nil
					return
				}
				// Got the full valid response.
				if shouldStop {
					errChan <- nil
					resultChan <- responseBody
					return
				}
			}
		}
	}()
	// Make sure we don't wait forever
	go func() {
		time.Sleep(runner.Timeout)
		timeout <- true
	}()
	// Wait for result or error
	err := <-errChan
	body := <-resultChan
	ticker.Stop()
	return body, err
}

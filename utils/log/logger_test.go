package log

import (
	"bytes"
	"sync"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/stretchr/testify/assert"
)

func TestLoggerRemoveEmojis(t *testing.T) {
	testLoggerWithEmojis(t, false, expectedLogOutputWithoutEmojis)
}

func TestLoggerLeaveEmojis(t *testing.T) {
	expected := expectedLogOutputWithEmojis
	if io.IsWindows() {
		// Should not print emojis on Windows
		expected = expectedLogOutputWithoutEmojis
	}
	testLoggerWithEmojis(t, true, expected)
}

func testLoggerWithEmojis(t *testing.T, mockIsTerminalFlags bool, expected string) {
	previousLog := Logger
	// Restore previous logger when the function returns.
	defer SetLogger(previousLog)

	// Set new logger with output redirection to buffer.
	buffer := &bytes.Buffer{}
	SetLogger(NewLogger(DEBUG, buffer))
	if mockIsTerminalFlags {
		// Mock logger with isTerminal flags set to true
		revertFlags := SetIsTerminalFlagsWithCallback(true)
		// Revert to previous status
		defer revertFlags()
	}
	Debug("111", 111, "", "111😀111👻🪶")
	Info("222", 222, "", "222😀222👻🪶")
	Warn("333", 333, "", "333😀333👻🪶")
	Error("444", 444, "", "444😀444👻🪶")
	Output("555", 555, "", "555😀555👻🪶")

	// Compare output.
	logOutput := buffer.Bytes()
	compareResult := bytes.Compare(logOutput, []byte(expected))
	assert.Equal(t, 0, compareResult)
}

const expectedLogOutputWithoutEmojis = `[Debug] 111 111  111111
[Info] 222 222  222222
[Warn] 333 333  333333
[Error] 444 444  444444
555 555  555555
`

const expectedLogOutputWithEmojis = `[Debug] 111 111  111😀111👻🪶
[Info] 222 222  222😀222👻🪶
[Warn] 333 333  333😀333👻🪶
[Error] 444 444  444😀444👻🪶
555 555  555😀555👻🪶
`

// BufferedLogger tests

func TestBufferedLogger_CapturesAllLevels(t *testing.T) {
	logger := NewBufferedLogger(DEBUG)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	logger.Output("output message")

	assert.Equal(t, 5, logger.Len())

	output := logger.String()
	assert.Contains(t, output, "[DEBUG] debug message")
	assert.Contains(t, output, "[INFO] info message")
	assert.Contains(t, output, "[WARN] warn message")
	assert.Contains(t, output, "[ERROR] error message")
	assert.Contains(t, output, "output message")
}

func TestBufferedLogger_RespectsLogLevel(t *testing.T) {
	logger := NewBufferedLogger(WARN)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Only WARN and ERROR should be captured (DEBUG and INFO filtered out)
	assert.Equal(t, 2, logger.Len())

	output := logger.String()
	assert.NotContains(t, output, "debug message")
	assert.NotContains(t, output, "info message")
	assert.Contains(t, output, "[WARN] warn message")
	assert.Contains(t, output, "[ERROR] error message")
}

func TestBufferedLogger_OutputBypassesLogLevel(t *testing.T) {
	logger := NewBufferedLogger(ERROR)

	logger.Debug("debug message")
	logger.Output("output message")

	// Output should always be captured regardless of log level
	assert.Equal(t, 1, logger.Len())
	assert.Contains(t, logger.String(), "output message")
}

func TestBufferedLogger_Clear(t *testing.T) {
	logger := NewBufferedLogger(DEBUG)

	logger.Info("message 1")
	logger.Info("message 2")
	assert.Equal(t, 2, logger.Len())

	logger.Clear()
	assert.Equal(t, 0, logger.Len())
	assert.Empty(t, logger.String())
}

func TestBufferedLogger_ReplayTo(t *testing.T) {
	source := NewBufferedLogger(DEBUG)
	source.Debug("debug msg")
	source.Info("info msg")
	source.Warn("warn msg")
	source.Error("error msg")
	source.Output("output msg")

	target := NewBufferedLogger(DEBUG)
	source.ReplayTo(target)

	assert.Equal(t, 5, target.Len())
	output := target.String()
	assert.Contains(t, output, "[DEBUG] debug msg")
	assert.Contains(t, output, "[INFO] info msg")
	assert.Contains(t, output, "[WARN] warn msg")
	assert.Contains(t, output, "[ERROR] error msg")
	assert.Contains(t, output, "output msg")
}

func TestBufferedLogger_ThreadSafe(t *testing.T) {
	logger := NewBufferedLogger(DEBUG)
	var wg sync.WaitGroup

	// Spawn multiple goroutines writing concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			logger.Info("message from goroutine")
		}(i)
	}

	wg.Wait()
	assert.Equal(t, 100, logger.Len())
}

// Goroutine-local logger tests

func TestSetLoggerForGoroutine_IsolatesLoggers(t *testing.T) {
	previousLog := Logger
	defer SetLogger(previousLog)

	globalBuffer := &bytes.Buffer{}
	SetLogger(NewLogger(DEBUG, globalBuffer))

	goroutineBuffer := NewBufferedLogger(DEBUG)

	done := make(chan bool)
	go func() {
		SetLoggerForGoroutine(goroutineBuffer)
		defer ClearLoggerForGoroutine()

		// This should go to goroutineBuffer, not globalBuffer
		GetLogger().Info("goroutine message")
		done <- true
	}()

	<-done

	// Global logger should not have the message
	assert.NotContains(t, globalBuffer.String(), "goroutine message")

	// Goroutine-local logger should have the message
	assert.Contains(t, goroutineBuffer.String(), "goroutine message")
}

func TestGetLogger_ReturnsGlobalWhenNoGoroutineLogger(t *testing.T) {
	previousLog := Logger
	defer SetLogger(previousLog)

	buffer := &bytes.Buffer{}
	globalLogger := NewLogger(DEBUG, buffer)
	SetLogger(globalLogger)

	// No goroutine-local logger set, should return global
	assert.Equal(t, globalLogger, GetLogger())
}

func TestClearLoggerForGoroutine_ReturnsToGlobal(t *testing.T) {
	previousLog := Logger
	defer SetLogger(previousLog)

	globalBuffer := &bytes.Buffer{}
	SetLogger(NewLogger(DEBUG, globalBuffer))

	goroutineBuffer := NewBufferedLogger(DEBUG)

	done := make(chan bool)
	go func() {
		SetLoggerForGoroutine(goroutineBuffer)
		GetLogger().Info("before clear")

		ClearLoggerForGoroutine()
		GetLogger().Info("after clear")
		done <- true
	}()

	<-done

	// "before clear" should be in goroutine buffer
	assert.Contains(t, goroutineBuffer.String(), "before clear")
	assert.NotContains(t, goroutineBuffer.String(), "after clear")

	// "after clear" should be in global buffer
	assert.Contains(t, globalBuffer.String(), "after clear")
}

func TestGoroutineLoggers_DifferentGoroutinesDifferentLoggers(t *testing.T) {
	buffer1 := NewBufferedLogger(DEBUG)
	buffer2 := NewBufferedLogger(DEBUG)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		SetLoggerForGoroutine(buffer1)
		defer ClearLoggerForGoroutine()
		GetLogger().Info("message from goroutine 1")
	}()

	go func() {
		defer wg.Done()
		SetLoggerForGoroutine(buffer2)
		defer ClearLoggerForGoroutine()
		GetLogger().Info("message from goroutine 2")
	}()

	wg.Wait()

	// Each buffer should only have its own message
	assert.Contains(t, buffer1.String(), "message from goroutine 1")
	assert.NotContains(t, buffer1.String(), "message from goroutine 2")

	assert.Contains(t, buffer2.String(), "message from goroutine 2")
	assert.NotContains(t, buffer2.String(), "message from goroutine 1")
}

package log

import (
	"bytes"
	"github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestLoggerForceEmojis(t *testing.T) {
	// Force emojis
	assert.False(t, allowEmojis)
	testLoggerWithEmojis(t, true, expectedLogOutputWithEmojis)
	callback := SetAllowEmojiFlagWithCallback(true)
	assert.True(t, allowEmojis)
	testLoggerWithEmojis(t, true, expectedLogOutputWithEmojis)
	callback()
	assert.False(t, allowEmojis)
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

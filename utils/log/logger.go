package log

import (
	"fmt"
	"github.com/forPelevin/gomoji"
	"github.com/gookit/color"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"runtime"
)

var Logger Log

type LevelType int
type LogFormat string

// Used for coloring sections of the log message. For example log.Format.Path("...")
var Format LogFormat

// Determines whether the terminal is available. This variable should not be accessed directly,
// but through the 'isTerminalMode' function.
var terminalMode *bool

// Determines whether emoji is supported. This variable should not be accessed directly,
// but through the 'isEmojiSupported' function.
var emojiSupported *bool

// defaultLogger is the default logger instance in case the user does not set one
var defaultLogger = NewLogger(INFO, nil)

const (
	ERROR LevelType = iota
	WARN
	INFO
	DEBUG
)

// Creates a new logger with a given LogLevel.
// All logs are written to Stderr by default (output to Stdout).
// If logToWriter != nil, logging is done to the provided writer instead.
// Log flags to modify the log prefix as described in https://pkg.go.dev/log#pkg-constants.
func NewLoggerWithFlags(logLevel LevelType, logToWriter io.Writer, logFlags int) *jfrogLogger {
	logger := new(jfrogLogger)
	logger.SetLogLevel(logLevel)
	logger.SetOutputWriter(logToWriter)
	logger.SetLogsWriter(logToWriter, logFlags)
	return logger
}

// Same as NewLoggerWithFlags, with log flags turned off.
func NewLogger(logLevel LevelType, logToWriter io.Writer) *jfrogLogger {
	return NewLoggerWithFlags(logLevel, logToWriter, 0)
}

type jfrogLogger struct {
	LogLevel  LevelType
	OutputLog *log.Logger
	DebugLog  *log.Logger
	InfoLog   *log.Logger
	WarnLog   *log.Logger
	ErrorLog  *log.Logger
}

func SetLogger(newLogger Log) {
	Logger = newLogger
}

func GetLogger() Log {
	if Logger != nil {
		return Logger
	}
	return defaultLogger
}

func (logger *jfrogLogger) SetLogLevel(LevelEnum LevelType) {
	logger.LogLevel = LevelEnum
}

func (logger *jfrogLogger) SetOutputWriter(writer io.Writer) {
	if writer == nil {
		writer = os.Stdout
	}
	logger.OutputLog = log.New(writer, "", 0)
}

func (logger *jfrogLogger) Println(log *log.Logger, values ...interface{}) {
	if !isEmojiSupported() {
		for _, value := range values {
			switch value.(type) {
			case string:
				value = gomoji.RemoveEmojis(value.(string))
			}
		}
	}
	log.Println(values...)
}

// Set the logs' writer to Stderr unless an alternative one is provided.
// In case the writer is set for file, colors will not be in use.
// Log flags to modify the log prefix as described in https://pkg.go.dev/log#pkg-constants.
func (logger *jfrogLogger) SetLogsWriter(writer io.Writer, logFlags int) {
	writerIsStdErr := false
	if writer == nil {
		writer = os.Stderr
		writerIsStdErr = true
	}
	logger.DebugLog = log.New(writer, getLogPrefix(DEBUG, writerIsStdErr), logFlags)
	logger.InfoLog = log.New(writer, getLogPrefix(INFO, writerIsStdErr), logFlags)
	logger.WarnLog = log.New(writer, getLogPrefix(WARN, writerIsStdErr), logFlags)
	logger.ErrorLog = log.New(writer, getLogPrefix(ERROR, writerIsStdErr), logFlags)
}

type prefixStyle struct {
	logLevel string
	color    color.Color
	emoji    string
}

var prefixStyles = map[LevelType]*prefixStyle{
	DEBUG: {logLevel: "Debug", color: color.Cyan},
	INFO:  {logLevel: "Info", emoji: "🔵", color: color.Blue},
	WARN:  {logLevel: "Warn", emoji: "🟠", color: color.Yellow},
	ERROR: {logLevel: "Error", emoji: "🚨", color: color.Red},
}

func getLogPrefix(logType LevelType, writerIsStdErr bool) string {
	if logPrefixStyle, ok := prefixStyles[logType]; ok {
		prefix := logPrefixStyle.logLevel
		// Use colors only on stdErr terminal output
		if writerIsStdErr && isTerminalMode() {
			prefix = logPrefixStyle.color.Render(prefix)
		}
		if isEmojiSupported() {
			prefix = logPrefixStyle.emoji + prefix
		}
		return fmt.Sprintf("[%s] ", prefix)
	}

	return ""
}

func Debug(a ...interface{}) {
	GetLogger().Debug(a...)
}

func Info(a ...interface{}) {
	GetLogger().Info(a...)
}

func Warn(a ...interface{}) {
	GetLogger().Warn(a...)
}

func Error(a ...interface{}) {
	GetLogger().Error(a...)
}

func Output(a ...interface{}) {
	GetLogger().Output(a...)
}

func (logger jfrogLogger) GetLogLevel() LevelType {
	return logger.LogLevel
}

func (logger jfrogLogger) Debug(a ...interface{}) {
	if logger.GetLogLevel() >= DEBUG {
		logger.Println(logger.DebugLog, a...)
	}
}

func (logger jfrogLogger) Info(a ...interface{}) {
	if logger.GetLogLevel() >= INFO {
		logger.Println(logger.InfoLog, a...)
	}
}

func (logger jfrogLogger) Warn(a ...interface{}) {
	if logger.GetLogLevel() >= WARN {
		logger.Println(logger.WarnLog, a...)
	}
}

func (logger jfrogLogger) Error(a ...interface{}) {
	if logger.GetLogLevel() >= ERROR {
		logger.Println(logger.ErrorLog, a...)
	}
}

func (logger jfrogLogger) Output(a ...interface{}) {
	logger.Println(logger.OutputLog, a...)
}

type Log interface {
	Debug(a ...interface{})
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
	Output(a ...interface{})
}

// Check if Stderr is a terminal
func isTerminalMode() bool {
	if terminalMode == nil {
		t := term.IsTerminal(int(os.Stderr.Fd()))
		terminalMode = &t
	}
	return *terminalMode
}

// Check if Emoji is supported
func isEmojiSupported() bool {
	if emojiSupported == nil {
		t := isTerminalMode() && runtime.GOOS != "windows"
		emojiSupported = &t
	}
	return *emojiSupported
}

// Predefined color formatting functions
func (f *LogFormat) Path(message string) string {
	if isTerminalMode() {
		return color.Green.Render(message)
	}
	return message
}

func (f *LogFormat) URL(message string) string {
	if isTerminalMode() {
		return color.Cyan.Render(message)
	}
	return message
}

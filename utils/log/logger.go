package log

import (
	"fmt"
	"github.com/forPelevin/gomoji"
	"github.com/gookit/color"
	"golang.org/x/term"
	"io"
	termUtils "k8s.io/kubectl/pkg/util/term"
	"log"
	"os"
)

var Logger Log

type LevelType int
type LogFormat string

// Used for coloring sections of the log message. For example log.Format.Path("...")
var Format LogFormat

// Determines whether the terminal is available. This variable should not be accessed directly,
// but through the 'isTerminalMode' function.
var terminalMode *bool

// Determines whether colors supported. This variable should not be accessed directly,
// but through the 'colorsSupported' function.
var colorsSupported *bool

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
	if !isColorsSupported(log.Writer()) {
		for _, value := range values {
			if str, ok := value.(string); ok {
				value = gomoji.RemoveEmojis(str)
			}
		}
	}
	log.Println(values...)
}

// Set the logs' writer to Stderr unless an alternative one is provided.
// In case the writer is set for file, colors will not be in use.
// Log flags to modify the log prefix as described in https://pkg.go.dev/log#pkg-constants.
func (logger *jfrogLogger) SetLogsWriter(writer io.Writer, logFlags int) {
	if writer == nil {
		writer = os.Stderr
	}
	logger.DebugLog = getLogWriter(writer, DEBUG, logFlags)
	logger.InfoLog = getLogWriter(writer, INFO, logFlags)
	logger.WarnLog = getLogWriter(writer, WARN, logFlags)
	logger.ErrorLog = getLogWriter(writer, ERROR, logFlags)
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

func getLogWriter(writer io.Writer, logType LevelType, logFlags int) *log.Logger {
	if logPrefixStyle, ok := prefixStyles[logType]; ok {
		prefix := logPrefixStyle.logLevel
		if isColorsSupported(writer) {
			prefix = logPrefixStyle.emoji + logPrefixStyle.color.Render(prefix)
		}
		return log.New(writer, fmt.Sprintf("[%s] ", prefix), logFlags)
	}
	return log.New(writer, "", logFlags)
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

// Check if Color is supported
func isColorsSupported(writer io.Writer) bool {
	if colorsSupported == nil {
		t := termUtils.AllowsColorOutput(writer)
		colorsSupported = &t
	}
	return *colorsSupported
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

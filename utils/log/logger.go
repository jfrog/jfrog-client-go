package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gookit/color"
	"golang.org/x/term"
)

var Logger Log

type LevelType int
type LogFormat string

// Used for coloring sections of the log message. For example log.Format.Path("...")
var Format LogFormat

// Determines whether the terminal is available. This variable should not be accessed directly,
// but through the 'isTerminalMode' function.
var terminalMode *bool

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

func (logger *jfrogLogger) SetLogLevel(LevelEnum LevelType) {
	logger.LogLevel = LevelEnum
}

func (logger *jfrogLogger) SetOutputWriter(writer io.Writer) {
	if writer == nil {
		writer = os.Stdout
	}
	logger.OutputLog = log.New(writer, "", 0)
}

// Set the logs' writer to Stderr unless an alternative one is provided.
// In case the writer is set for file, colors will not be in use.
// Log flags to modify the log prefix as described in https://pkg.go.dev/log#pkg-constants.
func (logger *jfrogLogger) SetLogsWriter(writer io.Writer, logFlags int) {
	if writer == nil {
		writer = os.Stderr
		if isTerminalMode() {
			logger.DebugLog = log.New(writer, fmt.Sprintf("[%s] ", color.Cyan.Render("Debug")), logFlags)
			logger.InfoLog = log.New(writer, fmt.Sprintf("[🔵%s] ", color.Blue.Render("Info")), logFlags)
			logger.WarnLog = log.New(writer, fmt.Sprintf("[🟠%s] ", color.Yellow.Render("Warn")), logFlags)
			logger.ErrorLog = log.New(writer, fmt.Sprintf("[🚨%s] ", color.Red.Render("Error")), logFlags)
			return
		}
	}
	logger.DebugLog = log.New(writer, "[Debug] ", logFlags)
	logger.InfoLog = log.New(writer, "[Info] ", logFlags)
	logger.WarnLog = log.New(writer, "[Warn] ", logFlags)
	logger.ErrorLog = log.New(writer, "[Error] ", logFlags)
}

func validateLogInit() {
	if Logger == nil {
		panic("Logger not initialized. See API documentation.")
	}
}

func Debug(a ...interface{}) {
	validateLogInit()
	Logger.Debug(a...)
}

func Info(a ...interface{}) {
	validateLogInit()
	Logger.Info(a...)
}

func Warn(a ...interface{}) {
	validateLogInit()
	Logger.Warn(a...)
}

func Error(a ...interface{}) {
	validateLogInit()
	Logger.Error(a...)
}

func Output(a ...interface{}) {
	validateLogInit()
	Logger.Output(a...)
}

func (logger jfrogLogger) GetLogLevel() LevelType {
	return logger.LogLevel
}

func (logger jfrogLogger) Debug(a ...interface{}) {
	if logger.GetLogLevel() >= DEBUG {
		logger.DebugLog.Println(a...)
	}
}

func (logger jfrogLogger) Info(a ...interface{}) {
	if logger.GetLogLevel() >= INFO {
		logger.InfoLog.Println(a...)
	}
}

func (logger jfrogLogger) Warn(a ...interface{}) {
	if logger.GetLogLevel() >= WARN {
		logger.WarnLog.Println(a...)
	}
}

func (logger jfrogLogger) Error(a ...interface{}) {
	if logger.GetLogLevel() >= ERROR {
		logger.ErrorLog.Println(a...)
	}
}

func (logger jfrogLogger) Output(a ...interface{}) {
	logger.OutputLog.Println(a...)
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

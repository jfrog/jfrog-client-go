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

// Determines whether the Stdout is terminal. This variable should not be accessed directly,
// but through the 'IsStdOutTerminal' function.
var stdOutIsTerminal *bool

// Determines whether the Stderr is terminal. This variable should not be accessed directly,
// but through the 'IsStdErrTerminal' function.
var stdErrIsTerminal *bool

// Determines whether colors are supported. This variable should not be accessed directly,
// but through the 'colorsSupported' function.
var colorsSupported *bool

// defaultLogger is the default logger instance in case the user does not set one
var defaultLogger = NewLogger(INFO, nil)

var logWriter io.Writer
var outputWriter io.Writer

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
func NewLoggerWithFlags(logLevel LevelType, writer io.Writer, logFlags int) *jfrogLogger {
	logger := new(jfrogLogger)
	logger.SetLogLevel(logLevel)
	logger.SetOutputWriter(writer)
	logger.SetLogsWriter(writer, logFlags)
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

func (logger *jfrogLogger) SetLogLevel(levelEnum LevelType) {
	logger.LogLevel = levelEnum
}

func (logger *jfrogLogger) SetOutputWriter(writer io.Writer) {
	if writer != nil {
		outputWriter = writer
	} else {
		outputWriter = io.Writer(os.Stdout)
	}
	// Reset outIsTerminal flag
	stdOutIsTerminal = nil
	logger.OutputLog = log.New(outputWriter, "", 0)
}

// Set the logs' writer to Stderr unless an alternative one is provided.
// In case the writer is set for file, colors will not be in use.
// Log flags to modify the log prefix as described in https://pkg.go.dev/log#pkg-constants.
func (logger *jfrogLogger) SetLogsWriter(writer io.Writer, logFlags int) {
	if writer != nil {
		logWriter = writer
	} else {
		logWriter = io.Writer(os.Stderr)
	}
	// reset errIsTerminal flag
	stdErrIsTerminal = nil
	logger.DebugLog = log.New(logWriter, getLogPrefix(DEBUG), logFlags)
	logger.InfoLog = log.New(logWriter, getLogPrefix(INFO), logFlags)
	logger.WarnLog = log.New(logWriter, getLogPrefix(WARN), logFlags)
	logger.ErrorLog = log.New(logWriter, getLogPrefix(ERROR), logFlags)
}

var prefixStyles = map[LevelType]struct {
	logLevel string
	color    color.Color
	emoji    string
}{
	DEBUG: {logLevel: "Debug", color: color.Cyan},
	INFO:  {logLevel: "Info", emoji: "🔵", color: color.Blue},
	WARN:  {logLevel: "Warn", emoji: "🟠", color: color.Yellow},
	ERROR: {logLevel: "Error", emoji: "🚨", color: color.Red},
}

func getLogPrefix(logType LevelType) string {
	if logPrefixStyle, ok := prefixStyles[logType]; ok {
		prefix := logPrefixStyle.logLevel
		// Add emoji and color only if it's a terminal that supports it
		if IsStdErrTerminal() && IsColorsSupported() {
			prefix = logPrefixStyle.emoji + logPrefixStyle.color.Render(prefix)
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
		logger.Println(logger.DebugLog, IsStdErrTerminal(), a...)
	}
}

func (logger jfrogLogger) Info(a ...interface{}) {
	if logger.GetLogLevel() >= INFO {
		logger.Println(logger.InfoLog, IsStdErrTerminal(), a...)
	}
}

func (logger jfrogLogger) Warn(a ...interface{}) {
	if logger.GetLogLevel() >= WARN {
		logger.Println(logger.WarnLog, IsStdErrTerminal(), a...)
	}
}

func (logger jfrogLogger) Error(a ...interface{}) {
	if logger.GetLogLevel() >= ERROR {
		logger.Println(logger.ErrorLog, IsStdErrTerminal(), a...)
	}
}

func (logger jfrogLogger) Output(a ...interface{}) {
	logger.Println(logger.OutputLog, IsStdOutTerminal(), a...)
}

func (logger *jfrogLogger) Println(log *log.Logger, isTerminal bool, values ...interface{}) {
	// Remove emojis from all strings if it's not a terminal or if the terminal is not supporting colors
	if !(IsColorsSupported() && isTerminal) {
		for i, value := range values {
			if str, ok := value.(string); ok {
				if gomoji.ContainsEmoji(str) {
					values[i] = gomoji.RemoveEmojis(str)
				}
			}
		}
	}
	log.Println(values...)
}

type Log interface {
	Debug(a ...interface{})
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
	Output(a ...interface{})
	GetLogLevel() LevelType
}

// Check if StdErr is a terminal
func IsStdErrTerminal() bool {
	if stdErrIsTerminal == nil {
		isTerminal := false
		if v, ok := (logWriter).(*os.File); ok {
			isTerminal = term.IsTerminal(int(v.Fd()))
		}
		stdErrIsTerminal = &isTerminal
	}
	return *stdErrIsTerminal
}

// Check if Stdout is a terminal
func IsStdOutTerminal() bool {
	if stdOutIsTerminal == nil {
		isTerminal := false
		if v, ok := (outputWriter).(*os.File); ok {
			isTerminal = term.IsTerminal(int(v.Fd()))
		}
		stdOutIsTerminal = &isTerminal
	}
	return *stdOutIsTerminal
}

// SetIsTerminalFlagsWithCallback changes IsTerminal flags to the given value and return function that changes the flags back to the original values.
func SetIsTerminalFlagsWithCallback(isTerminal bool) func() {
	stdoutIsTerminalPrev := stdOutIsTerminal
	stdErrIsTerminalPrev := stdErrIsTerminal

	stdOutIsTerminal = &isTerminal
	stdErrIsTerminal = &isTerminal

	return func() {
		stdOutIsTerminal = stdoutIsTerminalPrev
		stdErrIsTerminal = stdErrIsTerminalPrev
	}
}

func IsColorsSupported() bool {
	if colorsSupported == nil {
		supported := true
		if os.Getenv("TERM") == "dumb" ||

			// On Windows WT_SESSION is set by the modern terminal component.
			// Older terminals have poor support for UTF-8, VT escape codes, etc.
			(runtime.GOOS == "windows" && os.Getenv("WT_SESSION") == "") ||

			// https://no-color.org/
			func() bool { _, noColorEnvExists := os.LookupEnv("NO_COLOR"); return noColorEnvExists }() {
			supported = false
		}

		colorsSupported = &supported
	}
	return *colorsSupported
}

// Predefined color formatting functions
func (f *LogFormat) Path(message string) string {
	if IsStdErrTerminal() {
		return color.Green.Render(message)
	}
	return message
}

func (f *LogFormat) URL(message string) string {
	if IsStdErrTerminal() {
		return color.Cyan.Render(message)
	}
	return message
}

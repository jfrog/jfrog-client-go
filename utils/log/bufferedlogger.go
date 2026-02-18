package log

import (
	"fmt"
	"strings"
	"sync"
)

type logEntry struct {
	level LevelType
	msg   string
}

// BufferedLogger captures logs for isolated parallel operations.
type BufferedLogger struct {
	entries  []logEntry
	logLevel LevelType
	mu       sync.Mutex
}

// NewBufferedLogger creates a logger that captures entries for later replay.
func NewBufferedLogger(level LevelType) *BufferedLogger {
	return &BufferedLogger{
		logLevel: level,
	}
}

func (b *BufferedLogger) GetLogLevel() LevelType {
	return b.logLevel
}

func (b *BufferedLogger) SetLogLevel(level LevelType) {
	b.logLevel = level
}

func (b *BufferedLogger) append(level LevelType, a ...interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = append(b.entries, logEntry{level: level, msg: fmt.Sprint(a...)})
}

func (b *BufferedLogger) Verbose(a ...interface{}) {
	if b.logLevel >= VERBOSE {
		b.append(VERBOSE, a...)
	}
}

func (b *BufferedLogger) Debug(a ...interface{}) {
	if b.logLevel >= DEBUG {
		b.append(DEBUG, a...)
	}
}

func (b *BufferedLogger) Info(a ...interface{}) {
	if b.logLevel >= INFO {
		b.append(INFO, a...)
	}
}

func (b *BufferedLogger) Warn(a ...interface{}) {
	if b.logLevel >= WARN {
		b.append(WARN, a...)
	}
}

func (b *BufferedLogger) Error(a ...interface{}) {
	if b.logLevel >= ERROR {
		b.append(ERROR, a...)
	}
}

func (b *BufferedLogger) Output(a ...interface{}) {
	b.append(-1, a...)
}

// ReplayTo outputs captured logs through the target logger (preserving colors).
func (b *BufferedLogger) ReplayTo(target Log) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, entry := range b.entries {
		switch entry.level {
		case VERBOSE:
			target.Verbose(entry.msg)
		case DEBUG:
			target.Debug(entry.msg)
		case INFO:
			target.Info(entry.msg)
		case WARN:
			target.Warn(entry.msg)
		case ERROR:
			target.Error(entry.msg)
		default:
			target.Output(entry.msg)
		}
	}
}

func (b *BufferedLogger) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}

func (b *BufferedLogger) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}

func (b *BufferedLogger) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	var sb strings.Builder
	for _, entry := range b.entries {
		var levelStr string
		switch entry.level {
		case VERBOSE:
			levelStr = "VERBOSE"
		case DEBUG:
			levelStr = "DEBUG"
		case INFO:
			levelStr = "INFO"
		case WARN:
			levelStr = "WARN"
		case ERROR:
			levelStr = "ERROR"
		default:
			sb.WriteString(entry.msg)
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n", levelStr, entry.msg))
	}
	return sb.String()
}

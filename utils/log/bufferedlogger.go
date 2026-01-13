package log

import (
	"fmt"
	"strings"
	"sync"
)

// logEntry represents a single captured log message with its level.
type logEntry struct {
	level LevelType
	msg   string
}

// BufferedLogger implements the Log interface and captures all logs as structured entries.
// This enables isolated log capture for parallel operations - each operation can have
// its own BufferedLogger. Use ReplayTo() to output the captured logs through another logger.
type BufferedLogger struct {
	entries  []logEntry
	logLevel LevelType
	mu       sync.Mutex
}

// NewBufferedLogger creates a new logger that captures log entries.
// Use ReplayTo() to replay the captured logs through another logger (preserving colors).
func NewBufferedLogger(level LevelType) *BufferedLogger {
	return &BufferedLogger{
		logLevel: level,
		entries:  make([]logEntry, 0, 100), // Pre-allocate for typical usage
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
	b.mu.Lock()
	defer b.mu.Unlock()
	// Output is always captured regardless of level
	b.entries = append(b.entries, logEntry{level: -1, msg: fmt.Sprint(a...)})
}

// ReplayTo replays all captured log entries through the target logger.
// This preserves colors, formatting, and timestamps from the target logger.
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
			// Output (level -1) or unknown
			target.Output(entry.msg)
		}
	}
}

// Clear removes all captured log entries.
func (b *BufferedLogger) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}

// Len returns the number of captured log entries.
func (b *BufferedLogger) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}

// String returns all captured log entries as a formatted string.
// For colored output, use ReplayTo() instead.
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
			// Output entries (level -1)
			sb.WriteString(entry.msg)
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n", levelStr, entry.msg))
	}
	return sb.String()
}

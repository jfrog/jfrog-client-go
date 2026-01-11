package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// BufferedLogger implements the Log interface and writes all logs to a custom io.Writer.
// This enables isolated log capture for parallel operations - each operation can have
// its own BufferedLogger writing to its own buffer.
type BufferedLogger struct {
	writer   io.Writer
	logLevel LevelType
	mu       sync.Mutex
}

// NewBufferedLogger creates a new logger that writes to the provided writer.
// Use this with a bytes.Buffer to capture logs for later retrieval.
func NewBufferedLogger(writer io.Writer, level LevelType) *BufferedLogger {
	return &BufferedLogger{
		writer:   writer,
		logLevel: level,
	}
}

func (b *BufferedLogger) GetLogLevel() LevelType {
	return b.logLevel
}

func (b *BufferedLogger) SetLogLevel(level LevelType) {
	b.logLevel = level
}

func (b *BufferedLogger) logf(level string, a ...interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
	msg := fmt.Sprint(a...)
	fmt.Fprintf(b.writer, "[%s] [%s] %s\n", timestamp, level, msg)
}

func (b *BufferedLogger) Verbose(a ...interface{}) {
	if b.logLevel >= VERBOSE {
		b.logf("VERBOSE", a...)
	}
}

func (b *BufferedLogger) Debug(a ...interface{}) {
	if b.logLevel >= DEBUG {
		b.logf("DEBUG", a...)
	}
}

func (b *BufferedLogger) Info(a ...interface{}) {
	if b.logLevel >= INFO {
		b.logf("INFO", a...)
	}
}

func (b *BufferedLogger) Warn(a ...interface{}) {
	if b.logLevel >= WARN {
		b.logf("WARN", a...)
	}
}

func (b *BufferedLogger) Error(a ...interface{}) {
	if b.logLevel >= ERROR {
		b.logf("ERROR", a...)
	}
}

func (b *BufferedLogger) Output(a ...interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	msg := fmt.Sprint(a...)
	fmt.Fprintln(b.writer, msg)
}

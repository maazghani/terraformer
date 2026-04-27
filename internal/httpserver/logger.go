package httpserver

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

// logEntry is the structured JSON log line emitted per request/event.
type logEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// jsonLogger writes newline-delimited JSON log entries to a writer.
type jsonLogger struct {
	enc      *json.Encoder
	mu       sync.Mutex
	minLevel int // 0=debug, 1=info, 2=warn, 3=error
}

// levelToInt converts a log level string to an integer for comparison.
func levelToInt(level string) int {
	switch level {
	case "debug":
		return 0
	case "info":
		return 1
	case "warn":
		return 2
	case "error":
		return 3
	default:
		return 1 // default to info
	}
}

// newJSONLogger returns a jsonLogger writing to w with the given minimum log level.
func newJSONLogger(w io.Writer, minLevelStr string) *jsonLogger {
	return &jsonLogger{
		enc:      json.NewEncoder(w),
		minLevel: levelToInt(minLevelStr),
	}
}

// log emits a single JSON log entry if the level meets the minimum threshold.
func (l *jsonLogger) log(level, message string) {
	if levelToInt(level) < l.minLevel {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.enc.Encode(logEntry{
		Level:   level,
		Message: message,
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}

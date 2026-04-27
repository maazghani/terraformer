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
	enc *json.Encoder
	mu  sync.Mutex
}

// newJSONLogger returns a jsonLogger writing to w.
func newJSONLogger(w io.Writer) *jsonLogger {
	return &jsonLogger{enc: json.NewEncoder(w)}
}

// log emits a single JSON log entry.
func (l *jsonLogger) log(level, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.enc.Encode(logEntry{
		Level:   level,
		Message: message,
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}

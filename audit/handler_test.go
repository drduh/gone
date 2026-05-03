package audit

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"testing"
)

// newTestAuditor creates a new Auditor writing to buf.
func newTestAuditor(buf *bytes.Buffer) *Auditor {
	a := &Auditor{
		Config: Config{
			TimeFormat: "2006-01-02 15:04:05",
		},
		Handler: slog.NewJSONHandler(io.Discard, nil),
		Logger:  log.New(buf, "", 0),
	}
	a.Log = slog.New(a)
	return a
}

// TestHandleMarshalable tests a log entry is written
// as valid JSON with expected fields.
func TestHandleMarshalable(t *testing.T) {
	var buf bytes.Buffer
	a := newTestAuditor(&buf)
	a.Log.Info("test event", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("expected valid json, got error: %v\noutput: %s",
			err, buf.String())
	}

	for _, field := range []string{"time", "level", "message", "data"} {
		if _, ok := entry[field]; !ok {
			t.Errorf("expected %q in log entry, got: %v",
				field, entry)
		}
	}

	if entry["message"] != "test event" {
		t.Errorf("expected message %q, got %q",
			"test event", entry["message"])
	}
}

// TestHandleUnmarshallable tests marshal errors are logged.
func TestHandleUnmarshallable(t *testing.T) {
	var buf bytes.Buffer
	a := newTestAuditor(&buf)
	a.Log.Error("forbidden", "user", func() {})

	output := buf.String()
	if output == "" {
		t.Fatal("expected output for marshal error, got nothing")
	}

	var entry map[string]any
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("expected valid json fallback, got error: %v\noutput: %s",
			err, output)
	}

	if entry["level"] != "ERROR" {
		t.Errorf("expected level ERROR, got %q", entry["level"])
	}
	if entry["event"] != "forbidden" {
		t.Errorf("expected event %q, got %q", "forbidden", entry["event"])
	}
	if _, ok := entry["error"]; !ok {
		t.Errorf("expected 'error' in fallback log entry, got: %v", entry)
	}
	if _, ok := entry["time"]; !ok {
		t.Errorf("expected 'time' in fallback log entry, got: %v", entry)
	}
}

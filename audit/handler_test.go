package audit

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestHandleMarshalable tests a log entry is written
// as valid JSON with expected fields.
func TestHandleMarshalable(t *testing.T) {
	var buf bytes.Buffer
	a := newTestAuditor(&buf)
	a.Log.Info("test event", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(
		buf.Bytes(), &entry); err != nil {
		t.Fatalf("expected json, got error: %v\noutput: %s",
			err, buf.String())
	}

	for _, field := range []string{
		"time", "level", "message", "data",
	} {
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
		t.Fatal("no output for marshal error")
	}

	var entry map[string]any
	if err := json.Unmarshal(
		[]byte(output), &entry); err != nil {
		t.Fatalf("expected json, got error: %v\noutput: %s",
			err, output)
	}

	if entry["level"] != "ERROR" {
		t.Errorf("expected level ERROR, got %q",
			entry["level"])
	}
	if entry["event"] != "forbidden" {
		t.Errorf("expected event %q, got %q",
			"forbidden", entry["event"])
	}
	if _, ok := entry["error"]; !ok {
		t.Errorf("expected 'error' in log, got: %v",
			entry)
	}
	if _, ok := entry["time"]; !ok {
		t.Errorf("expected 'time' in log, got: %v",
			entry)
	}
}

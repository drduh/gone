package storage

import (
	"mime"
	"net/http/httptest"
	"testing"
	"time"
)

// TestFileServeDispositionEscape tests Serve creates a valid
// disposition is set for File names with escape characters.
func TestFileServeDispositionEscape(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"basic", "file.txt"},
		{"space", "my file.txt"},
		{"quote", `foo"bar.txt`},
		{"semicolon", "foo;bar.txt"},
		{"equals", "foo=bar.txt"},
		{"comma", "foo,bar.txt"},
		{"backslash", `foo\bar.txt`},
		{"cr", "foo\rbar.txt"},
		{"lf", "foo\nbar.txt"},
		{"crlf", "foo\r\nbar.txt"},
		{"injection", "foo\r\nAnother-Header: v.txt"},
		{"unicode", "üñîçødé.txt"},
		{"percent", "foo%bar.txt"},
		{"empty", ""},
	}

	for _, tt := range tests {
		f := &File{Name: tt.filename}
		rec := httptest.NewRecorder()
		f.Serve(rec)

		got := rec.Header().Get("Content-Disposition")
		typ, params, err := mime.ParseMediaType(got)
		if err != nil {
			t.Fatalf("%s: invalid header: %v\nheader: %q",
				tt.name, err, got)
		}
		if typ != "attachment" {
			t.Fatalf("%s: invalid type: got %q, want %q",
				tt.name, typ, "attachment")
		}
		if params["filename"] != tt.filename {
			t.Fatalf("%s: invalid filename: got %q, want %q\nheader: %q",
				tt.name, params["filename"], tt.filename, got)
		}
	}
}

// TestServeMessages tests writing Messages to response.
func TestServeMessages(t *testing.T) {
	const timeFormat = "Monday Jan 2 15:04"

	s := &Storage{
		Messages: map[int]*Message{
			1: {
				Count: 1,
				Data:  "hello",
				Time: Time{Allow: time.Date(
					2026, 5, 22, 23, 0, 0, 0, time.UTC).Format(timeFormat)},
			},
			2: {
				Count: 2,
				Data:  "world",
				Time: Time{Allow: time.Date(
					2026, 5, 22, 23, 1, 0, 0, time.UTC).Format(timeFormat)},
			},
		},
	}

	rr := httptest.NewRecorder()
	s.ServeMessages(rr)

	disp := `attachment; filename="messages.txt"`
	if rr.Header().Get("Content-Disposition") != disp {
		t.Errorf("Content-Disposition = %q; want '%q'",
			rr.Header().Get("Content-Disposition"), disp)
	}
	if rr.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("Content-Type = %q; want 'text/plain'",
			rr.Header().Get("Content-Type"))
	}

	body := "" +
		"1 (Friday May 22 23:00) - hello\n" +
		"2 (Friday May 22 23:01) - world\n"
	if rr.Body.String() != body {
		t.Errorf("body = %q; want %q",
			rr.Body.String(), body)
	}
}

// TestServeWall tests writing Wall content to response.
func TestServeWall(t *testing.T) {
	const wallContent = "test wall content"
	s := &Storage{WallContent: wallContent}

	rr := httptest.NewRecorder()
	s.ServeWall(rr)

	disp := `attachment; filename="wall.txt"`
	if rr.Header().Get("Content-Disposition") != disp {
		t.Errorf("Content-Disposition = %q; want '%q'",
			rr.Header().Get("Content-Disposition"), disp)
	}
	if rr.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("Content-Type = %q; want 'text/plain'",
			rr.Header().Get("Content-Type"))
	}
	if rr.Body.String() != wallContent {
		t.Errorf("body = %q; want '%q'",
			rr.Body.String(), wallContent)
	}
}

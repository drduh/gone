package storage

import (
	"mime"
	"net/http/httptest"
	"testing"
)

// TestFileServeDispositionQuote tests Serve creates a valid
// Content-Disposition header for filenames containing quotes.
func TestFileServeDispositionQuote(t *testing.T) {
	filename := `foo"bar.txt`
	f := &File{Name: filename}

	rec := httptest.NewRecorder()
	f.Serve(rec)

	got := rec.Header().Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(got)
	if err != nil {
		t.Fatalf("invalid disposition: %v\nheader: %q", err, got)
	}
	if params["filename"] != filename {
		t.Fatalf("invalid filename: got %q, want %q\nheader: %q",
			params["filename"], filename, got)
	}
}

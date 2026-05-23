package storage

import (
	"mime"
	"net/http/httptest"
	"testing"
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

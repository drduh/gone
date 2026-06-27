package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/util"
)

// TestList tests a successful file list op.
func TestList(t *testing.T) {
	app := newTestApp()
	app.Require.List = false

	data := []byte("hello, world!\n")
	f := &storage.File{
		Name:      "test.txt",
		Data:      data,
		Downloads: storage.Downloads{Allow: 10},
		Time: storage.Time{
			Duration:   5 * time.Minute,
			UploadTime: time.Now(),
		},
	}

	f.Scan()
	f.ID = "1ABCDEF"
	app.Files = map[string]*storage.File{f.ID: f}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.List, nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Fatalf("unexpected Content-Type: %q", ct)
	}

	var files []storage.File
	if err := json.NewDecoder(
		rr.Body).Decode(&files); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if files[0].Remain != 10 {
		t.Errorf("expected %d downloads to remain, got %d",
			10, files[0].Remain)
	}

	if files[0].ID != "1ABCDEF" {
		t.Errorf("expected id %q, got %q",
			"1ABCDEF", files[0].ID)
	}

	if files[0].Name != "test.txt" {
		t.Errorf("expected filename %q, got %q",
			"test.txt", files[0].Name)
	}

	expectedSize := util.FormatSize(len(data))
	if files[0].Size != expectedSize {
		t.Errorf("expected size %q, got %q",
			expectedSize, files[0].Size)
	}

	if files[0].Sum != "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc" {
		t.Errorf("expected sum %q, got %q",
			"4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc",
			files[0].Sum)
	}
}

// TestListDeny tests denied List requests.
func TestListDeny(t *testing.T) {
	app := newTestApp()
	app.Require.List = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.List, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)
}

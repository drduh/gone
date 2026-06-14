package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

const (
	testAddrAndPort = "127.0.0.1:12345"
	testUserAgent   = "testGoneAgent"
	formContentType = "application/x-www-form-urlencoded"
	testContentMsgs = "hello, world!"
	testContentWall = "hello,\r\nworld!\r\n"
)

// newTestApp sets up a configured App for tests,
// ignoring logging and rate limiting.
func newTestApp() *config.App {
	app := config.Load()
	app.Log = slog.New(slog.DiscardHandler)
	app.ReqsPerMinute = 1000
	return app
}

// newTestAppWithStorage sets up a configured
// App with Storage content.
func newTestAppWithStorage() *config.App {
	app := newTestApp()
	app.Storage = storage.Storage{
		Files: map[string]*storage.File{
			"file1": {},
			"file2": {},
		},
		Messages: []*storage.Message{
			{Count: 1, Data: "hello"},
			{Count: 2, Data: "world"},
		},
		WallContent: "test wall content",
	}
	return app
}

// serveDeniedRequest helps serve requests to deny,
// by setting tarpit duration to 0.
func serveDeniedRequest(
	t *testing.T,
	handler http.HandlerFunc,
	req *http.Request,
) *httptest.ResponseRecorder {
	t.Helper()

	auth.SetTarpit(0)

	if req.RemoteAddr == "" {
		req.RemoteAddr = testAddrAndPort
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// assertDenied tests request denial.
func assertDenied(
	t *testing.T,
	rr *httptest.ResponseRecorder,
	want string) {
	t.Helper()

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d",
			http.StatusForbidden, rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(
		rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if got := body["error"]; got != want {
		t.Fatalf("expected error %q, got %q",
			want, got)
	}
}

// assertStorageClear tests Storage is empty.
func assertStorageClear(t *testing.T, app *config.App) {
	t.Helper()
	assertFilesClear(t, app)
	assertMessagesClear(t, app)
	assertWallClear(t, app)
}

// assertFilesClear tests Files is empty.
func assertFilesClear(t *testing.T, app *config.App) {
	t.Helper()

	if app.Files == nil {
		t.Fatalf("Files is nil; want empty map")
	}
	if got := len(app.Files); got != 0 {
		t.Fatalf("Files length = %d; want 0", got)
	}
}

// assertMessagesClear tests Messages is empty.
func assertMessagesClear(t *testing.T, app *config.App) {
	t.Helper()

	if app.Messages == nil {
		t.Fatalf("Messages is nil; want empty slice")
	}
	if got := len(app.Messages); got != 0 {
		t.Fatalf("Messages length = %d; want 0", got)
	}
}

// assertWallClear tests Wall is empty.
func assertWallClear(t *testing.T, app *config.App) {
	t.Helper()

	if app.WallContent != "" {
		t.Fatalf("WallContent = %q; want empty string",
			app.WallContent)
	}
}

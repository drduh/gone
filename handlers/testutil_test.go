package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drduh/gone/auth"
	"github.com/drduh/gone/config"
	"github.com/drduh/gone/storage"
)

const (
	formContentType = "application/x-www-form-urlencoded"
	testAddrAndPort = "127.0.0.1:12345"
	testContentMsgs = "hello, world! (see http://example.com)"
	testContentWall = "hello, world!\r\n(see http://example.com)\r\n"
	testUserAgent   = "test goneAgent-v1.2026"
)

// newTestApp sets up a configured App for tests,
// ignoring logging and rate limiting.
func newTestApp() *config.App {
	app := config.Load()

	app.Hostname = "testRunHost"
	app.Log = slog.New(slog.DiscardHandler)
	app.ReqsPerMinute = 99
	app.StartTime = time.Now().Add(-time.Second)

	app.MessageLimits.MaxCount = 32
	app.MessageLimits.LengthChars = 128
	app.WallLimits.LengthChars = 1024
	app.FileLimits.NameLength = 64

	auth.SetTarpit(0)

	return app
}

// newTestMux sets up a route handlers for tests.
func newTestMux(app *config.App) *http.ServeMux {
	mux := http.NewServeMux()

	for pattern, h := range Routes(app) {
		mux.HandleFunc(pattern, h)
	}

	return mux
}

// newTestAppWithStorage sets up a configured
// App with Storage content.
func newTestAppWithStorage() *config.App {
	app := newTestApp()

	app.Storage = storage.Storage{
		Files: map[string]*storage.File{
			"file1": {
				Name:  "file1",
				ID:    "12345",
				Bytes: 10,
			},
			"file2": {
				Name:  "file2",
				ID:    "67890",
				Bytes: 20,
			},
		},
		Messages: []*storage.Message{
			{Count: 1, Data: testContentMsgs + "1"},
			{Count: 2, Data: testContentMsgs + "2"},
		},
		WallContent: testContentWall,
	}

	return app
}

// serveRequest serves req through the app handler mux.
func serveRequest(
	t *testing.T,
	app *config.App,
	req *http.Request,
) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	newTestMux(app).ServeHTTP(rr, req)

	return rr
}

// serveDeniedRequest serves a request expected to be denied.
func serveDeniedRequest(
	t *testing.T,
	app *config.App,
	req *http.Request,
) *httptest.ResponseRecorder {
	t.Helper()

	auth.SetTarpit(0)

	if req.RemoteAddr == "" {
		req.RemoteAddr = testAddrAndPort
	}

	return serveRequest(t, app, req)
}

// assertDenied tests request denial.
func assertDenied(
	t *testing.T,
	rr *httptest.ResponseRecorder,
	want string) {
	t.Helper()

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d",
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

// assertStatus checks the response status code.
func assertStatus(
	t *testing.T,
	rr *httptest.ResponseRecorder,
	want int,
) {
	t.Helper()

	if got := rr.Code; got != want {
		t.Fatalf("status = %d; want %d", got, want)
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

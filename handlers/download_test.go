package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDownloadServe tests successful download.
func TestDownloadServe(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Download = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download+"?name=file1", nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}
}

// TestDownloadMissingFilename tests requests
// without a filename provided.
func TestDownloadMissingFilename(t *testing.T) {
	app := newTestApp()
	app.Require.Download = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download, nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d",
			http.StatusNotFound, rr.Code)
	}
}

// TestDownloadFileNotFound tests requests to
// download a File not found in Storage.
func TestDownloadFileNotFound(t *testing.T) {
	app := newTestApp()
	app.Require.Download = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download+"?name=none.txt", nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d",
			http.StatusNotFound, rr.Code)
	}
}

// TestDownloadDeny tests denied requests to
// download a File.
func TestDownloadDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Download = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)
}

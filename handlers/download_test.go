package handlers

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestDownloadServe tests successful download.
func TestDownloadServe(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Download = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download+"?name=file1", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)
}

// TestDownloadMissingFilename tests requests
// without a filename provided.
func TestDownloadMissingFilename(t *testing.T) {
	app := newTestApp()
	app.Require.Download = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download, nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusNotFound)
}

// TestDownloadFileNotFound tests requests to
// download a File not found in Storage.
func TestDownloadFileNotFound(t *testing.T) {
	app := newTestApp()
	app.Require.Download = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download+"?name=none.txt", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusNotFound)
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

// TestDownloadLastExpires tests a File is removed once
// its download limit is met, counted before serving.
func TestDownloadLastExpires(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Download = false
	app.Files["12345"].Allow = 1

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download+"12345", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	if _, ok := app.Files["12345"]; ok {
		t.Fatal("expected file1 removed after last download")
	}

	rr = serveRequest(t, app, httptest.NewRequestWithContext(
		t.Context(), http.MethodGet, app.Download+"12345", nil))
	assertStatus(t, rr, http.StatusNotFound)
}

// TestDownloadConcurrentLimit tests only one of many
// concurrent requests gets the last allowed download.
func TestDownloadConcurrentLimit(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Download = false
	app.Files["12345"].Allow = 1

	const n = 20
	codes := make(chan int, n)
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			req := httptest.NewRequestWithContext(t.Context(),
				http.MethodGet, app.Download+"12345", nil)
			codes <- serveRequest(t, app, req).Code
		})
	}
	wg.Wait()
	close(codes)

	served := 0
	for code := range codes {
		if code == http.StatusOK {
			served++
		}
	}
	if served != 1 {
		t.Fatalf("served %d downloads; want 1", served)
	}
}

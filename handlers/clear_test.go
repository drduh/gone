package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClearBrowser tests clear Storage and redirect.
func TestClearBrowser(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Clear = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Clear, nil)
	req.RemoteAddr = testAddrAndPort
	req.Header.Set("Accept", "text/html")

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected %d, got %d",
			http.StatusSeeOther, rr.Code)
	}

	if got := rr.Header().Get("Location"); got != app.Root {
		t.Fatalf("Location = %q; want %q", got, app.Root)
	}

	assertStorageClear(t, app)
}

// TestClearJSON tests clear Storage and JSON response.
func TestClearJSON(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Clear = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Clear, nil)
	req.RemoteAddr = testAddrAndPort
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusSeeOther, rr.Code)
	}

	want := "\"storage cleared\"\n"
	if got := rr.Body.String(); got != want {
		t.Fatalf("body = %q; want %q", got, want)
	}

	assertStorageClear(t, app)
}

// TestClearDeny tests denied Clear requests.
func TestClearDeny(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Clear = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Clear, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)

	if len(app.Files) != 2 {
		t.Fatalf("expected Files unchanged, got %d",
			len(app.Files))
	}
	if len(app.Messages) != 2 {
		t.Fatalf("expected Messages unchanged, got %d",
			len(app.Messages))
	}
	if app.WallContent == "" {
		t.Fatal("expected WallContent unchanged")
	}
}

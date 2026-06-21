package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drduh/gone/templates"
)

// TestStatic tests serving embedded static content.
func TestStatic(t *testing.T) {
	app := newTestApp()
	app.Require.Static = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Static, nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got !=
		"application/json; charset=utf-8" {
		t.Fatalf("unexpected Content-Type: %q", got)
	}

	var resp templates.Static
	if err := json.NewDecoder(
		rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.Data != templates.StaticData {
		t.Fatalf("unexpected static data")
	}
}

// TestStaticDeny tests denied Static requests.
func TestStaticDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Static = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Static, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)
}

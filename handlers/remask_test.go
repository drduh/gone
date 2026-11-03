package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUserRemask tests a new mask is assigned and reported.
func TestUserRemask(t *testing.T) {
	app := newTestApp()
	app.Require.UserRemask = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.UserRemask, nil)
	req.RemoteAddr = testAddrAndPort
	before := getMask(req.RemoteAddr)

	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	after := getMask(req.RemoteAddr)
	if before == after {
		t.Fatalf("mask unchanged: %q", after)
	}
	if !strings.Contains(rr.Body.String(),
		before+" is now "+after) {
		t.Fatalf("body = %q", rr.Body.String())
	}

	req = httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.UserRemask, nil)
	req.RemoteAddr = testAddrAndPort
	req.Header.Set("Accept", "text/html")
	rr = serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusSeeOther)
}

// TestUndefined tests unknown paths get a 404.
func TestUndefined(t *testing.T) {
	app := newTestApp()

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/no/such/path", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusNotFound)

	if !strings.Contains(rr.Body.String(), app.NoPath) {
		t.Fatalf("body = %q; want %q", rr.Body.String(), app.NoPath)
	}
}

// TestRateLimit tests requests over the per-client
// limit are rejected, while another client passes.
func TestRateLimit(t *testing.T) {
	app := newTestApp()
	app.Require.Status = false
	app.ReqsPerMinute = 1

	get := func(addr string) *httptest.ResponseRecorder {
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodGet, app.Status, nil)
		req.RemoteAddr = addr
		return serveRequest(t, app, req)
	}

	assertStatus(t, get("10.0.0.1:1"), http.StatusOK)
	rr := get("10.0.0.1:2")
	assertStatus(t, rr, http.StatusTooManyRequests)
	if !strings.Contains(rr.Body.String(), app.RateLimit) {
		t.Fatalf("body = %q", rr.Body.String())
	}
	assertStatus(t, get("10.0.0.2:1"), http.StatusOK)
}

// TestClientAddr tests host extraction with and without a port.
func TestClientAddr(t *testing.T) {
	cases := map[string]string{
		"127.0.0.1:12345": "127.0.0.1",
		"[::1]:80":        "::1",
		"unix-socket":     "unix-socket",
	}

	for in, want := range cases {
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodGet, "/", nil)
		req.RemoteAddr = in
		if got := clientAddr(req); got != want {
			t.Fatalf("clientAddr(%q) = %q; want %q", in, got, want)
		}
	}
}

// TestIsAuthenticatedNoPolicy tests an unregistered pattern.
func TestIsAuthenticatedNoPolicy(t *testing.T) {
	app := newTestApp()

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/nope", nil)
	req.Pattern = "/nope"
	if isAuthenticated(app, req) {
		t.Fatal("expected deny with no policy")
	}
}

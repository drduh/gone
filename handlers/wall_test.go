package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestWallGet tests reading wall content.
func TestWallGet(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false
	app.WallContent = testContentWall

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Wall, nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}

	var got string
	if err := json.NewDecoder(
		rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got != testContentWall {
		t.Errorf("expected wall content %q, got %q",
			testContentWall, got)
	}
}

// TestWallPostUpdate tests updating wall content.
func TestWallPostUpdate(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false
	values := "wall=new content"

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Wall, strings.NewReader(values))
	req.Header.Set("Content-Type", formContentType)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}

	var got string
	if err := json.NewDecoder(
		rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode wall response: %v", err)
	}
	if got != "new content" {
		t.Errorf("expected wall content %q, got %q",
			"new content", got)
	}
}

// TestWallPostClear tests clearing wall content.
func TestWallPostClear(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false
	app.WallContent = testContentWall

	values := formFieldClear + "=1"
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Wall, strings.NewReader(values))
	req.Header.Set("Content-Type", formContentType)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}

	var got string
	if err := json.NewDecoder(
		rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode wall response: %v", err)
	}
	if got != "" {
		t.Errorf("expected wall content cleared, got %q", got)
	}

	assertWallClear(t, app)
}

// TestWallGetDownloadAll tests downloading wall content.
func TestWallGetDownloadAll(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false
	app.WallContent = testContentWall

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Wall+"?download=wall", nil)
	req.RemoteAddr = testAddrAndPort

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}

	disp := rr.Header().Get("Content-Disposition")
	if disp != `attachment; filename="wall.txt"` {
		t.Fatalf("invalid Content-Disposition: %q", disp)
	}

	body := rr.Body.String()
	if body != testContentWall {
		t.Fatalf("expected wall content %q, got %q",
			"downloaded wall content", body)
	}
}

// TestWallDeny tests denied Wall requests.
func TestWallDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = true
	app.WallContent = testContentWall

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost,
		app.Wall, strings.NewReader("wall=new content"))
	req.Header.Set("Content-Type", formContentType)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)

	if app.WallContent != testContentWall {
		t.Fatalf("expected wall content unchanged, got %q",
			app.WallContent)
	}
}

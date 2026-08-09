package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestWallGet tests reading wall content.
func TestWallGet(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false
	app.WallContent = testContentWall

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.WallModify, nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

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
		http.MethodPost, app.WallModify,
		strings.NewReader(values))
	req.Header.Set("Content-Type", formContentType)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	var got string
	if err := json.NewDecoder(
		rr.Body).Decode(&got); err != nil {
		t.Errorf("failed to decode wall response: %v", err)
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
		http.MethodPost, app.WallModify,
		strings.NewReader(values))
	req.Header.Set("Content-Type", formContentType)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

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

// TestWallPostDownload tests downloading wall content.
func TestWallPostDownload(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false
	app.WallContent = testContentWall

	form := url.Values{}
	form.Set(formFieldDownload, formFieldWall)

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.WallModify,
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", formContentType)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	got := rr.Header().Get("Content-Disposition")
	want := `attachment; filename="wall.txt"`
	if got != want {
		t.Fatalf("invalid Content-Disposition: got %q, want %q",
			got, want)
	}

	if got := rr.Body.String(); got != testContentWall {
		t.Fatalf("expected wall content %q, got %q",
			testContentWall, got)
	}
}

// TestWallDeny tests denied Wall requests.
func TestWallDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = true
	app.WallContent = testContentWall

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.WallModify,
		strings.NewReader("wall=new content"))
	req.Header.Set("Content-Type", formContentType)
	rr := serveDeniedRequest(t, app, req)
	assertDenied(t, rr, app.Deny)

	if app.WallContent != testContentWall {
		t.Fatalf("expected wall content unchanged, got %q",
			app.WallContent)
	}
}

package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestGetDefaultTheme tests fixed and time-based themes.
func TestGetDefaultTheme(t *testing.T) {
	if got := getDefaultTheme("lime"); got != "lime" {
		t.Fatalf("fixed theme = %q", got)
	}

	if got := getDefaultTheme(autoTheme); got != "light" &&
		got != "dark" {
		t.Fatalf("auto theme = %q", got)
	}
}

// TestGetTheme tests form selection, fallback, and cookie.
func TestGetTheme(t *testing.T) {
	themes := []string{"light", "dark", "lime"}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/?theme=lime", nil)
	rr := httptest.NewRecorder()
	if got := getTheme(rr, req, "light", "id", time.Hour,
		themes); got != "lime" {
		t.Fatalf("selected theme = %q", got)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Value != "lime" {
		t.Fatalf("cookie = %v", cookies)
	}

	req = httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/?theme=nope", nil)
	rr = httptest.NewRecorder()
	if got := getTheme(rr, req, "lime", "id", time.Hour,
		themes); got != "light" && got != "dark" {
		t.Fatalf("invalid theme fell back to %q", got)
	}

	req = httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "id", Value: "dark"})
	if got := getTheme(httptest.NewRecorder(), req, "light",
		"id", time.Hour, themes); got != "dark" {
		t.Fatalf("cookie theme = %q", got)
	}
}

// TestIndexOptions tests index rendering with CSP and
// theme picking disabled, and with a theme selected.
func TestIndexOptions(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Root = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/?theme=dark", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)
	if rr.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("expected CSP header")
	}
	if !strings.Contains(rr.Body.String(), "file1") {
		t.Fatal("expected file list in index")
	}

	app.CSP = false
	app.Style.AllowPick = false
	app.Style.Theme = "lime"
	rr = serveRequest(t, app, httptest.NewRequestWithContext(
		t.Context(), http.MethodGet, "/", nil))
	assertStatus(t, rr, http.StatusOK)
	if rr.Header().Get("Content-Security-Policy") != "" {
		t.Fatal("expected no CSP header")
	}
	if len(rr.Result().Cookies()) != 0 {
		t.Fatal("expected no cookie with picking disabled")
	}
}

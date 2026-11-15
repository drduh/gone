package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
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

// TestThemeRules tests each default theme has one stylesheet rule.
func TestThemeRules(t *testing.T) {
	css, err := os.ReadFile("../assets/style.css")
	if err != nil {
		t.Fatal(err)
	}

	themes := newTestApp().Style.Available
	if len(themes) == 0 {
		t.Fatal("expected default themes")
	}

	for _, name := range themes {
		rule := `[data-theme="` + name + `"]`
		if got := strings.Count(string(css), rule); got != 1 {
			t.Errorf("%s rules = %d; want 1", rule, got)
		}
	}
}

// TestIndexTheme tests the index names and marks the current theme.
func TestIndexTheme(t *testing.T) {
	app := newTestApp()
	app.Require.Root = false

	rr := serveRequest(t, app, httptest.NewRequestWithContext(
		t.Context(), http.MethodGet, "/?theme=lime", nil))
	assertStatus(t, rr, http.StatusOK)

	body := rr.Body.String()
	if !strings.Contains(body,
		`<html lang="en" data-theme="lime">`) {
		t.Error("expected data-theme on html")
	}
	if strings.Contains(body, "<style") {
		t.Error("expected no style element")
	}
	if got := strings.Count(body, `aria-current="true"`); got != 1 {
		t.Errorf("current swatches = %d; want 1", got)
	}
	if !strings.Contains(body,
		`aria-label="lime" aria-current="true"`) {
		t.Error("expected the lime swatch marked current")
	}
}

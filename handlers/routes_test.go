package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drduh/gone/auth"
)

// TestRegisterDisabled tests an empty path is not routed
// and every other path denies without auth.
func TestRegisterDisabled(t *testing.T) {
	app := newTestApp()
	app.Random = ""

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/random/coin", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusNotFound)

	app.Require.MessageGet = true
	app.Require.MessageDelete = true
	app.Require.UserRemask = true
	for _, p := range []string{
		app.Message + "/1",
		app.Message + "/delete/1",
		app.UserRemask,
	} {
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodPost, p, nil)
		rr := serveDeniedRequest(t, app, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("%s: status = %d; want 403", p, rr.Code)
		}
	}
}

// TestStateChangeRequiresPost tests GET cannot modify state.
func TestStateChangeRequiresPost(t *testing.T) {
	app := newTestAppWithStorage()
	auth.SetTarpit(0)

	for _, p := range []string{
		app.MessageClear,
		app.Message + "/delete/1",
		app.Clear,
		app.UserRemask,
		app.Upload,
		app.MessageAdd + "?message=x",
	} {
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodGet, p, nil)
		req.Header.Set(app.Basic.Field, app.Basic.Token)
		rr := serveRequest(t, app, req)
		if rr.Code == http.StatusOK ||
			rr.Code == http.StatusSeeOther {
			t.Fatalf("GET %s: status = %d; want rejection",
				p, rr.Code)
		}
	}

	if len(app.Messages) != 2 || len(app.Files) != 2 ||
		app.WallContent == "" {
		t.Fatal("expected storage unchanged by GET")
	}
}

// TestPatternPath tests method prefix and suffix removal.
func TestPatternPath(t *testing.T) {
	cases := map[string]string{
		"POST /msg/clear": "/msg/clear",
		"/{$}":            "/",
		"/download/":      "/download/",
		"":                "",
	}

	for in, want := range cases {
		if got := patternPath(in); got != want {
			t.Fatalf("patternPath(%q) = %q; want %q", in, got, want)
		}
	}

	if post("") != "" || post("/x") != "POST /x" {
		t.Fatal("post() prefix incorrect")
	}
}

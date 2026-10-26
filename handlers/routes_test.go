package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

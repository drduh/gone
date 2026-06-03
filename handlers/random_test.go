package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRandomDeny tests denied Random requests.
func TestRandomDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Random = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Random+"test", nil)
	rr := serveDeniedRequest(t, Random(app), req)

	assertDenied(t, rr, app.Deny)
}

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestStatusDeny tests denied Status requests.
func TestStatusDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Status = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Status, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)
}

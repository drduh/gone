package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIndexDeny tests denied Index requests.
func TestIndexDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Root = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Root, nil)
	rr := serveDeniedRequest(t, Index(app), req)

	assertDenied(t, rr, app.Deny)
}

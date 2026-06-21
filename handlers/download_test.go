package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDownloadDeny tests denied Download requests.
func TestDownloadDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Download = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Download, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)
}

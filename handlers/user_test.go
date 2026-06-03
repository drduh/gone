package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUserInfoDeny tests Denied UserInfo requests.
func TestUserInfoForbidden(t *testing.T) {
	app := newTestApp()
	app.Require.UserInfo = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.UserInfo, nil)
	rr := serveDeniedRequest(t, UserInfo(app), req)

	assertDenied(t, rr, app.Deny)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drduh/gone/templates"
)

// TestUserInfoDeny tests denied UserInfo requests.
func TestUserInfoForbidden(t *testing.T) {
	app := newTestApp()
	app.Require.UserInfo = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.UserInfo, nil)
	rr := serveDeniedRequest(t, UserInfo(app), req)

	assertDenied(t, rr, app.Deny)
}

// TestUserInfoJSON tests JSON-based UserInfo requests.
func TestUserInfoJSON(t *testing.T) {
	app := newTestApp()

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.UserInfo, nil)
	req.RemoteAddr = testAddrAndPort
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", testUserAgent)
	rr := httptest.NewRecorder()

	UserInfo(app).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d",
			rr.Code, http.StatusOK)
	}

	var got templates.User
	if err := json.NewDecoder(
		rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Address != testAddrAndPort {
		t.Fatalf("user address = %q; want %q",
			got.Address, testAddrAndPort)
	}

	if got.Mask == "" {
		t.Fatal("expected address mask")
	}

	if got.IsBrowser {
		t.Fatal("expected isBrowser to be false")
	}

	if ua := got.Headers.Get(
		"User-Agent"); ua != testUserAgent {
		t.Fatalf("User-Agent = %q; want %q",
			ua, testUserAgent)
	}
}

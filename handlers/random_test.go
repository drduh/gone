package handlers

import (
	"encoding/json"
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

// TestRandomCoin tests requests for random coin values.
func TestRandomCoin(t *testing.T) {
	app := newTestApp()
	app.RandomLimits.StrCount = 3

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Random+"coin", nil)
	rr := httptest.NewRecorder()

	Random(app).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rr.Code)
	}

	var got []string
	if err := json.NewDecoder(
		rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 strings, got %d", len(got))
	}

	for _, v := range got {
		if v != "heads" && v != "tails" {
			t.Fatalf("expected coin result, got %q", v)
		}
	}
}

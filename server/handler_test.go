package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drduh/gone/auth"
)

// TestGetHandlerDisabledPaths tests empty/disabled paths
// are skipped, not registered.
func TestGetHandlerDisabledPaths(t *testing.T) {
	app := newTestApp(nil)
	auth.SetTarpit(0)
	app.Message = ""
	app.WallModify = ""

	h := getHandler(app)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(
		http.MethodGet, "/wall", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("disabled path status = %d, want %d",
			rec.Code, http.StatusNotFound)
	}
}

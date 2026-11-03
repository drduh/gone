package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/wall", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("disabled path status = %d, want %d",
			rec.Code, http.StatusNotFound)
	}
}

// TestCrossOriginDenied tests a cross-site browser POST
// is rejected while same-origin and non-browser pass.
func TestCrossOriginDenied(t *testing.T) {
	app := newTestApp(nil)
	auth.SetTarpit(0)
	app.ReqsPerMinute = 99
	h := getHandler(app)

	cases := []struct {
		name string
		site string
		want int
	}{
		{"cross-site", "cross-site", http.StatusForbidden},
		{"same-origin", "same-origin", http.StatusSeeOther},
		{"non-browser", "", http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(),
				http.MethodPost, "/msg/add", strings.NewReader("message=hi"))
			req.Header.Set("Content-Type",
				"application/x-www-form-urlencoded")
			if tc.site != "" {
				req.Header.Set("Sec-Fetch-Site", tc.site)
				req.Header.Set("Accept", "text/html")
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d",
					rec.Code, tc.want)
			}
		})
	}
}

// TestRegisterAssets tests assets are served or skipped.
func TestRegisterAssets(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "assets"),
		0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets",
		"style.css"), []byte("body{}"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Chdir(dir)

	app := newTestApp(nil)
	auth.SetTarpit(0)
	app.ReqsPerMinute = 99
	h := getHandler(app)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(
		http.MethodGet, "/assets/style.css", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "body{}" {
		t.Fatalf("assets: %d %q", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q", got)
	}

	app.Require.Assets = true
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(
		http.MethodGet, "/assets/style.css", nil))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("assets denied status = %d", rec.Code)
	}

	app.Assets = ""
	rec = httptest.NewRecorder()
	getHandler(app).ServeHTTP(rec, httptest.NewRequest(
		http.MethodGet, "/assets/style.css", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("assets disabled status = %d", rec.Code)
	}
}

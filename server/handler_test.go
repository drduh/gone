package server

import (
	"net/http"
	"net/http/httptest"
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
	h.ServeHTTP(rec, httptest.NewRequest(
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
			req := httptest.NewRequest(http.MethodPost,
				"/msg/add", strings.NewReader("message=hi"))
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

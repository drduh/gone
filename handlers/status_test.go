package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drduh/gone/storage"
	"github.com/drduh/gone/version"
)

// TestStatus tests successful Status responses.
func TestStatus(t *testing.T) {
	tests := []struct {
		name      string
		showBuild bool
		wantBuild map[string]string
	}{
		{
			name:      "redacted status",
			showBuild: false,
			wantBuild: version.Get(version.Redacted),
		},
		{
			name:      "full status",
			showBuild: true,
			wantBuild: version.Get(version.Full),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestAppWithStorage()
			app.ShowBuild = tt.showBuild

			req := httptest.NewRequestWithContext(
				t.Context(), http.MethodGet, app.Status, nil)
			rr := serveRequest(t, app, req)
			assertStatus(t, rr, http.StatusOK)

			if got := rr.Header().Get("Content-Type"); got !=
				"application/json; charset=utf-8" {
				t.Fatalf("Content-Type = %q; want JSON", got)
			}

			var response map[string]json.RawMessage
			if err := json.NewDecoder(
				rr.Body).Decode(&response); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			var gotBuild map[string]string
			if err := json.Unmarshal(
				response["buildInfo"], &gotBuild); err != nil {
				t.Fatalf("decode buildInfo: %v", err)
			}

			if !reflect.DeepEqual(gotBuild, tt.wantBuild) {
				t.Errorf("buildInfo = %#v; want %#v",
					gotBuild, tt.wantBuild)
			}

			var gotAddr string
			if err := json.Unmarshal(
				response["serverAddr"], &gotAddr); err != nil {
				t.Fatalf("addr unavailable: %v", err)
			}
			if gotAddr != app.ServerAddr {
				t.Errorf("addr = %q; want %q",
					gotAddr, app.ServerAddr)
			}

			var gotPort int
			if err := json.Unmarshal(
				response["serverPort"], &gotPort); err != nil {
				t.Fatalf("port unavailable: %v", err)
			}
			if gotPort != app.ServerPort {
				t.Errorf("port = %d; want %d",
					gotPort, app.ServerPort)
			}

			var gotHostname string
			if err := json.Unmarshal(
				response["hostname"], &gotHostname); err != nil {
				t.Fatalf("decode hostname: %v", err)
			}
			if gotHostname != app.Hostname {
				t.Errorf("hostname = %q; want %q",
					gotHostname, app.Hostname)
			}

			var gotUptime string
			if err := json.Unmarshal(
				response["uptime"], &gotUptime); err != nil {
				t.Fatalf("decode uptime: %v", err)
			}
			if gotUptime == "" || gotUptime <= "0s" {
				t.Errorf("uptime = %q; want >0s duration",
					gotUptime)
			}

			var gotSizes storage.Sizes
			if err := json.Unmarshal(
				response["storageSizes"], &gotSizes); err != nil {
				t.Fatalf("decode storageSizes: %v", err)
			}

			if gotSizes.NumFiles != 2 {
				t.Errorf("numFiles = %d; want 2",
					gotSizes.NumFiles)
			}
			if gotSizes.SizeFiles != 30 {
				t.Errorf("sizeFiles = %d; want 30",
					gotSizes.SizeFiles)
			}
			if gotSizes.NumMessages != 2 {
				t.Errorf("numMessages = %d; want 2",
					gotSizes.NumMessages)
			}

			var gotWall storage.WallMeta
			if err := json.Unmarshal(
				response["wallMeta"], &gotWall); err != nil {
				t.Fatalf("decode wall meta: %v", err)
			}

			if gotWall.WallChars != len(app.WallContent) {
				t.Errorf("wallChars = %d; want %d",
					gotWall.WallChars, len(app.WallContent))
			}
		})
	}
}

// TestStatusDeny tests denied Status requests.
func TestStatusDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Status = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Status, nil)
	rr := serveDeniedRequest(t, app, req)
	assertDenied(t, rr, app.Deny)
}

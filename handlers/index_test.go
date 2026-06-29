package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/drduh/gone/storage"
)

// TestIndexDeny tests denied Index requests.
func TestIndexDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Root = true

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Root, nil)
	rr := serveDeniedRequest(t, app, req)

	assertDenied(t, rr, app.Deny)
}

// TestIndexMessageEscape tests Message URL encoding.
func TestIndexMessageEscape(t *testing.T) {
	app := newTestApp()
	app.Require.Root = false

	app.Messages = append(app.Messages, &storage.Message{
		Count: 1,
		Data:  `<script>alert("xss")</script> https://example.com`,
		Owner: storage.Owner{
			Mask: testAddrAndPort,
		},
		Time: storage.Time{
			UploadTimeFmt: "now",
		},
	})

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Root, nil)
	req.Header.Set("Accept", "text/html")

	rr := httptest.NewRecorder()
	mux := newTestMux(app)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",
			http.StatusOK, rr.Code)
	}

	body := rr.Body.String()

	if strings.Contains(body,
		`<script>alert("xss")</script>`) {
		t.Fatalf("expected escaped script tag, got: %q", body)
	}

	if !strings.Contains(body,
		`&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;`) {
		t.Fatalf("expected escaped script tag, got: %q", body)
	}

	if !strings.Contains(body,
		`<a href="https://example.com"`) {
		t.Fatalf("expected rendered link, got body: %q", body)
	}
}

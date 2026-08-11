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
	})

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Root, nil)
	req.Header.Set("Accept", "text/html")
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

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

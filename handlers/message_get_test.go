package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drduh/gone/storage"
)

// TestMessageGet tests reading one Message by count.
func TestMessageGet(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.MessageGet = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Message+"/2", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	var got storage.Message
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Count != 2 || got.Data != testContentMsgs+"2" {
		t.Fatalf("got %+v; want count 2", got)
	}
}

// TestMessageGetErrors tests bad and unknown counts.
func TestMessageGetErrors(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.MessageGet = false

	cases := []struct {
		path string
		want int
	}{
		{"/0", http.StatusBadRequest},
		{"/abc", http.StatusBadRequest},
		{"/9", http.StatusNotFound},
	}

	for _, tc := range cases {
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodGet, app.Message+tc.path, nil)
		rr := serveRequest(t, app, req)
		if rr.Code != tc.want {
			t.Fatalf("%s: status = %d; want %d",
				tc.path, rr.Code, tc.want)
		}
	}
}

// TestMessageDelete tests deleting a Message by count,
// renumbering, JSON and browser responses.
func TestMessageDelete(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.MessageDelete = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message+"/delete/1", nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	if len(app.Messages) != 1 || app.Messages[0].Count != 1 ||
		app.Messages[0].Data != testContentMsgs+"2" {
		t.Fatalf("messages = %+v; want renumbered 2nd",
			app.Messages)
	}

	req = httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message+"/delete/1", nil)
	req.Header.Set("Accept", "text/html")
	rr = serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusSeeOther)
	if len(app.Messages) != 0 {
		t.Fatal("expected no messages left")
	}
}

// TestMessageDeleteErrors tests bad and unknown counts.
func TestMessageDeleteErrors(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.MessageDelete = false

	cases := []struct {
		path string
		want int
	}{
		{"/delete/x", http.StatusBadRequest},
		{"/delete/7", http.StatusNotFound},
	}

	for _, tc := range cases {
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodPost, app.Message+tc.path, nil)
		rr := serveRequest(t, app, req)
		if rr.Code != tc.want {
			t.Fatalf("%s: status = %d; want %d",
				tc.path, rr.Code, tc.want)
		}
	}
	if len(app.Messages) != 2 {
		t.Fatal("expected messages unchanged")
	}
}

// TestMessageJSON tests reading all Messages as JSON.
func TestMessageJSON(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Message = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, app.Message, nil)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	var got []storage.Message
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d messages; want 2", len(got))
	}
}

// TestMessageClearBrowser tests clear redirects browsers.
func TestMessageClearBrowser(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.MessageClear = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.MessageClear, nil)
	req.Header.Set("Accept", "text/html")
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusSeeOther)
	if len(app.Messages) != 0 {
		t.Fatal("expected messages cleared")
	}
}

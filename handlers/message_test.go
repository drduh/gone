package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/drduh/gone/storage"
)

// TestMessageHandlerValid tests a valid message post.
func TestMessageHandlerValid(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", testContentMsgs)
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()

	handler := Message(app)
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	if len(app.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d",
			len(app.Messages))
	}

	msg := app.Messages[0]
	if msg == nil {
		t.Fatal("message not found")
	}

	if msg.Data != testContentMsgs {
		t.Errorf("expected message %q, got %q",
			testContentMsgs, msg.Data)
	}
}

// TestMessageHandlerExceedLength tests message
// exceeding configured length causes an error.
func TestMessageHandlerExceedLength(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", strings.Repeat("a", 1000))
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()

	handler := Message(app)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	want := `{"error":"` + app.MsgLength + `"}` + "\n"
	if rr.Body.String() != want {
		t.Errorf("expected response %q, got %q",
			want, rr.Body.String())
	}
}

// TestMessageHandlerClear tests messages are cleared.
func TestMessageHandlerClear(t *testing.T) {
	app := newTestApp()

	app.Messages = append(app.Messages, &storage.Message{
		Count: 1, Data: testContentMsgs,
	})

	form := url.Values{}
	form.Set("clear", "true")
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()

	handler := Message(app)
	handler.ServeHTTP(rr, req)

	if len(app.Messages) != 0 {
		t.Errorf("expected messages cleared, got %d",
			len(app.Messages))
	}
}

// TestMessageHandlerDownloadAll test all messages download.
func TestMessageHandlerDownloadAll(t *testing.T) {
	app := newTestApp()

	app.Messages = append(app.Messages, &storage.Message{
		Count: 1, Data: testContentMsgs + "1",
	})
	app.Messages = append(app.Messages, &storage.Message{
		Count: 2, Data: testContentMsgs + "2",
	})
	app.Messages = append(app.Messages, &storage.Message{
		Count: 3, Data: testContentMsgs + "3",
	})

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/?download=all", nil)

	rr := httptest.NewRecorder()

	handler := Message(app)
	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	for i := 1; i <= 3; i++ {
		want := testContentMsgs + strconv.Itoa(i)
		if !strings.Contains(body, want) {
			t.Errorf("expected message %q, got %q",
				want, body)
		}
	}

	disp := rr.Header().Get("Content-Disposition")
	if disp != `attachment; filename="messages.txt"` {
		t.Errorf("invalid Content-Disposition: %q", disp)
	}
}

// TestMessageHandlerExceedCount tests message
// exceeding configured count causes an error.
func TestMessageHandlerExceedCount(t *testing.T) {
	app := newTestApp()

	handler := Message(app)
	form := url.Values{}

	for i := range app.MessageLimits.MaxCount {
		form.Set("message", fmt.Sprintf("msg %d", i+1))
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodPost, "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", formContentType)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		if rr.Code != 200 {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	}

	form.Set("message", testContentMsgs)
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	want := `{"error":"` + app.MsgCount + `"}` + "\n"
	if rr.Body.String() != want {
		t.Errorf("expected %q, got %q",
			want, rr.Body.String())
	}
}

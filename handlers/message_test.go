package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/drduh/gone/storage"
)

const (
	contentType    = "application/x-www-form-urlencoded"
	messageContent = "hello, world!"
)

// TestMessageHandlerValid tests a valid message post.
func TestMessageHandlerValid(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", messageContent)
	req := httptest.NewRequest(
		"POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)

	rr := httptest.NewRecorder()
	handler := Message(app)
	handler.ServeHTTP(rr, req)

	if rr.Code < 200 || rr.Code >= 400 {
		t.Errorf("expected 2xx, got %d", rr.Code)
	}

	if len(app.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d",
			len(app.Messages))
	}

	msg := app.Messages[1]
	if msg == nil {
		t.Fatal("message not found")
	}

	if msg.Data != messageContent {
		t.Errorf("expected message %q, got %q",
			messageContent, msg.Data)
	}
}

// TestMessageHandlerExceedLength tests message
// exceeding configured length causes an error.
func TestMessageHandlerExceedLength(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", strings.Repeat("a", 1000))
	req := httptest.NewRequest(
		"POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)

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

	app.Messages[1] = &storage.Message{
		Count: 1, Data: messageContent,
	}
	app.NumMessages = 1

	form := url.Values{}
	form.Set("clear", "true")
	req := httptest.NewRequest(
		"POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)

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

	app.Messages[1] = &storage.Message{
		Count: 1, Data: messageContent + "1",
	}
	app.Messages[2] = &storage.Message{
		Count: 2, Data: messageContent + "2",
	}
	app.Messages[3] = &storage.Message{
		Count: 3, Data: messageContent + "3",
	}
	app.NumMessages = 3

	req := httptest.NewRequest("GET", "/?download=all", nil)
	rr := httptest.NewRecorder()
	handler := Message(app)
	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	for i := 1; i <= 3; i++ {
		want := messageContent + fmt.Sprint(i)
		if !strings.Contains(body, want) {
			t.Errorf("expected message %q, got %q",
				want, body)
		}
	}
}

// TestMessageHandlerExceedCount tests message
// exceeding configured count causes an error.
func TestMessageHandlerExceedCount(t *testing.T) {
	app := newTestApp()

	handler := Message(app)
	form := url.Values{}

	for i := 0; i < app.MessageLimits.MaxCount; i++ {
		form.Set("message", fmt.Sprintf("msg %d", i+1))
		req := httptest.NewRequest(
			"POST", "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code < 200 || rr.Code >= 400 {
			t.Errorf("expected 2xx, got %d", rr.Code)
		}
	}

	form.Set("message", messageContent)
	req := httptest.NewRequest(
		"POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", contentType)
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

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

// TestMessageAdd tests Message add.
func TestMessageAdd(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", testContentMsgs)
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
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

// TestMessageExceedLength tests Message add exceeding
// the configured length.
func TestMessageExceedLength(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", strings.Repeat("a", 1000))
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
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

// TestMessageExceedCount tests Message add exceeding
// the configured count.
func TestMessageExceedCount(t *testing.T) {
	app := newTestApp()

	handler := Message(app)
	form := url.Values{}

	for i := range app.MessageLimits.MaxCount {
		form.Set("message", fmt.Sprintf("msg %d", i+1))
		req := httptest.NewRequestWithContext(t.Context(),
			http.MethodPost, app.Message,
			strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", formContentType)

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		if rr.Code != 200 {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	}

	form.Set("message", testContentMsgs)
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
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

// TestMessageDeny tests Message add with no auth.
func TestMessageDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Message = true

	app.Messages = append(
		app.Messages, &storage.Message{
			Count: 1, Data: "existing"})

	form := url.Values{}
	form.Set("message", testContentMsgs)

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := serveDeniedRequest(t, Message(app), req)

	assertDenied(t, rr, app.Deny)

	if got := len(app.Messages); got != 1 {
		t.Fatalf("expected messages unchanged, got %d",
			got)
	}
	if app.Messages[0].Data != "existing" {
		t.Fatalf("expected existing message, got %q",
			app.Messages[0].Data)
	}
}

// TestMessageClear tests Messages clear.
func TestMessageClear(t *testing.T) {
	app := newTestApp()

	app.Messages = append(
		app.Messages, &storage.Message{
			Count: 1, Data: testContentMsgs})

	form := url.Values{}
	form.Set("clear", "true")
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()

	handler := Message(app)
	handler.ServeHTTP(rr, req)

	assertMessagesClear(t, app)
}

// TestMessageDownloadAll test Messages download all.
func TestMessageDownloadAll(t *testing.T) {
	app := newTestApp()

	app.Messages = append(
		app.Messages, &storage.Message{
			Count: 1, Data: testContentMsgs + "1"})
	app.Messages = append(
		app.Messages, &storage.Message{
			Count: 2, Data: testContentMsgs + "2"})
	app.Messages = append(
		app.Messages, &storage.Message{
			Count: 3, Data: testContentMsgs + "3"})

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

// TestMessageBrowser tests Message add with browser.
func TestMessageBrowser(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", testContentMsgs)
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)
	req.Header.Set("Accept", "text/html")

	rr := httptest.NewRecorder()
	Message(app).ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected %d, got %d",
			http.StatusSeeOther, rr.Code)
	}

	if got := rr.Header().Get("Location"); got != app.Root {
		t.Fatalf("expected redirect to %q, got %q",
			app.Root, got)
	}

	if got := len(app.Messages); got != 1 {
		t.Fatalf("expected 1 message, got %d", got)
	}
}

// TestMessageSpaces tests Message with spaces add.
func TestMessageSpaces(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", "  \n\t hello, world! \r\n ")
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()

	handler := Message(app)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if got := len(app.Messages); got != 1 {
		t.Fatalf("expected 1 message, got %d", got)
	}

	msg := app.Messages[0]
	if msg == nil {
		t.Fatal("message not found")
	}

	if got := msg.Data; got != testContentMsgs {
		t.Fatalf("expected trim message %q, got %q",
			testContentMsgs, got)
	}
}

// TestMessageSpacesOnly tests spaces-only Message add.
func TestMessageSpacesOnly(t *testing.T) {
	app := newTestApp()

	form := url.Values{}
	form.Set("message", "   \n\t  ")
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()
	Message(app).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if got := len(app.Messages); got != 0 {
		t.Fatalf("expected 0 messages, got %d", got)
	}
}

// TestMessageClearAdd tests Message clear and add.
func TestMessageClearAdd(t *testing.T) {
	app := newTestApp()
	app.Messages = append(
		app.Messages, &storage.Message{
			Count: 1, Data: "previous message"})

	form := url.Values{}
	form.Set("clear", "true")
	form.Set("message", testContentMsgs)

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Message,
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formContentType)

	rr := httptest.NewRecorder()
	Message(app).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if got := len(app.Messages); got != 1 {
		t.Fatalf("expected 1 message, got %d", got)
	}

	if got := app.Messages[0].Data; got != testContentMsgs {
		t.Fatalf("expected %q, got %q", testContentMsgs, got)
	}

	if got := app.Messages[0].Count; got != 1 {
		t.Fatalf("expected count of 1, got %d", got)
	}
}

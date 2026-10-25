package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// newMultipartUpload builds a multipart body with the
// token as a form field (as the HTML form sends it)
// and one file of size bytes.
func newMultipartUpload(t *testing.T, field, token string,
	size int) (*bytes.Buffer, string) {
	t.Helper()

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	if err := mw.WriteField(field, token); err != nil {
		t.Fatalf("write token field: %v", err)
	}
	fw, err := mw.CreateFormFile("file", "big.bin")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write(bytes.Repeat([]byte("x"), size)); err != nil {
		t.Fatalf("write file content: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	return &b, mw.FormDataContentType()
}

// TestUploadFormTokenChunkedTooLarge tests the per-file limit
// holds when the token is a form field and no Content-Length
// is sent, so the body is bounded before auth parses it.
func TestUploadFormTokenChunkedTooLarge(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = true
	app.FileLimits.SizeEachMb = 1

	body, ctype := newMultipartUpload(t,
		app.Basic.Field, app.Basic.Token, 1_500_000)

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, body)
	req.Header.Set("Content-Type", ctype)
	req.ContentLength = -1
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusRequestEntityTooLarge)

	if len(app.Files) != 0 {
		t.Fatalf("expected no files stored, got %d",
			len(app.Files))
	}
}

// TestUploadFormTokenOK tests a form-field token upload
// within the limit succeeds.
func TestUploadFormTokenOK(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = true
	app.FileLimits.SizeEachMb = 1

	body, ctype := newMultipartUpload(t,
		app.Basic.Field, app.Basic.Token, 500_000)

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, body)
	req.Header.Set("Content-Type", ctype)
	req.ContentLength = -1
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	if len(app.Files) != 1 {
		t.Fatalf("expected 1 file stored, got %d",
			len(app.Files))
	}
}

// TestWallBodyTooLarge tests a form post beyond the
// derived form limit is rejected before parsing.
func TestWallBodyTooLarge(t *testing.T) {
	app := newTestApp()
	app.Require.Wall = false

	limit := maxFormBytes(app)
	body := "wall=" + strings.Repeat("a", int(limit))

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.WallModify, strings.NewReader(body))
	req.Header.Set("Content-Type", formContentType)
	req.ContentLength = -1
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusRequestEntityTooLarge)

	if app.WallContent != "" {
		t.Fatal("expected wall unchanged")
	}
}

// TestTransferDeadline tests deadline scaling with size.
func TestTransferDeadline(t *testing.T) {
	base := 20 * time.Second
	n := int64(10 * minTransferRate)

	before := time.Now()
	got := transferDeadline(base, n)
	want := before.Add(base + 10*time.Second)

	if got.Before(want) {
		t.Fatalf("deadline %v before %v", got, want)
	}
	if got.After(want.Add(time.Second)) {
		t.Fatalf("deadline %v too far after %v", got, want)
	}
}

// TestBodyLimit tests upload vs form body limits.
func TestBodyLimit(t *testing.T) {
	app := newTestApp()

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, nil)
	req.Pattern = app.Upload
	if got := bodyLimit(app, req); got != app.GetMaxFileBytes() {
		t.Fatalf("upload limit = %d, want %d",
			got, app.GetMaxFileBytes())
	}

	req.Pattern = app.WallModify
	if got := bodyLimit(app, req); got != maxFormBytes(app) {
		t.Fatalf("form limit = %d, want %d",
			got, maxFormBytes(app))
	}
}

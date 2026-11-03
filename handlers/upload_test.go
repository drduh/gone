package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/drduh/gone/storage"
)

// TestUploadFileTooLarge test file uploads exceeding size limit.
func TestUploadFileTooLarge(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = false

	reqBody := strings.NewReader("upload_test")
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, reqBody)
	req.ContentLength = (app.FileLimits.SizeEachMb << 20) + 1
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusRequestEntityTooLarge)

	var resp map[string]string
	if err := json.Unmarshal(
		rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp["error"] != app.FileSize {
		t.Fatalf("expected error %q, got %q",
			app.FileSize, resp["error"])
	}
}

// TestUploadNoFile tests uploads without a file selected.
func TestUploadNoFile(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = false

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, &b)
	req.Header.Set("Content-Type",
		mw.FormDataContentType())
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusBadRequest)

	var resp map[string]string
	if err := json.Unmarshal(
		rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp["error"] != app.Form {
		t.Fatalf("expected error %q, got %q",
			app.Form, resp["error"])
	}
}

// TestUploadSuccess tests successful file uploads.
func TestUploadSuccess(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = false

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("file", "upload_test.txt")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write([]byte("hello world")); err != nil {
		t.Fatalf("write file content: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, &b)
	req.Header.Set("Content-Type",
		mw.FormDataContentType())
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusOK)

	var uploads []storage.File
	if err := json.Unmarshal(
		rr.Body.Bytes(), &uploads); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(uploads) != 1 {
		t.Fatalf("expected 1 upload, got %d",
			len(uploads))
	}

	if uploads[0].Name != "upload_test.txt" {
		t.Fatalf("expected filename %q, got %q",
			"upload_test.txt", uploads[0].Name)
	}

	found := false

	for _, f := range app.Files {
		if f.Name == "upload_test.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("file not found in storage")
	}
}

// TestUploadDeny tests denied Upload requests.
func TestUploadDeny(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = true

	reqBody := strings.NewReader("upload_test")
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, reqBody)
	rr := serveDeniedRequest(t, app, req)
	assertDenied(t, rr, app.Deny)

	if len(app.Files) != 0 {
		t.Fatalf("expected no files uploaded, got %d",
			len(app.Files))
	}
}

// TestUploadTooManyFiles tests the selected files limit.
func TestUploadTooManyFiles(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = false
	app.FileLimits.MaxSelectFiles = 1

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	for _, name := range []string{"a.txt", "b.txt"} {
		fw, err := mw.CreateFormFile("file", name)
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		if _, err := fw.Write([]byte("x")); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusBadRequest)
	if !strings.Contains(rr.Body.String(), app.UploadCount) {
		t.Fatalf("body = %q", rr.Body.String())
	}
}

// TestUploadTotalLimit tests the total stored size limit.
func TestUploadTotalLimit(t *testing.T) {
	app := newTestAppWithStorage()
	app.Require.Upload = false
	app.FileLimits.SizeTotalMb = 1

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("file", "big.bin")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write(bytes.Repeat([]byte("x"),
		1_000_000)); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusRequestEntityTooLarge)
	if !strings.Contains(rr.Body.String(), app.FileSizeAll) {
		t.Fatalf("body = %q", rr.Body.String())
	}
}

// TestUploadNotMultipart tests a non-multipart body.
func TestUploadNotMultipart(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = false

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, strings.NewReader("file=x"))
	req.Header.Set("Content-Type", formContentType)
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusBadRequest)
}

// TestUploadBrowser tests browser uploads redirect.
func TestUploadBrowser(t *testing.T) {
	app := newTestApp()
	app.Require.Upload = false

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("file", "a.txt")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write([]byte("x")); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodPost, app.Upload, &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Accept", "text/html")
	rr := serveRequest(t, app, req)
	assertStatus(t, rr, http.StatusSeeOther)
	if len(app.Files) != 1 {
		t.Fatalf("files = %d; want 1", len(app.Files))
	}
}

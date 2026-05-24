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

	reqBody := strings.NewReader("dummy")
	req := httptest.NewRequest(http.MethodPost, "/upload", reqBody)
	req.ContentLength = (app.FileLimits.SizeEachMb << 20) + 1

	w := httptest.NewRecorder()
	handler := Upload(app)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d",
			http.StatusRequestEntityTooLarge, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
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

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	handler := Upload(app)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
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

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	handler := Upload(app)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, w.Code)
	}

	var uploads []storage.File
	if err := json.Unmarshal(w.Body.Bytes(), &uploads); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(uploads) != 1 {
		t.Fatalf("expected 1 upload, got %d", len(uploads))
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

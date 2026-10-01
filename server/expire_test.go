package server

import (
	"testing"
	"time"

	"github.com/drduh/gone/storage"
)

// TestExpireFilesDownloads tests File downloads limit.
func TestExpireFilesDownloads(t *testing.T) {
	f := &storage.File{
		ID:         "file1",
		Name:       "file1.txt",
		Allow:      1,
		Count:      1,
		UploadTime: time.Now(),
	}

	app := newTestApp(map[string]*storage.File{f.ID: f})
	expireFiles(app)

	if _, ok := app.Files[f.ID]; ok {
		t.Fatalf("expected %q to be removed", f.ID)
	}
}

// TestExpireFilesDuration tests File duration limit.
func TestExpireFilesDuration(t *testing.T) {
	f := &storage.File{
		ID:         "file2",
		Name:       "file2.txt",
		Duration:   10 * time.Second,
		UploadTime: time.Now().Add(-30 * time.Second),
	}

	app := newTestApp(map[string]*storage.File{f.ID: f})
	expireFiles(app)

	if _, ok := app.Files[f.ID]; ok {
		t.Fatalf("expected %q to be removed", f.ID)
	}
}

// TestExpireFilesValid tests File limits not met.
func TestExpireFilesValid(t *testing.T) {
	f := &storage.File{
		ID:         "file3",
		Name:       "file3.txt",
		Allow:      3,
		Count:      1,
		Duration:   5 * time.Minute,
		UploadTime: time.Now().Add(-30 * time.Second),
	}

	app := newTestApp(map[string]*storage.File{f.ID: f})
	expireFiles(app)

	if _, ok := app.Files[f.ID]; !ok {
		t.Fatalf("expected %q to remain", f.ID)
	}
}

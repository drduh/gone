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

// TestRunExpiryLoop tests Files expire on each tick
// and the loop ends when the channel closes.
func TestRunExpiryLoop(t *testing.T) {
	f := &storage.File{
		ID:         "tick",
		Name:       "tick.txt",
		Allow:      1,
		Count:      1,
		UploadTime: time.Now(),
	}

	app := newTestApp(map[string]*storage.File{f.ID: f})

	ticks := make(chan time.Time, 1)
	done := make(chan struct{})
	go func() {
		runExpiryLoop(app, ticks)
		close(done)
	}()

	ticks <- time.Now()
	close(ticks)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("expiry loop did not stop")
	}

	if _, ok := app.Files[f.ID]; ok {
		t.Fatalf("expected %q removed on tick", f.ID)
	}
}

// TestExpiryWorker tests the worker removes an expired File.
func TestExpiryWorker(t *testing.T) {
	f := &storage.File{
		ID:         "worker",
		Name:       "worker.txt",
		Allow:      1,
		Count:      1,
		UploadTime: time.Now(),
	}

	app := newTestApp(map[string]*storage.File{f.ID: f})
	app.FileLimits.ExpiryCheck.Duration = 10 * time.Millisecond

	go expiryWorker(app)

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if app.FindFile(f.ID) == nil {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("expected %q removed by worker", f.ID)
}

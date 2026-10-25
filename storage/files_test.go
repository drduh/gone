package storage

import (
	"sync"
	"testing"
	"time"
)

func newFile(id string, size int, allow int) *File {
	return &File{
		ID:         id,
		Name:       id + ".txt",
		Data:       make([]byte, size),
		Bytes:      size,
		Allow:      allow,
		Duration:   time.Hour,
		UploadTime: time.Now(),
	}
}

// TestAddFilesLimit tests the total size limit is
// all-or-nothing.
func TestAddFilesLimit(t *testing.T) {
	var s Storage

	if !s.AddFiles([]*File{newFile("a", 60, 1)}, 100) {
		t.Fatal("expected first add within limit")
	}

	if s.AddFiles([]*File{
		newFile("b", 20, 1), newFile("c", 30, 1)}, 100) {
		t.Fatal("expected add over limit rejected")
	}

	if len(s.Files) != 1 || s.SizeFiles != 60 {
		t.Fatalf("files = %d size = %d; want 1, 60",
			len(s.Files), s.SizeFiles)
	}
	if !s.AddFiles([]*File{newFile("b", 20, 1)}, 0) {
		t.Fatal("expected add with no limit")
	}
}

// TestReserveDownload tests counting, removal on the
// last download, and not-found cases.
func TestReserveDownload(t *testing.T) {
	var s Storage
	s.AddFiles([]*File{newFile("a", 10, 2)}, 0)

	f, reason := s.ReserveDownload("a.txt")
	if f == nil || reason != "" || f.Count != 1 || f.Remain != 1 {
		t.Fatalf("first: f=%v reason=%q", f, reason)
	}

	if len(f.Data) != 10 {
		t.Fatal("expected Data on served copy")
	}

	f, reason = s.ReserveDownload("a")
	if f == nil || reason != expiryDownloadLimit || f.Count != 2 {
		t.Fatalf("last: f=%v reason=%q", f, reason)
	}

	if _, ok := s.Files["a"]; ok {
		t.Fatal("expected file removed after last download")
	}

	if f, _ := s.ReserveDownload("a"); f != nil {
		t.Fatal("expected nil after removal")
	}
}

// TestReserveDownloadExpired tests an expired File is
// removed rather than served.
func TestReserveDownloadExpired(t *testing.T) {
	var s Storage
	f := newFile("a", 10, 5)
	f.UploadTime = time.Now().Add(-2 * time.Hour)
	s.AddFiles([]*File{f}, 0)

	got, reason := s.ReserveDownload("a")
	if got != nil || reason != expiryDurationLimit {
		t.Fatalf("got=%v reason=%q", got, reason)
	}

	if len(s.Files) != 0 {
		t.Fatal("expected expired file removed")
	}
}

// TestReserveDownloadConcurrent tests only Allow
// reservations succeed under concurrency.
func TestReserveDownloadConcurrent(t *testing.T) {
	var s Storage
	s.AddFiles([]*File{newFile("a", 10, 3)}, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex
	served := 0
	for range 50 {
		wg.Go(func() {
			if f, _ := s.ReserveDownload("a"); f != nil {
				mu.Lock()
				served++
				mu.Unlock()
			}
		})
	}
	wg.Wait()

	if served != 3 {
		t.Fatalf("served %d; want 3", served)
	}
}

// TestExpireFiles tests expired Files are removed.
func TestExpireFiles(t *testing.T) {
	var s Storage
	old := newFile("old", 1, 5)
	old.UploadTime = time.Now().Add(-2 * time.Hour)
	used := newFile("used", 1, 1)
	used.Count = 1
	s.AddFiles([]*File{old, used, newFile("live", 1, 5)}, 0)

	expired := s.ExpireFiles()
	if len(expired) != 2 {
		t.Fatalf("expired %d; want 2", len(expired))
	}

	for _, e := range expired {
		if e.Data != nil {
			t.Fatal("expected no Data in expired copy")
		}
		if e.Reason == "" {
			t.Fatalf("%s: expected reason", e.ID)
		}
	}

	if _, ok := s.Files["live"]; !ok || len(s.Files) != 1 {
		t.Fatal("expected only live file to remain")
	}
}

// TestSnapshot tests a detached, ordered copy.
func TestSnapshot(t *testing.T) {
	var s Storage
	b := newFile("b", 5, 1)
	a := newFile("a", 5, 1)
	a.UploadTime = b.UploadTime.Add(-time.Minute)
	s.AddFiles([]*File{b, a}, 0)
	s.AddMessage(&Message{Data: "hi"}, 0)
	s.SetWall("x\ny", WallMeta{WallModifiedUser: "u"}, 10)

	snap := s.Snapshot()

	if len(snap.Files) != 2 || snap.Files[0].ID != "a" {
		t.Fatalf("files = %v; want a first", snap.Files)
	}

	if snap.Files[0].Data != nil {
		t.Fatal("expected no Data in snapshot")
	}

	if snap.NumFiles != 2 || snap.SizeFiles != 10 {
		t.Fatalf("sizes = %+v", snap.Sizes)
	}

	if len(snap.Messages) != 1 || snap.Messages[0].Count != 1 {
		t.Fatalf("messages = %v", snap.Messages)
	}

	if snap.WallContent != "x\ny" || snap.WallLines != 2 ||
		snap.WallCap != 30 || snap.WallModifiedUser != "u" {
		t.Fatalf("wall = %q meta = %+v",
			snap.WallContent, snap.WallMeta)
	}

	snap.Messages[0].Data = "changed"
	if s.Messages[0].Data != "hi" {
		t.Fatal("expected snapshot detached from storage")
	}
}

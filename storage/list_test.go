package storage

import (
	"slices"
	"testing"
	"time"
)

// TestListFiles tests active and expired Files.
func TestListFiles(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"durationExpire": {
				ID:   "durationExpire",
				Name: "durationExpire.txt",
				Downloads: Downloads{
					Allow: 5,
					Count: 0,
				},
				Time: Time{
					Duration: time.Second,
					UploadTime: time.Date(
						2006, 1, 1, 0, 0, 0, 0, time.UTC,
					),
				},
			},
			"downloadExpire": {
				ID:   "downloadExpire",
				Name: "downloadExpire.txt",
				Downloads: Downloads{
					Allow: 1,
					Count: 1,
				},
			},
			"active1": {
				ID:   "active1",
				Name: "active1.txt",
				Downloads: Downloads{
					Allow: 2,
					Count: 0,
				},
				Time: Time{
					Duration:   time.Minute,
					UploadTime: time.Now(),
				},
			},
			"active2": {
				ID:   "active2",
				Name: "active2.txt",
				Downloads: Downloads{
					Allow: 3,
					Count: 1,
				},
				Time: Time{
					Duration:   time.Hour,
					UploadTime: time.Now(),
				},
			},
		},
	}

	got := s.ListFiles()

	if len(got) != 2 {
		t.Fatalf("listed %d files; want 2", len(got))
	}

	gotByID := make(map[string]File, len(got))
	for _, f := range got {
		gotByID[f.ID] = f
	}

	gotIDs := make([]string, 0, len(gotByID))
	for id := range gotByID {
		gotIDs = append(gotIDs, id)
	}
	slices.Sort(gotIDs)

	wantIDs := []string{"active1", "active2"}
	if !slices.Equal(gotIDs, wantIDs) {
		t.Fatalf("got files %v; want %v", gotIDs, wantIDs)
	}

	for _, id := range []string{
		"durationExpire", "downloadExpire"} {
		if _, ok := s.Files[id]; ok {
			t.Fatalf("expired file %q present", id)
		}
	}
	if got := len(s.Files); got != 2 {
		t.Fatalf("%d files stored; want 2", got)
	}

	if got := gotByID["active1"].Remain; got != 2 {
		t.Errorf("active1 downloads = %d; want 2", got)
	}
	if got := gotByID["active2"].Remain; got != 2 {
		t.Errorf("active2 downloads = %d; want 2", got)
	}
}

// TestListFilesEmpty tests no Files to list.
func TestListFilesEmpty(t *testing.T) {
	s := &Storage{Files: nil}

	got := s.ListFiles()

	if len(got) != 0 {
		t.Fatalf("got %d files; want 0", len(got))
	}
}

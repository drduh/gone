package storage

import (
	"slices"
	"testing"
	"time"
)

// TestListFiles tests active and expired Files.
func TestListFiles(t *testing.T) {
	now := time.Now()

	s := &Storage{
		Files: map[string]*File{
			"durationExpired": {
				Id:   "durationExpire",
				Name: "durationExpire.txt",
				Downloads: Downloads{
					Allow: 5,
					Count: 0,
				},
				Time: Time{
					Duration: time.Second,
					Upload:   now.Add(-time.Minute),
				},
			},
			"downloadExpired": {
				Id:   "downloadExpire",
				Name: "downloadExpire.txt",
				Downloads: Downloads{
					Allow: 1,
					Count: 1,
				},
				Time: Time{
					Duration: time.Hour,
					Upload:   now,
				},
			},
			"active1": {
				Id:   "active1",
				Name: "active1.txt",
				Downloads: Downloads{
					Allow: 2,
					Count: 0,
				},
				Time: Time{
					Duration: time.Minute,
					Upload:   now,
				},
			},
			"active2": {
				Id:   "active2",
				Name: "active2.txt",
				Downloads: Downloads{
					Allow: 3,
					Count: 1,
				},
				Time: Time{
					Duration: time.Hour,
					Upload:   now,
				},
			},
		},
	}

	got := s.ListFiles()

	if len(got) != 2 {
		t.Fatalf("listed %d files; want 2", len(got))
	}

	gotIDs := make([]string, 0, len(got))
	for _, f := range got {
		gotIDs = append(gotIDs, f.Id)
	}
	slices.Sort(gotIDs)

	wantIDs := []string{"active1", "active2"}
	if !slices.Equal(gotIDs, wantIDs) {
		t.Fatalf("got files %v; want %v",
			gotIDs, wantIDs)
	}

	if _, ok := s.Files["expired"]; ok {
		t.Fatalf("expired files still present")
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

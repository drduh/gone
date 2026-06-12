package storage

import "testing"

// TestFindFileID tests a File found in Storage by id.
func TestFindFileID(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"1": {ID: "id1", Name: "file1"},
			"2": {ID: "id2", Name: "file2"},
		},
	}

	got := s.FindFile("id2")
	if got == nil {
		t.Fatalf("FindFile returned nil")
	}
	if got.ID != "id2" {
		t.Fatalf("FindFile returned %q; want %q",
			got.ID, "id2")
	}
	if got.Name != "file2" {
		t.Fatalf("FindFile returned %q; want %q",
			got.Name, "file2")
	}
}

// TestFindFileName tests a File found in Storage by name.
func TestFindFileName(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"1": {ID: "id1", Name: "file1"},
			"2": {ID: "id2", Name: "file2"},
		},
	}

	got := s.FindFile("file1")
	if got == nil {
		t.Fatalf("FindFile returned nil")
	}
	if got.ID != "id1" {
		t.Fatalf("FindFile returned %q; want %q",
			got.ID, "id1")
	}
	if got.Name != "file1" {
		t.Fatalf("FindFile returned %q; want %q",
			got.Name, "file1")
	}
}

// TestFindFileNotFound tests a File not found in Storage.
func TestFindFileNotFound(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"1": {ID: "id1", Name: "file1"},
		},
	}

	got := s.FindFile("file2")
	if got != nil {
		t.Fatalf("FindFile returned %#v; want nil", got)
	}
}

// TestFindFileEmptyStorage tests empty Storage (no Files).
func TestFindFileEmptyStorage(t *testing.T) {
	s := &Storage{Files: nil}

	got := s.FindFile("file1")
	if got != nil {
		t.Fatalf("FindFile returned %#v; want nil", got)
	}
}

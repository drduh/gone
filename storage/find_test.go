package storage

import "testing"

// TestFindFileId tests a File found in Storage by id.
func TestFindFileId(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"1": {Id: "id1", Name: "file1"},
			"2": {Id: "id2", Name: "file2"},
		},
	}
	got := s.FindFile("id2")
	if got == nil {
		t.Fatalf("FindFile returned nil")
	}
	if got.Id != "id2" {
		t.Fatalf("FindFile returned %q; want %q", got.Id, "id2")
	}
}

// TestFindFileName tests a File found in Storage by name.
func TestFindFileName(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"1": {Id: "id1", Name: "file1"},
			"2": {Id: "id2", Name: "file2"},
		},
	}
	got := s.FindFile("file1")
	if got == nil {
		t.Fatalf("FindFile returned nil")
	}
	if got.Name != "file1" {
		t.Fatalf("FindFile returned %q; want %q", got.Name, "file1")
	}
}

// TestFindFileNotFound tests a File not found in Storage by name.
func TestFindFileNotFound(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"1": {Id: "id1", Name: "file1"},
		},
	}
	got := s.FindFile("file2")
	if got != nil {
		t.Fatalf("FindFile returned %#v; want nil", got)
	}
}

// TestFindFileNil tests no Files are found in empty Storage.
func TestFindFileNil(t *testing.T) {
	s := &Storage{
		Files: nil,
	}
	got := s.FindFile("file1")
	if got != nil {
		t.Fatalf("FindFile returned %#v; want nil", got)
	}
}

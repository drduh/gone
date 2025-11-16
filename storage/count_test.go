package storage

import "testing"

// TestCountFiles tests File counts in Storage.
func TestCountFiles(t *testing.T) {
	var s Storage
	if got := s.CountFiles(); got != 0 {
		t.Fatalf("CountFiles = %d; want 0", got)
	}
	s.Files = map[string]*File{
		"file1": {},
		"file2": {},
	}
	if got := s.CountFiles(); got != 2 {
		t.Fatalf("CountFiles = %d; want 2", got)
	}
}

// TestCountMessages tests Messages counts in Storage.
func TestCountMessages(t *testing.T) {
	var s Storage
	if got := s.CountMessages(); got != 0 {
		t.Fatalf("CountMessages = %d; want 0", got)
	}
	s.Messages = map[int]*Message{
		1: {},
	}
	if got := s.CountMessages(); got != 1 {
		t.Fatalf("CountMessages = %d; want 1", got)
	}
	s.Messages[2] = &Message{}
	if got := s.CountMessages(); got != 2 {
		t.Fatalf("CountMessages = %d; want 2", got)
	}
}

package storage

import "testing"

// TestCountFiles tests File counts in Storage.
func TestCountFiles(t *testing.T) {
	var s Storage
	s.CountFiles()
	if s.NumFiles != 0 {
		t.Fatalf("NumFiles = %d; want 0", s.NumFiles)
	}
	s.Files = map[string]*File{
		"file1": {},
		"file2": {},
	}
	s.CountFiles()
	if s.NumFiles != 2 {
		t.Fatalf("NumFiles = %d; want 2", s.NumFiles)
	}
}

// TestCountMessages tests Messages counts in Storage.
func TestCountMessages(t *testing.T) {
	var s Storage
	s.CountMessages()
	if s.NumMessages != 0 {
		t.Fatalf("NumMessages = %d; want 0", s.NumMessages)
	}
	s.Messages = map[int]*Message{
		1: {},
	}
	s.CountMessages()
	if s.NumMessages != 1 {
		t.Fatalf("NumMessages = %d; want 1", s.NumMessages)
	}
	s.Messages[2] = &Message{}
	s.CountMessages()
	if s.NumMessages != 2 {
		t.Fatalf("NumMessages = %d; want 2", s.NumMessages)
	}
}

// TestCountWall tests Wall content counts in Storage.
func TestCountWall(t *testing.T) {
	var s Storage
	s.CountWall()
	if s.CharsWall != 0 {
		t.Fatalf("CountWall = %d; want 0", s.CharsWall)
	}
	s.WallContent = "test\r\nwall"
	s.CountWall()
	if s.CharsWall != 10 {
		t.Fatalf("CharsWall = %d; want 10", s.CharsWall)
	}
	s.WallContent = ""
	s.CountWall()
	if s.CharsWall != 0 {
		t.Fatalf("CharsWall = %d; want 0", s.CharsWall)
	}
}

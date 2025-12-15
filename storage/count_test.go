package storage

import "testing"

// TestCountFiles tests File counts in Storage.
func TestCountFiles(t *testing.T) {
	var s Storage
	s.CountFiles()
	if s.NumFiles != 0 {
		t.Fatalf("NumFiles = %d; want 0", s.NumFiles)
	}
	if s.SizeFiles != 0 {
		t.Fatalf("SizeFiles = %d; want 0", s.SizeFiles)
	}
	if s.SizeFilesFmt != "" {
		t.Fatalf("SizeFilesFmt = %s; want empty string",
			s.SizeFilesFmt)
	}

	s.Files = map[string]*File{
		"file1": {
			Data:  []byte("test content 1"),
			Bytes: 14,
		},
		"file2": {
			Data:  []byte("test content 2"),
			Bytes: 14,
		},
	}
	s.CountFiles()
	if s.NumFiles != 2 {
		t.Fatalf("NumFiles = %d; want 2", s.NumFiles)
	}
	if s.SizeFiles != 28 {
		t.Fatalf("SizeFiles = %d; want 0", s.SizeFiles)
	}
	if s.SizeFilesFmt != "28 bytes" {
		t.Fatalf("SizeFilesFmt = %s; want '28 bytes'",
			s.SizeFilesFmt)
	}
}

// TestCountMessages tests Messages counts in Storage.
func TestCountMessages(t *testing.T) {
	var s Storage
	s.CountMessages()
	if s.CharsMessages != 0 {
		t.Fatalf("CharsMessages = %d; want 0", s.CharsMessages)
	}
	if s.NumMessages != 0 {
		t.Fatalf("NumMessages = %d; want 0", s.NumMessages)
	}

	s.Messages = map[int]*Message{
		1: {Data: "test message 1"},
	}
	s.CountMessages()
	if s.NumMessages != 1 {
		t.Fatalf("NumMessages = %d; want 1", s.NumMessages)
	}
	if s.CharsMessages != 14 {
		t.Fatalf("CharsMessages = %d; want 14", s.CharsMessages)
	}

	s.Messages[2] = &Message{
		Data: "test message 2",
	}
	s.CountMessages()
	if s.CharsMessages != 28 {
		t.Fatalf("CharsMessages = %d; want 28", s.CharsMessages)
	}
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
	if s.LinesWall != 0 {
		t.Fatalf("LinesWall = %d; want 0", s.LinesWall)
	}

	s.WallContent = "test\r\nwall"
	s.CountWall()
	if s.CharsWall != 10 {
		t.Fatalf("CharsWall = %d; want 10", s.CharsWall)
	}
	if s.LinesWall != 2 {
		t.Fatalf("LinesWall = %d; want 2", s.LinesWall)
	}

	s.WallContent = ""
	s.CountWall()
	if s.CharsWall != 0 {
		t.Fatalf("CharsWall = %d; want 0", s.CharsWall)
	}
	if s.LinesWall != 0 {
		t.Fatalf("LinesWall = %d; want 0", s.LinesWall)
	}
}

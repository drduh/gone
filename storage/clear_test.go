package storage

import "testing"

// TestClearFiles tests removal of Files in Storage.
func TestClearFiles(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"file1": {},
			"file2": {},
		},
	}
	s.ClearFiles()
	if s.Files == nil {
		t.Fatalf("Files is nil; want empty map")
	}
	if got := len(s.Files); got != 0 {
		t.Fatalf("Files length = %d; want 0", got)
	}
	s = &Storage{}
	s.ClearFiles()
	if s.Files == nil {
		t.Fatalf("Files is nil; want empty map")
	}
	if got := len(s.Files); got != 0 {
		t.Fatalf("Files length = %d; want 0", got)
	}
}

// TestClearMessages tests removal of Messages in Storage.
func TestClearMessages(t *testing.T) {
	s := &Storage{
		Messages: map[int]*Message{
			1: {},
			2: {},
		},
	}
	s.ClearMessages()
	if s.Messages == nil {
		t.Fatalf("Messages is nil; want empty map")
	}
	if got := len(s.Messages); got != 0 {
		t.Fatalf("Messages length = %d; want 0", got)
	}
	s = &Storage{}
	s.ClearMessages()
	if s.Messages == nil {
		t.Fatalf("Messages is nil; want empty map")
	}
	if got := len(s.Messages); got != 0 {
		t.Fatalf("Messages length = %d; want 0", got)
	}
}

// TestClearWall tests removal of Wall content in Storage.
func TestClearWall(t *testing.T) {
	s := &Storage{
		WallContent: "test wall content",
	}
	s.ClearWall()
	if s.WallContent != "" {
		t.Fatalf("WallContent = %q; want empty string", s.WallContent)
	}
}

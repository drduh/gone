package storage

import "testing"

// TestClearFiles tests clearing Files in Storage.
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
}

// TestClearMessages tests clearing Messages in Storage.
func TestClearMessages(t *testing.T) {
	s := &Storage{
		Messages: []*Message{
			{Count: 1, Data: "hello"},
			{Count: 2, Data: "world"},
		},
	}
	s.ClearMessages()
	if s.Messages == nil {
		t.Fatalf("Messages is nil; want empty map")
	}
	if got := len(s.Messages); got != 0 {
		t.Fatalf("Messages length = %d; want 0", got)
	}
}

// TestClearWall tests clearing Wall content in Storage.
func TestClearWall(t *testing.T) {
	s := &Storage{WallContent: "test wall content"}
	s.ClearWall()
	if s.WallContent != "" {
		t.Fatalf("WallContent = %q; want empty string",
			s.WallContent)
	}
}

// TestClearStorage tests clearing all Storage content.
func TestClearStorage(t *testing.T) {
	s := &Storage{
		Files: map[string]*File{
			"file1": {},
			"file2": {},
		},
		Messages: []*Message{
			{Count: 1, Data: "hello"},
			{Count: 2, Data: "world"},
		},
		WallContent: "test wall content",
	}
	s.ClearStorage()
	if s.Files == nil {
		t.Fatalf("Files is nil; want empty map")
	}
	if got := len(s.Files); got != 0 {
		t.Fatalf("Files length = %d; want 0", got)
	}
	if s.Messages == nil {
		t.Fatalf("Messages is nil; want empty map")
	}
	if got := len(s.Messages); got != 0 {
		t.Fatalf("Messages length = %d; want 0", got)
	}
	if s.WallContent != "" {
		t.Fatalf("WallContent = %q; want empty string",
			s.WallContent)
	}
}

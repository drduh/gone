package storage

import (
	"regexp"
	"strings"
	"testing"
)

// TestSetId tests File id is a valid hex string.
func TestSetId(t *testing.T) {
	f := &File{}
	f.setId(32)
	if f.Id == "" {
		t.Fatalf("Id is empty; want hex string")
	}
	valid := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	if !valid.MatchString(f.Id) {
		t.Fatalf("Id = %q; want hex characters", f.Id)
	}
}

// TestSetSize tests File length and size are valid strings.
func TestSetSize(t *testing.T) {
	data := []byte("hello")
	length := "5"
	size := "5.00 Bytes"
	f := &File{Data: data}
	f.setSize()
	if got := f.Length; got != length {
		t.Fatalf("Length = %q; want %q", got, length)
	}
	if got := f.Size; got != size {
		t.Fatalf("Size = %q; want %q", got, size)
	}
}

// TestSetSum tests File hash is a valid sum string.
func TestSetSum(t *testing.T) {
	f := &File{Data: []byte("hello")}
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	f.setSum()
	if f.Sum != expected {
		t.Fatalf("Sum = %q; want %q", f.Sum, expected)
	}
}

// TestSetType test File type is a valid string based on extension.
func TestSetType(t *testing.T) {
	f1 := &File{Name: "notes.txt"}
	f1.setType()
	if !strings.HasPrefix(f1.Type, "text/plain") {
		t.Fatalf("Type = %q; want text/plain", f1.Type)
	}
	f2 := &File{Name: "notes.whatever"}
	f2.setType()
	if f2.Type != "application/octet-stream" {
		t.Fatalf("Type = %q; want application/octet-stream", f2.Type)
	}
}

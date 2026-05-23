package storage

import "testing"

const fileIdLength = 44 // storageVersion "1" + 32 bytes

// TestFileScan tests all scan attributes are set correctly.
func TestFileScan(t *testing.T) {
	files := []struct {
		Name   string
		Data   []byte
		Length string
		Size   int
		Sum    string
		Type   string
	}{
		{
			Name:   "test.txt",
			Data:   []byte("hello, world!\n"),
			Length: "14",
			Size:   14,
			Sum:    "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc",
			Type:   "text/plain; charset=utf-8",
		},
		{
			Name:   "empty.bin",
			Data:   []byte{},
			Length: "0",
			Size:   0,
			Sum:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			Type:   "application/octet-stream",
		},
	}

	for _, tc := range files {
		f := File{
			Name: tc.Name,
			Data: tc.Data,
		}
		f.Scan()

		if len(f.Id) != fileIdLength {
			t.Errorf("id length for %s: got %d, want %d",
				tc.Name, len(f.Id), fileIdLength)
		}
		if f.Bytes != tc.Size {
			t.Errorf("bytes for %s: got %d, want %d",
				tc.Name, f.Bytes, tc.Size)
		}
		if f.Length != tc.Length {
			t.Errorf("length for %s: got %s, want %s",
				tc.Name, f.Length, tc.Length)
		}
		if f.Sum != tc.Sum {
			t.Errorf("sum for %s: got %s, want %s",
				tc.Name, f.Sum, tc.Sum)
		}
		if f.Type != tc.Type {
			t.Errorf("type for %s: got %s, want %s",
				tc.Name, f.Type, tc.Type)
		}
	}
}

// TestSetId tests File id has valid length.
func TestSetId(t *testing.T) {
	f := &File{}
	f.setId()
	if f.Id == "" {
		t.Fatalf("id is empty")
	}
	if len(f.Id) != fileIdLength {
		t.Fatalf("id is incorrect length: %v", len(f.Id))
	}
}

// TestSetSize tests File length and size are valid.
func TestSetSize(t *testing.T) {
	data := []byte("hello")
	length := "5"
	size := "5 bytes"
	f := &File{Data: data}
	f.setSize()
	if got := f.Length; got != length {
		t.Fatalf("length = %q; want %q", got, length)
	}
	if got := f.Size; got != size {
		t.Fatalf("size = %q; want '%q'", got, size)
	}
}

// TestSetSum tests File hash is a valid sum.
func TestSetSum(t *testing.T) {
	f := &File{Data: []byte("hello, world!\n")}
	expected := "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc"
	f.setSum()
	if f.Sum != expected {
		t.Fatalf("sum = %q; want %q", f.Sum, expected)
	}
}

// TestSetType tests File type is set correctly.
func TestSetType(t *testing.T) {
	files := []struct {
		Name string
		Type string
	}{
		{"app.apk", "application/vnd.android.package-archive"},
		{"archive.zip", "application/zip"},
		{"data.json", "application/json"},
		{"data.xml", "application/xml"},
		{"document.pdf", "application/pdf"},
		{"image.gif", "image/gif"},
		{"image.jpg", "image/jpeg"},
		{"index.html", "text/html; charset=utf-8"},
		{"notes.txt", "text/plain; charset=utf-8"},
		{"picture.png", "image/png"},
		{"style.css", "text/css; charset=utf-8"},
		{"table.csv", "text/csv; charset=utf-8"},
		{"whatever.foo", "application/octet-stream"},
	}
	for _, tt := range files {
		f := &File{Name: tt.Name}
		f.setType()
		if f.Type != tt.Type {
			t.Errorf("setType(%q) = %q; want %q",
				tt.Name, f.Type, tt.Type)
		}
	}
}

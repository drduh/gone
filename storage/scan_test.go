package storage

import (
	"strings"
	"testing"
)

const fileIDLength = 44 // storageVersion "1" + 32 bytes

// TestFileScan tests all scan attributes are set.
func TestFileScan(t *testing.T) {
	tests := []struct {
		name   string
		file   File
		length string
		size   int
		sum    string
		typ    string
	}{
		{
			name: "text file",
			file: File{
				Name: "test.txt",
				Data: []byte("hello, world!\n"),
			},
			length: "14",
			size:   14,
			sum:    "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc",
			typ:    "text/plain; charset=utf-8",
		},
		{
			name: "empty binary file",
			file: File{
				Name: "empty.bin",
				Data: []byte{},
			},
			length: "0",
			size:   0,
			sum:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			typ:    "application/octet-stream",
		},
		{
			name: "html file",
			file: File{
				Name: "page.html",
				Data: []byte("<!DOCTYPE html><html><head></head><body></body></html>"),
			},
			length: "54",
			size:   54,
			sum:    "a68ced9bf3600aac812a146702e5e951eadb5e0312a8eaa7b54484534c979067",
			typ:    "text/html; charset=utf-8",
		},
		{
			name: "png file",
			file: File{
				Name: "image.png",
				Data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // rfc2083
			},
			length: "8",
			size:   8,
			sum:    "4c4b6a3be1314ab86138bef4314dde022e600960d8689a2c8f8631802d20dab6",
			typ:    "image/png",
		},
		{
			name: "pdf file",
			file: File{
				Name: "document.pdf",
				Data: []byte("%PDF-1.4\n"), // rfc8118
			},
			length: "9",
			size:   9,
			sum:    "e5c62df5dab5c87b6a015ef3d43597074d1eec433b15f51aec63b8582d0e4ab4",
			typ:    "application/pdf",
		},
		{
			name: "apk file override",
			file: File{
				Name: "app.apk",
				Data: []byte("not an apk"),
			},
			length: "10",
			size:   10,
			sum:    "25d5e232f8362e7f9ed050b7c6c5c2e6ae378ef7280843e09ac7a379dbad7649",
			typ:    "application/vnd.android.package-archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			f.Scan()

			if len(f.ID) != fileIDLength {
				t.Fatalf("id length = %d; want %d",
					len(f.ID), fileIDLength)
			}
			if f.Bytes != tt.size {
				t.Fatalf("bytes = %d; want %d",
					f.Bytes, tt.size)
			}
			if f.Length != tt.length {
				t.Fatalf("length = %q; want %q",
					f.Length, tt.length)
			}
			if f.Sum != tt.sum {
				t.Fatalf("sum = %q; want %q",
					f.Sum, tt.sum)
			}
			if f.Type != tt.typ {
				t.Fatalf("type = %q; want %q",
					f.Type, tt.typ)
			}
		})
	}
}

// TestSetID tests File id has valid length.
func TestSetID(t *testing.T) {
	f := &File{}
	f.setID()
	if f.ID == "" {
		t.Fatalf("id is empty")
	}
	if len(f.ID) != fileIDLength {
		t.Fatalf("id length is incorrect: %v",
			len(f.ID))
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
		t.Fatalf("length = %q; want %q",
			got, length)
	}
	if got := f.Size; got != size {
		t.Fatalf("size = %q; want '%q'",
			got, size)
	}
}

// TestSetSum tests File hash is a valid sum.
func TestSetSum(t *testing.T) {
	f := &File{Data: []byte("hello, world!\n")}
	expected := "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc"
	f.setSum()
	if f.Sum != expected {
		t.Fatalf("sum = %q; want %q",
			f.Sum, expected)
	}
}

// TestSetTypeEmptyFile tests File type assignment
// for empty files.
func TestSetTypeEmptyFile(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"app.apk", "application/vnd.android.package-archive"},
		{"archive.zip", "application/zip"},
		{"data.json", "application/json"},
		{"document.pdf", "application/pdf"},
		{"image.gif", "image/gif"},
		{"image.jpg", "image/jpeg"},
		{"index.html", "text/html"},
		{"notes.txt", "text/plain"},
		{"picture.png", "image/png"},
		{"style.css", "text/css"},
		{"table.csv", "text/csv"},
		{"whatever.foo", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{Name: tt.name}
			f.setType()

			got := f.Type
			if i := strings.Index(got, ";"); i != -1 {
				got = got[:i]
			}

			if got != tt.want {
				t.Fatalf("setType(%q) = %q; want %q",
					tt.name, got, tt.want)
			}
		})
	}
}

// TestSetTypeOverrides tests extension-based override.
func TestSetTypeOverrides(t *testing.T) {
	tests := []struct {
		name string
		file File
		want string
	}{
		{
			name: "apk empty file",
			file: File{Name: "app.apk"},
			want: "application/vnd.android.package-archive",
		},
		{
			name: "apk html data",
			file: File{
				Name: "app.apk",
				Data: []byte("<!DOCTYPE html><html><body></body></html>"),
			},
			want: "application/vnd.android.package-archive",
		},
		{
			name: "apk case insensitive",
			file: File{Name: "APP.APK"},
			want: "application/vnd.android.package-archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			f.setType()

			if got := f.Type; got != tt.want {
				t.Fatalf("setType(%q) = %q; want %q",
					f.Name, got, tt.want)
			}
		})
	}
}

// TestSetTypeDetectContent tests content type detection
// based on File data only.
func TestSetTypeDetectContent(t *testing.T) {
	tests := []struct {
		name string
		file File
		want string
	}{
		{
			name: "html data",
			file: File{
				Name: "page.bin",
				Data: []byte("<!DOCTYPE html><html><head></head><body></body></html>"),
			},
			want: "text/html; charset=utf-8",
		},
		{
			name: "png data",
			file: File{
				Name: "image.bin",
				Data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // rfc2083
			},
			want: "image/png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			f.setType()
			if f.Type != tt.want {
				t.Fatalf("setType(%q) = %q; want %q",
					f.Name, f.Type, tt.want)
			}
		})
	}
}

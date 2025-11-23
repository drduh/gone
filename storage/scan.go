package storage

import (
	"mime"
	"path/filepath"
	"strconv"

	"github.com/drduh/gone/util"
)

// Scan identifies and sets File attributes.
func (f *File) Scan() {
	f.setId()
	f.setSize()
	f.setSum()
	f.setType()
}

// SetId sets versioned File id with encoded entropy bytes.
func (f *File) setId() {
	f.Id = storageVersion + util.RandomId()
}

// SetSize sets the content length and readable file size.
func (f *File) setSize() {
	size := len(f.Data)
	f.Length = strconv.Itoa(size)
	f.Size = util.FormatSize(size)
}

// SetSum sets the content SHA-256 hash sum.
func (f *File) setSum() {
	f.Sum = util.Sum(f.Data)
}

// SetType sets File content type based on filename extension.
func (f *File) setType() {
	t := mime.TypeByExtension(filepath.Ext(f.Name))
	if t == "" {
		t = "application/octet-stream"
	}
	f.Type = t
}

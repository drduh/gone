package storage

import (
	"mime"
	"net/http"
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

// SetId sets a versioned File id using encoded
// random bytes.
func (f *File) setId() {
	f.Id = storageVersion + util.RandomId()
}

// SetSize sets File byte count, content length and
// and formatted size.
func (f *File) setSize() {
	f.Bytes = len(f.Data)
	f.Length = strconv.Itoa(f.Bytes)
	f.Size = util.FormatSize(f.Bytes)
}

// SetSum sets the content hash sum.
func (f *File) setSum() {
	f.Sum = util.Sum(f.Data)
}

// SetType sets File content type based on content type,
// or filename extension for empty files.
func (f *File) setType() {
	if len(f.Data) > 0 {
		f.Type = http.DetectContentType(f.Data)
	} else {
		ext := filepath.Ext(f.Name)
		if t := mime.TypeByExtension(ext); t != "" {
			f.Type = t
		}
	}
	if f.Type == "" {
		f.Type = "application/octet-stream"
	}
}

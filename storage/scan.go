package storage

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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
// formatted size.
func (f *File) setSize() {
	f.Bytes = len(f.Data)
	f.Length = strconv.Itoa(f.Bytes)
	f.Size = util.FormatSize(f.Bytes)
}

// SetSum sets the content hash sum.
func (f *File) setSum() {
	f.Sum = util.Sum(f.Data)
}

// SetType sets File content type based on extension
// override, contents, or filename extension.
func (f *File) setType() {
	const defaultType = "application/octet-stream"

	overrides := map[string]string{
		".apk": "application/vnd.android.package-archive",
	}

	ext := strings.ToLower(filepath.Ext(f.Name))

	if t, ok := overrides[ext]; ok {
		f.Type = t
		return
	}

	if len(f.Data) > 0 {
		f.Type = http.DetectContentType(f.Data)
		return
	}

	if t := mime.TypeByExtension(ext); t != "" {
		f.Type = t
		return
	}

	f.Type = defaultType
}

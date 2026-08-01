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
	f.setID()
	f.setSize()
	f.setSum()
	f.setType()
}

// SetID sets a versioned unique File ID using
// encoded random bytes.
func (f *File) setID() {
	f.ID = storageVersion + util.RandomID()
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
	ext := strings.ToLower(filepath.Ext(f.Name))

	if !f.setTypeOverride(ext) {
		if len(f.Data) > 0 {
			f.Type = http.DetectContentType(f.Data)
		} else if t := mime.TypeByExtension(ext); t != "" {
			f.Type = t
		} else {
			f.Type = "application/octet-stream"
		}
	}

	f.setTypeFmt()
}

// setTypeOverride sets Type when an extension-specific
// override is available.
func (f *File) setTypeOverride(ext string) bool {
	overrides := map[string]string{
		".apk": "application/vnd.android.package-archive",
	}

	t, ok := overrides[ext]
	if !ok {
		return false
	}

	f.Type = t
	return true
}

// setTypeFmt sets TypeFmt to Type unless a type-specific
// format override is available.
func (f *File) setTypeFmt() {
	overrides := map[string]string{
		"application/vnd.android.package-archive": "android package",
		"application/zip": "zip archive",
		"text/html; charset=utf-8": "html document",
		"text/plain; charset=utf-8": "text file",
	}

	if t, ok := overrides[f.Type]; ok {
		f.TypeFmt = t
		return
	}

	f.TypeFmt = f.Type
}

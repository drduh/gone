package storage

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
)

// Serve writes a File as an HTTP response. Downloads are
// counted by ReserveDownload, not here.
func (f *File) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition",
		mime.FormatMediaType(
			"attachment", map[string]string{"filename": f.Name}))
	w.Header().Set("Content-Length", f.Length)
	w.Header().Set("Content-Type", f.Type)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	_, _ = w.Write(f.Data)
}

// ServeMessages writes all Messages as a text file.
// Content is built under the read lock, then written.
func (s *Storage) ServeMessages(w http.ResponseWriter) {
	const msgFmt = "%d (%s) - %s\n"

	var buf bytes.Buffer
	s.Lease.RLock()
	for _, msg := range s.Messages {
		fmt.Fprintf(&buf, msgFmt,
			msg.Count, msg.UploadTimeFmt, msg.Data)
	}
	s.Lease.RUnlock()

	w.Header().Set("Content-Disposition",
		`attachment; filename="`+filenameMsgs+`"`)
	w.Header().Set("Content-Type", "text/plain")

	_, _ = w.Write(buf.Bytes())
}

// ServeWall writes Wall content as a text file.
func (s *Storage) ServeWall(w http.ResponseWriter) {
	content, _ := s.GetWall()

	w.Header().Set("Content-Disposition",
		`attachment; filename="`+filenameWall+`"`)
	w.Header().Set("Content-Type", "text/plain")

	_, _ = fmt.Fprint(w, content)
}

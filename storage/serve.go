package storage

import (
	"fmt"
	"mime"
	"net/http"
)

// Serve writes a File as an HTTP response.
func (f *File) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition",
		mime.FormatMediaType(
			"attachment", map[string]string{"filename": f.Name}))
	w.Header().Set("Content-Length", f.Length)
	w.Header().Set("Content-Type", f.Type)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	n, err := w.Write(f.Data)
	if err != nil {
		return
	}

	if n == f.Bytes {
		f.Count++
	}
}

// ServeMessages writes all Messages as a text file.
func (s *Storage) ServeMessages(w http.ResponseWriter) {
	const msgFmt = "%d (%s) - %s\n"

	w.Header().Set("Content-Disposition",
		`attachment; filename="`+filenameMsgs+`"`)
	w.Header().Set("Content-Type", "text/plain")

	for _, msg := range s.Messages {
		_, err := fmt.Fprintf(
			w, msgFmt, msg.Count, msg.UploadFmt, msg.Data)
		if err != nil {
			return
		}
	}
}

// ServeWall writes all Wall content as a text file.
func (s *Storage) ServeWall(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition",
		`attachment; filename="`+filenameWall+`"`)
	w.Header().Set("Content-Type", "text/plain")

	_, err := fmt.Fprintf(w, "%s", s.WallContent)
	if err != nil {
		return
	}
}

package storage

import (
	"fmt"
	"net/http"
)

// Serve writes a File as an HTTP response.
func (f *File) Serve(w http.ResponseWriter) {
	disposition := "attachment; filename=\"" + f.Name + "\""
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Length", f.Length)
	w.Header().Set("Content-Type", f.Type)
	n, err := w.Write(f.Data)
	if err != nil {
		return
	}
	if n == f.Bytes {
		f.Total++
	}
}

// ServeMessages writes all Messages as a text file.
func (s *Storage) ServeMessages(w http.ResponseWriter) {
	disposition := "attachment; filename=\"messages.txt\""
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", "text/plain")
	msgFormat := "%d (%s) - %s\n"
	for _, msg := range s.Messages {
		_, err := fmt.Fprintf(w, msgFormat, msg.Count, msg.Allow, msg.Data)
		if err != nil {
			return
		}
	}
}

// ServeWall writes all Wall content as a text file.
func (s *Storage) ServeWall(w http.ResponseWriter) {
	disposition := "attachment; filename=\"wall.txt\""
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintf(w, "%s", s.WallContent)
	if err != nil {
	    return
	}
}

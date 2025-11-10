package storage

import (
	"fmt"
	"net/http"
)

// Serve writes File as an HTTP response.
func (f *File) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name)
	w.Header().Set("Content-Length", f.Length)
	w.Header().Set("Content-Type", f.Type)
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(f.Data)
	if err != nil {
		return
	}
	if n == len(f.Data) {
		f.Total++
	}
}

// ServeMessages writes all Messages as an HTTP response.
func (s *Storage) ServeMessages(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=messages.txt")
	msgFormat := "%d (%s) - %s\n"
	for _, msg := range s.Messages {
		_, err := fmt.Fprintf(w, msgFormat, msg.Count, msg.Allow, msg.Data)
		if err != nil {
			return
		}
	}
}

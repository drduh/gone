package storage

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

// ClearStorage removes all Files and Messages from Storage.
func (s *Storage) ClearStorage() {
	s.ClearFiles()
	s.ClearMessages()
	s.ClearWall()
}

// ClearFiles removes all Files from Storage.
func (s *Storage) ClearFiles() {
	s.Files = make(map[string]*File)
}

// ClearMessages removes all Messages from Storage.
func (s *Storage) ClearMessages() {
	s.Messages = make(map[int]*Message)
}

// ClearWall removes all Wall content from Storage.
func (s *Storage) ClearWall() {
	s.WallContent = ""
}

// CountFiles returns the number of Files in Storage.
func (s *Storage) CountFiles() int {
	return len(s.Files)
}

// CountMessages returns the number of Messages in Storage.
func (s *Storage) CountMessages() int {
	return len(s.Messages)
}

// Expire removes a File from Storage.
func (s *Storage) Expire(f *File) {
	delete(s.Files, f.Name)
}

// IsExpires returns a reason if the File is expired.
func (f *File) IsExpired() string {
	if f.Total >= f.Downloads.Allow {
		return "limit downloads"
	}
	if f.GetLifetime() > f.Duration {
		return "limit duration"
	}
	return ""
}

// GetLifetime returns the duration since a File was uploaded.
func (f *File) GetLifetime() time.Duration {
	return time.Since(f.Time.Upload).Round(time.Second)
}

// GetType returns the File content type based on extension.
func (f *File) GetType() string {
	t := mime.TypeByExtension(filepath.Ext(f.Name))
	if t == "" {
		t = "application/octet-stream"
	}
	return t
}

// NumRemaining returns the number of downloads remaining
// until File expiration.
func (f *File) NumRemaining() int {
	return f.Downloads.Allow - f.Total
}

// TimeRemaining returns the relative duration remaining unitl
// until File expiration.
func (f *File) TimeRemaining() time.Duration {
	return time.Until(
		f.Time.Upload.Add(f.Time.Duration)).Round(time.Second)
}

// FindFile returns the requested File, if found by name.
func (s *Storage) FindFile(name string) *File {
	var file *File
	for _, f := range s.Files {
		if f.Name == name {
			file = f
			break
		}
	}
	return file
}

// Serve writes File as an HTTP response.
func (f *File) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Type", f.GetType())
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(f.Data)
	f.Total++
}

// ServeMessages writes all Messages as an HTTP response.
func (s *Storage) ServeMessages(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=messages.txt")
	for _, msg := range s.Messages {
		_, err := fmt.Fprintf(w, "%d. %s\n", msg.Count, msg.Data)
		if err != nil {
			return
		}
	}
}

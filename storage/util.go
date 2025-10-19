package storage

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/drduh/gone/util"
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

// GetSize sets content length and readable file size.
func (f *File) GetSize() {
	size := len(f.Data)
	f.Length = strconv.Itoa(size)
	f.Size = util.FormatSize(size)
}

// GetType sets File content type based on filename extension.
func (f *File) GetType() {
	t := mime.TypeByExtension(filepath.Ext(f.Name))
	if t == "" {
		t = "application/octet-stream"
	}
	f.Type = t
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

// UpdateTime updates time until expiration of each File in Storage.
func (s *Storage) UpdateTime() {
	for _, file := range s.Files {
		file.Time.Remain = file.TimeRemaining().String()
		s.Files[file.Name] = file
	}
}

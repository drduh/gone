package storage

import "time"

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

// FindFile returns a File found by id or name.
func (s *Storage) FindFile(id string) *File {
	var file *File
	for _, f := range s.Files {
		if f.Id == id || f.Name == id {
			file = f
			break
		}
	}
	return file
}

// UpdateTime updates time until expiration of each File in Storage.
func (s *Storage) UpdateTime() {
	for _, file := range s.Files {
		file.Time.Remain = file.TimeRemaining().String()
		s.Files[file.Name] = file
	}
}

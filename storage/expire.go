package storage

import "time"

// GetLifetime returns the rounded duration since a
// File was uploaded.
func (f *File) GetLifetime() time.Duration {
	return time.Since(f.Time.Upload).Round(time.Second)
}

// IsExpires returns a reason if File is expired,
// or an empty string if File is not expired.
func (f *File) IsExpired() string {
	if f.Total >= f.Downloads.Allow {
		return "limit downloads"
	}
	if f.GetLifetime() > f.Duration {
		return "limit duration"
	}
	return ""
}

// NumRemaining returns the number of downloads remaining
// until File expiration.
func (f *File) NumRemaining() int {
	return f.Downloads.Allow - f.Total
}

// TimeRemaining returns the relative duration remaining
// until File expiration.
func (f *File) TimeRemaining() time.Duration {
	return time.Until(
		f.Time.Upload.Add(f.Time.Duration)).Round(time.Second)
}

// Expire removes a File from Storage by id.
func (s *Storage) Expire(f *File) {
	delete(s.Files, f.Id)
}

// UpdateTimeRemaining updates the time remaining until
// expiration of each File in Storage.
func (s *Storage) UpdateTimeRemaining() {
	for _, file := range s.Files {
		file.Time.Remain = file.TimeRemaining().String()
		s.Files[file.Id] = file
	}
}

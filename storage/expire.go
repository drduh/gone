package storage

import "time"

// Expire removes a File from Storage by ID.
func (s *Storage) Expire(f *File) {
	delete(s.Files, f.ID)
}

// GetLifetime returns the duration since a File
// was uploaded, rounded to the nearest second.
func (f *File) GetLifetime() time.Duration {
	return time.Since(f.Time.Upload).Round(time.Second)
}

// IsExpired returns a reason the File is expired,
// or an empty string when the File isn't expired.
func (f *File) IsExpired() string {
	if f.Allow > 0 && f.Count >= f.Allow {
		return "limit downloads"
	}
	if f.Duration > 0 &&
		f.GetLifetime() > f.Duration {
		return "limit duration"
	}
	return ""
}

// TimeRemaining returns the relative duration remaining
// until File expiration.
func (f *File) TimeRemaining() time.Duration {
	return time.Until(
		f.Time.Upload.Add(f.Time.Duration)).Round(time.Second)
}

// UpdateRemainingDownloads sets the number of downloads
// remaining until File expiration.
func (f *File) UpdateRemainingDownloads() {
	f.Remain = f.Allow - f.Count
}

// UpdateRemainingFileLimits updates remaining limits
// of each File in Storage.
func (s *Storage) UpdateRemainingFileLimits() {
	for _, file := range s.Files {
		file.UpdateRemainingDownloads()
		file.DurationFmt = file.TimeRemaining().String()
	}
}

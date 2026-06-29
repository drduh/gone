package storage

import "time"

// Lifetime returns the duration since a File
// was uploaded, rounded to the nearest second.
func (f *File) Lifetime() time.Duration {
	return time.Since(f.Time.UploadTime).Round(time.Second)
}

// SetRemainingDuration sets the remaining duration
// until File expiration, or to "0s" when expired.
func (f *File) SetRemainingDuration() {
	if r := f.Duration - f.Lifetime(); r > 0 {
		f.DurationRemaining = r.Round(time.Second).String()
	} else {
		f.DurationRemaining = "0s"
	}
}

// SetRemainingDownloads sets the number of downloads
// remaining for a File.
func (f *File) SetRemainingDownloads() {
	f.Remain = f.Allow - f.Count
}

// IsExpired returns a reason the File is expired,
// or an empty string when the File isn't expired.
func (f *File) IsExpired() string {
	if f.Allow > 0 && f.Count >= f.Allow {
		return "limit downloads"
	}
	if f.Duration > 0 && f.Lifetime() > f.Duration {
		return "limit duration"
	}
	return ""
}

// Expire removes a File from Storage by ID.
func (s *Storage) Expire(f *File) {
	delete(s.Files, f.ID)
}

// UpdateRemainingFileLimits updates remaining limits
// of each File in Storage.
func (s *Storage) UpdateRemainingFileLimits() {
	for _, file := range s.Files {
		file.SetRemainingDownloads()
		file.SetRemainingDuration()
	}
}

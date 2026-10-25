package storage

import (
	"slices"
	"strings"
)

// ListFiles returns non-expired Files in Storage, oldest
// first, removing expired Files in the process. Owner
// address and content are omitted.
func (s *Storage) ListFiles() []File {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	files := make([]File, 0, len(s.Files))
	for _, file := range s.sorted() {
		file.SetRemainingDownloads()
		file.SetRemainingDuration()

		if file.IsExpired() != "" {
			s.expire(file)
			continue
		}

		f := File{
			ID:                file.ID,
			Name:              file.Name,
			Size:              file.Size,
			Sum:               file.Sum,
			Type:              file.Type,
			Agent:             file.Agent,
			Mask:              file.Mask,
			DurationRemaining: file.DurationRemaining,
			UploadTimeFmt:     file.UploadTimeFmt,
			Allow:             file.Allow,
			Count:             file.Count,
			Remain:            file.Remain,
		}

		files = append(files, f)
	}

	return files
}

// sorted returns Files ordered by upload time, then id.
func (s *Storage) sorted() []*File {
	files := make([]*File, 0, len(s.Files))
	for _, f := range s.Files {
		files = append(files, f)
	}

	slices.SortFunc(files, func(a, b *File) int {
		if c := a.UploadTime.Compare(b.UploadTime); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})

	return files
}

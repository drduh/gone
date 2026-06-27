package storage

// ListFiles returns a list of non-expired Files in Storage,
// removing expired Files in the process.
func (s *Storage) ListFiles() []File {
	files := make([]File, 0, len(s.Files))
	for _, file := range s.Files {
		file.SetRemainingDownloads()
		file.SetRemainingDuration()

		if file.IsExpired() != "" {
			s.Expire(file)
			continue
		}

		f := File{
			ID:   file.ID,
			Name: file.Name,
			Size: file.Size,
			Sum:  file.Sum,
			Type: file.Type,
			Owner: Owner{
				Agent: file.Agent,
				Mask:  file.Mask,
			},
			Time: Time{
				DurationRemaining: file.DurationRemaining,
				UploadTimeFmt:     file.UploadTimeFmt,
			},
			Downloads: Downloads{
				Allow:  file.Allow,
				Count:  file.Count,
				Remain: file.Remain,
			},
		}

		files = append(files, f)
	}

	return files
}

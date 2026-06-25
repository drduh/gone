package storage

// ListFiles returns a list of non-expired Files in Storage,
// removing expired Files in the process.
func (s *Storage) ListFiles() []File {
	s.UpdateRemainingFileLimits()

	files := make([]File, 0, len(s.Files))
	for _, file := range s.Files {
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
				DurationFmt: file.DurationFmt,
				UploadFmt:   file.UploadFmt,
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

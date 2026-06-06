package storage

// ListFiles returns a list of non-expired Files in Storage,
// removing expired Files in the process.
func (s *Storage) ListFiles() []File {
	s.UpdateTimeRemaining()

	files := make([]File, 0, len(s.Files))
	for _, file := range s.Files {
		if file.IsExpired() != "" {
			s.Expire(file)
			continue
		}

		f := File{
			Id:   file.Id,
			Name: file.Name,
			Size: file.Size,
			Sum:  file.Sum,
			Type: file.Type,
			Owner: Owner{
				Agent: file.Agent,
				Mask:  file.Mask,
			},
			Time: Time{
				Remain: file.Time.Remain,
				Upload: file.Upload,
			},
			Downloads: Downloads{
				Allow:  file.Downloads.Allow,
				Remain: file.NumRemaining(),
				Count:  file.Count,
			},
		}
		files = append(files, f)
	}

	return files
}

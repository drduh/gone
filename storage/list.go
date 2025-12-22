package storage

// ListFiles returns a list of non-expired Files in Storage,
// removing expired Files in the process.
func (s *Storage) ListFiles() []File {
	files := make([]File, 0, len(s.Files))
	for _, file := range s.Files {
		reason := file.IsExpired()
		if reason != "" {
			s.Expire(file)
			break
		}

		s.UpdateTimeRemaining()
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
				Total:  file.Total,
			},
		}
		files = append(files, f)
	}

	return files
}

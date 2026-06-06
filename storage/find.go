package storage

// FindFile returns a File found in Storage by ID or name.
func (s *Storage) FindFile(query string) *File {
	var file *File

	for _, f := range s.Files {
		if f.ID == query || f.Name == query {
			file = f
			break
		}
	}
	return file
}

package storage

// FindFile returns a File found in Storage by id or name.
func (s *Storage) FindFile(query string) *File {
	var file *File
	for _, f := range s.Files {
		if f.Id == query || f.Name == query {
			file = f
			break
		}
	}
	return file
}

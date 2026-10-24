package storage

// FindFile returns a File from Storage by id (map key),
// falling back to a name scan; nil if not found.
func (s *Storage) FindFile(query string) *File {
	if f, ok := s.Files[query]; ok {
		return f
	}

	for _, f := range s.Files {
		if f.ID == query || f.Name == query {
			return f
		}
	}
	return nil
}

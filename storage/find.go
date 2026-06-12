package storage

// FindFile returns a File from Storage if found
// by id or name, or nil if the File isn't found.
func (s *Storage) FindFile(query string) *File {
	for _, f := range s.Files {
		if f.ID == query || f.Name == query {
			return f
		}
	}
	return nil
}

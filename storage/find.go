package storage

// FindFile returns a copy of the File matching query by id
// (map key) or name; nil if not found. Names are not
// unique, so prefer ids.
func (s *Storage) FindFile(query string) *File {
	s.Lease.RLock()
	defer s.Lease.RUnlock()

	f := s.findFile(query)
	if f == nil {
		return nil
	}

	c := *f
	return &c
}

func (s *Storage) findFile(query string) *File {
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

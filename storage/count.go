package storage

// CountFiles returns the number of Files in Storage.
func (s *Storage) CountFiles() int {
	return len(s.Files)
}

// CountMessages returns the number of Messages in Storage.
func (s *Storage) CountMessages() int {
	return len(s.Messages)
}

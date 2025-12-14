package storage

// CountStorage performs all Storage counts.
func (s *Storage) CountStorage() {
	s.CountFiles()
	s.CountMessages()
	s.CountWall()
}

// CountFiles counts the number of Files in Storage.
func (s *Storage) CountFiles() {
	s.NumFiles = len(s.Files)
}

// CountMessages counts the number of Messages in Storage.
func (s *Storage) CountMessages() {
	s.NumMessages = len(s.Messages)
}

// CountWall counts the length of Wall contents in Storage.
func (s *Storage) CountWall() {
	s.CharsWall = len(s.WallContent)
}

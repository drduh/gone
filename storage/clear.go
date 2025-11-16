package storage

// ClearStorage removes all Files and Messages from Storage.
func (s *Storage) ClearStorage() {
	s.ClearFiles()
	s.ClearMessages()
	s.ClearWall()
}

// ClearFiles removes all Files from Storage.
func (s *Storage) ClearFiles() {
	s.Files = make(map[string]*File)
}

// ClearMessages removes all Messages from Storage.
func (s *Storage) ClearMessages() {
	s.Messages = make(map[int]*Message)
}

// ClearWall removes all Wall content from Storage.
func (s *Storage) ClearWall() {
	s.WallContent = ""
}

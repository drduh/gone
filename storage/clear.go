package storage

// ClearStorage removes all content from Storage.
func (s *Storage) ClearStorage() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.clearFiles()
	s.clearMessages()
	s.clearWall()
}

// ClearFiles removes all Files from Storage.
func (s *Storage) ClearFiles() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.clearFiles()
}

// ClearMessages removes all Messages from Storage.
func (s *Storage) ClearMessages() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.clearMessages()
}

// ClearWall removes all Wall content from Storage.
func (s *Storage) ClearWall() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.clearWall()
}

func (s *Storage) clearFiles() {
	s.Files = make(map[string]*File)
	s.countFiles()
}

func (s *Storage) clearMessages() {
	s.Messages = []*Message{}
	s.countMessages()
}

func (s *Storage) clearWall() {
	s.WallContent = ""
	s.countWall()
	s.WallCap = 0
}

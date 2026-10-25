package storage

// SetWall replaces Wall content and its modification
// metadata, then recounts chars, lines and capacity
// against limit. Returns the updated metadata.
func (s *Storage) SetWall(content string, modified WallMeta,
	limit int) WallMeta {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.WallContent = content
	s.WallModifiedTime = modified.WallModifiedTime
	s.WallModifiedTimeFmt = modified.WallModifiedTimeFmt
	s.WallModifiedUser = modified.WallModifiedUser
	s.countWall()
	s.wallCap(limit)

	return s.WallMeta
}

// GetWall returns Wall content and metadata.
func (s *Storage) GetWall() (string, WallMeta) {
	s.Lease.RLock()
	defer s.Lease.RUnlock()

	return s.WallContent, s.WallMeta
}

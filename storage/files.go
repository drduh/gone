package storage

// AddFiles stores files if their combined size plus Files
// already stored is within maxTotal (0 for no limit).
// All or nothing, so concurrent uploads cannot overshoot.
func (s *Storage) AddFiles(files []*File, maxTotal int64) bool {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	if s.Files == nil {
		s.Files = make(map[string]*File)
	}

	var incoming int64
	for _, f := range files {
		incoming += int64(f.Bytes)
	}

	if maxTotal > 0 &&
		int64(s.filesBytes())+incoming > maxTotal {
		return false
	}

	for _, f := range files {
		s.Files[f.ID] = f
	}
	s.countFiles()

	return true
}

// ReserveDownload finds a File by id or name, counts one
// download and removes the File once a limit is met.
// Returns a copy to serve (Data is immutable once stored),
// or nil when not found or already expired. Reason is set
// when this call removed the File. Counting before the
// response is written means concurrent requests cannot
// both claim the last download.
func (s *Storage) ReserveDownload(query string) (*File, ExpiryReason) {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	f := s.findFile(query)
	if f == nil {
		return nil, expiryNone
	}

	if reason := f.IsExpired(); reason != "" {
		s.expire(f)
		return nil, reason
	}

	f.Count++
	f.SetRemainingDownloads()
	f.SetRemainingDuration()

	reason := f.IsExpired()
	if reason != "" {
		s.expire(f)
	}
	s.countFiles()

	c := *f
	return &c, reason
}

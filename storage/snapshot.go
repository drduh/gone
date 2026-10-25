package storage

// Snapshot returns a detached copy of Storage for
// rendering, with remaining limits and counts refreshed.
func (s *Storage) Snapshot() Snapshot {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.updateRemaining()
	s.countStorage()

	snap := Snapshot{
		Sizes:       s.Sizes,
		WallMeta:    s.WallMeta,
		WallContent: s.WallContent,
		Files:       make([]File, 0, len(s.Files)),
		Messages:    make([]Message, 0, len(s.Messages)),
	}

	for _, f := range s.sorted() {
		snap.Files = append(snap.Files, f.meta())
	}
	for _, m := range s.Messages {
		snap.Messages = append(snap.Messages, *m)
	}

	return snap
}

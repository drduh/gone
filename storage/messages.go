package storage

// AddMessage stores m numbered after the last Message,
// unless maxCount (0 for no limit) is already reached.
func (s *Storage) AddMessage(m *Message, maxCount int) bool {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	if maxCount > 0 && len(s.Messages) >= maxCount {
		return false
	}

	m.Count = len(s.Messages) + 1
	s.Messages = append(s.Messages, m)
	s.countMessages()

	return true
}

// GetMessage returns a copy of the Message numbered count.
func (s *Storage) GetMessage(count int) (Message, bool) {
	s.Lease.RLock()
	defer s.Lease.RUnlock()

	for _, m := range s.Messages {
		if m.Count == count {
			return *m, true
		}
	}

	return Message{}, false
}

// DeleteMessage removes the Message numbered count
// and renumbers the remainder.
func (s *Storage) DeleteMessage(count int) bool {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	for i, m := range s.Messages {
		if m.Count != count {
			continue
		}

		s.Messages = append(
			s.Messages[:i], s.Messages[i+1:]...)
		for i, remaining := range s.Messages {
			remaining.Count = i + 1
		}
		s.countMessages()

		return true
	}

	return false
}

// ListMessages returns copies of all Messages.
func (s *Storage) ListMessages() []Message {
	s.Lease.RLock()
	defer s.Lease.RUnlock()

	messages := make([]Message, 0, len(s.Messages))
	for _, m := range s.Messages {
		messages = append(messages, *m)
	}

	return messages
}

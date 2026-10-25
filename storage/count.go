package storage

import (
	"strings"

	"github.com/drduh/gone/util"
)

// CountStorage counts all Storage content.
func (s *Storage) CountStorage() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.countStorage()
}

// CountFiles counts the number of Files and
// their total combined size in bytes.
func (s *Storage) CountFiles() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.countFiles()
}

// CountMessages counts the number of Messages
// and total count of characters in all Messages.
func (s *Storage) CountMessages() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.countMessages()
}

// CountWall counts the number of characters
// and lines in Wall content.
func (s *Storage) CountWall() {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.countWall()
}

// GetWallCap determines the wall capacity
// percentage, rounded to the nearest int.
func (s *Storage) GetWallCap(limit int) {
	s.Lease.Lock()
	defer s.Lease.Unlock()

	s.wallCap(limit)
}

// FilesBytes returns the total size of stored Files.
func (s *Storage) FilesBytes() int {
	s.Lease.RLock()
	defer s.Lease.RUnlock()

	return s.filesBytes()
}

func (s *Storage) countStorage() {
	s.countFiles()
	s.countMessages()
	s.countWall()
}

func (s *Storage) filesBytes() int {
	total := 0
	for _, file := range s.Files {
		total += file.Bytes
	}

	return total
}

func (s *Storage) countFiles() {
	s.NumFiles = len(s.Files)
	s.SizeFiles = s.filesBytes()
	if s.SizeFiles > 0 {
		s.SizeFilesFmt = util.FormatSize(s.SizeFiles)
	} else {
		s.SizeFilesFmt = ""
	}
}

func (s *Storage) countMessages() {
	s.NumMessages = len(s.Messages)

	total := 0
	for _, message := range s.Messages {
		total += len(message.Data)
	}
	s.CharsMessages = total
}

func (s *Storage) countWall() {
	s.WallChars = len(s.WallContent)

	if s.WallContent == "" {
		s.WallLines = 0
	} else {
		s.WallLines = len(strings.Split(
			s.WallContent, "\n"))
	}
}

func (s *Storage) wallCap(limit int) {
	if limit <= 0 {
		s.WallCap = 0
		return
	}
	s.WallCap = s.WallChars * 100 / limit
}

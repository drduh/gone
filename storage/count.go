package storage

import (
	"strings"

	"github.com/drduh/gone/util"
)

// CountStorage counts all Storage content.
func (s *Storage) CountStorage() {
	s.CountFiles()
	s.CountMessages()
	s.CountWall()
}

// CountFiles counts the number of Files and
// their total combined size in bytes.
func (s *Storage) CountFiles() {
	s.NumFiles = len(s.Files)
	total := 0

	for _, file := range s.Files {
		total += file.Bytes
	}
	s.SizeFiles = total
	if s.SizeFiles > 0 {
		s.SizeFilesFmt = util.FormatSize(s.SizeFiles)
	} else {
		s.SizeFilesFmt = ""
	}
}

// CountMessages counts the number of Messages
// and total count of characters in all Messages.
func (s *Storage) CountMessages() {
	s.NumMessages = len(s.Messages)

	total := 0
	for _, message := range s.Messages {
		total += len(message.Data)
	}
	s.CharsMessages = total
}

// CountWall counts the number of characters
// and lines in Wall content.
func (s *Storage) CountWall() {
	s.WallChars = len(s.WallContent)

	if s.WallContent == "" {
		s.WallLines = 0
	} else {
		s.WallLines = len(strings.Split(
			s.WallContent, "\n"))
	}
}

// GetWallCap determines the wall capacity
// percentage, rounded to the nearest int.
func (s *Storage) GetWallCap(limit int) {
	s.WallCap = s.WallChars * 100 / limit
}

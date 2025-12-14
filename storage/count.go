package storage

import (
	"strings"

	"github.com/drduh/gone/util"
)

// CountStorage performs all Storage counts.
func (s *Storage) CountStorage() {
	s.CountFiles()
	s.CountMessages()
	s.CountWall()
}

// CountFiles counts the number of Files and
// their total combined size.
func (s *Storage) CountFiles() {
	s.NumFiles = len(s.Files)
	total := 0
	for _, file := range s.Files {
		total += file.Bytes
	}
	s.SizeFiles = total
	s.SizeFilesFmt = util.FormatSize(s.SizeFiles)
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
// and lines in Wall contents.
func (s *Storage) CountWall() {
	s.CharsWall = len(s.WallContent)
	if s.WallContent == "" {
		s.LinesWall = 0
	} else {
		s.LinesWall = len(
			strings.Split(s.WallContent, "\n"))
	}
}

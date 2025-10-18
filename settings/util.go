package settings

import "fmt"

// GetAddr returns the server address based on configured port.
func (s *Settings) GetAddr() string {
	return fmt.Sprintf(":%d", s.Port)
}

// GetMaxBytes returns the maximum allowed file size in bytes.
func (l *Limit) GetMaxBytes() int64 {
	return l.MaxSizeMb << 20
}

package settings

import "fmt"

// GetAddr returns the server address based on configured port.
func (s *Settings) GetAddr() string {
	return fmt.Sprintf(":%d", s.Port)
}

// GetMaxFileBytes returns the maximum allowed file size in bytes.
func (l *Limit) GetMaxFileBytes() int64 {
	return l.MaxSizeFileMb << 20
}

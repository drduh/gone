package settings

import "fmt"

// GetAddr returns the server address based on configured settings.
func (s *Settings) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.ServerAddr, s.ServerPort)
}

// GetMaxFileBytes returns the maximum allowed file size in bytes.
func (l *Limit) GetMaxFileBytes() int64 {
	return l.FileLimits.SizeEachMb << 20
}

// GetMaxTotalFilesBytes returns the maximum allowed size of all
// files in bytes.
func (l *Limit) GetMaxTotalFilesBytes() int64 {
	return l.FileLimits.SizeTotalMb << 20
}

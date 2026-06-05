package settings

import (
	"net"
	"strconv"
)

// GetAddr returns the configured server host
// and port as an address string.
func (s *Settings) GetAddr() string {
	return net.JoinHostPort(
		s.ServerAddr, strconv.Itoa(s.ServerPort))
}

// GetMaxFileBytes returns the configured
// per-file size limit in bytes.
func (l *Limit) GetMaxFileBytes() int64 {
	return l.FileLimits.SizeEachMb << 20
}

// GetMaxTotalFilesBytes returns the configured
// total file size limit in bytes.
func (l *Limit) GetMaxTotalFilesBytes() int64 {
	return l.FileLimits.SizeTotalMb << 20
}

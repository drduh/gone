package settings

import (
	"net"
	"strconv"
)

const bytesPerMB int64 = 1_000_000

// MegabytesToBytes converts decimal megabytes/MB
// to bytes (1 MB = 1,000,000 bytes).
func MegabytesToBytes(n int64) int64 {
	return n * bytesPerMB
}

// BytesToMegabytes converts bytes to decimal
// megabytes/MB (1 MB = 1,000,000 bytes).
func BytesToMegabytes(n int64) int64 {
	return n / bytesPerMB
}

// GetAddr returns the configured server host
// and port as an address string.
func (s *Settings) GetAddr() string {
	return net.JoinHostPort(
		s.ServerAddr, strconv.Itoa(s.ServerPort))
}

// GetMaxFileBytes returns the configured
// per-file size limit in bytes.
func (l *Limit) GetMaxFileBytes() int64 {
	return MegabytesToBytes(l.FileLimits.SizeEachMb)
}

// GetMaxTotalFilesBytes returns the configured
// total file size limit in bytes.
func (l *Limit) GetMaxTotalFilesBytes() int64 {
	return MegabytesToBytes(l.FileLimits.SizeTotalMb)
}

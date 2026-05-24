package util

import "fmt"

// FormatSize returns a formatted size from number of bytes.
func FormatSize(bytes int) string {
	var units = [...]string{
		"bytes", "kb", "mb", "gb", "tb",
	}

	if bytes <= 0 {
		return "0 bytes"
	}

	var unitIndex int
	size := float64(bytes)
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	if size == float64(int(size)) {
		return fmt.Sprintf("%d %s",
			int(size), units[unitIndex])
	}

	return fmt.Sprintf("%.2f %s",
		size, units[unitIndex])
}

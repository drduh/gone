package util

import "fmt"

// Returns formatted file size based on bytes
func FormatSize(bytes int) string {
	if bytes == 0 {
		return "0 Bytes"
	}
	units := []string{"Bytes", "KB", "MB", "GB"}
	var unitIndex int
	size := float64(bytes)
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

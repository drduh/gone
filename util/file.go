package util

import (
	"fmt"
	"io"
	"os"
)

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

// Returns file or stdout IO writer
func GetOutput(filename string) (io.Writer, error) {
	if filename == "" {
		return os.Stdout, nil
	}
	dest, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filename, err)
	}
	return dest, nil
}

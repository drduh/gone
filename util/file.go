package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var sizeUnits = []string{"bytes", "kb", "mb", "gb", "tb"}

// FormatSize returns a formatted size from number of bytes.
func FormatSize(bytes int) string {
	if bytes == 0 {
		return "0 bytes"
	}
	var unitIndex int
	size := float64(bytes)
	for size >= 1024 && unitIndex < len(sizeUnits)-1 {
		size /= 1024
		unitIndex++
	}
	if size == float64(int(size)) {
		return fmt.Sprintf("%d %s", int(size), sizeUnits[unitIndex])
	}
	return fmt.Sprintf("%.2f %s", size, sizeUnits[unitIndex])
}

// GetOutput returns the destination file or stdout IO writer.
func GetOutput(filename string) (io.Writer, error) {
	if filename == "" {
		return os.Stdout, nil
	}
	dest, err := os.OpenFile(
		filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %w", filename, err)
	}
	return dest, nil
}

// loadNames returns trimmed names from a file or
// the default names list.
func loadNames(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		return defaultNames
	}
	var fileNames []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			fileNames = append(fileNames, name)
		}
	}
	if err := f.Close(); err != nil {
		return defaultNames
	}
	if len(fileNames) == 0 {
		return defaultNames
	}
	return fileNames
}

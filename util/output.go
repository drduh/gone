package util

import (
	"fmt"
	"io"
	"os"
)

// GetOutput returns the destination file or stdout writer.
func GetOutput(filename string) (io.Writer, error) {
	if filename == "" {
		return os.Stdout, nil
	}

	dest, err := os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open '%s': %w", filename, err)
	}

	return dest, nil
}

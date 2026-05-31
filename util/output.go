package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GetOutput returns the destination file or stdout writer.
func GetOutput(p string) (io.Writer, error) {
	if p == "" {
		return os.Stdout, nil
	}

	dir := filepath.Dir(p)
	file := filepath.Base(p)

	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open root '%s': %w", dir, err)
	}
	defer func() { _ = root.Close() }()

	dest, err := root.OpenFile(file,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open file '%s': %w", file, err)
	}

	return dest, nil
}

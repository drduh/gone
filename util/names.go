package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// loadNames returns trimmed names from a file or
// the default names list.
func loadNames(dir, filename string) []string {
	f, err := os.OpenInRoot(dir, filename)
	if err != nil {
		return defaultNames
	}
	defer func() { _ = f.Close() }()

	return loadNamesFromReader(f)
}

func loadNamesFromReader(r io.Reader) []string {
	scanner := bufio.NewScanner(r)

	var names []string

	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names = append(names, name)
		}
	}

	if scanner.Err() != nil {
		return defaultNames
	}

	if len(names) == 0 {
		return defaultNames
	}

	return names
}

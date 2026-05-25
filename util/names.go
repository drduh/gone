package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var defaultNames = []string{
	"Alice", "Bob", "Charlie", "Diana", "Eve",
	"Frank", "Grace", "Henry", "Ivan", "Judy",
	"Ken", "Laura", "Mallory", "Nancy", "Olivia",
	"Peggy", "Quentin", "Rupert", "Sam", "Trent",
	"Uma", "Victor", "Wendy", "Xavier", "Yvonne", "Zack",
}

var names = loadNames("/etc/gone/assets", "names.txt")

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

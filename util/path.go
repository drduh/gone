package util

import "strings"

// GetBasePath returns the string up to and including the first "/".
func GetBasePath(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}

	i := strings.Index(s[1:], "/")
	if i == -1 {
		return s
	}

	return s[:i+2]
}

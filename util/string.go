package util

import (
	"strings"
	"unicode"
)

// GetBasePath returns the string up to and including the first "/".
func GetBasePath(s string) string {
	i := strings.Index(s[1:], "/")
	if i == -1 {
		return s
	}
	return s[:i+2]
}

// IsNumeric returns true if the string contains only numbers.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

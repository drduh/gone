package storage

import (
	"net/url"
	"path/filepath"
	"strings"
)

const (
	defaultName  = "default"
	maxExtLength = 5
)

// SanitizeName validates strings for use as filename.
func SanitizeName(input, extraChars string, maxLength int) string {
	if strings.TrimSpace(input) == "" {
		return defaultName
	}
	input, err := url.QueryUnescape(input)
	if err != nil {
		return defaultName
	}
	f := filepath.Base(input)
	f = removeInvalidChars(f, extraChars)
	ext := filepath.Ext(f)
	base := strings.TrimSuffix(f, ext)
	if strings.TrimSpace(base) == "" {
		base = defaultName
	}
	return truncateName(base, ext, maxLength)
}

// removeInvalidChars removes all invalid characters.
func removeInvalidChars(filename string, allowed string) string {
	var result strings.Builder
	for _, char := range filename {
		if (char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			isAllowedChar(char, allowed) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// truncateName trims a filename string to max size,
// preserving reasonably-sized original file extensions.
func truncateName(base string, ext string, maxLength int) string {
	ext = strings.ReplaceAll(ext, " ", "")
	if len(ext) > maxExtLength {
		ext = ext[:maxExtLength]
	}
	base = strings.TrimSpace(base)
	totalLength := len(base) + len(ext)
	if totalLength <= maxLength {
		return base + ext
	}
	allowedBaseLength := maxLength - len(ext)
	if allowedBaseLength > 0 {
		return base[:allowedBaseLength] + ext
	}
	return ext[:maxLength]
}

// isAllowedChar returns true if a character is allowed.
func isAllowedChar(r rune, allowed string) bool {
	return strings.ContainsRune(allowed, r)
}

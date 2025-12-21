package storage

import (
	"path/filepath"
	"strings"
	"unicode"
)

const defaultName = "default"

// SanitizeName validates strings for use as filename.
func SanitizeName(input string, maxLength int, extraChars string) string {
	f := filepath.Base(input)
	f = removeInvalidChars(f, extraChars)
	ext := filepath.Ext(f)
	base := strings.TrimSuffix(f, ext)
	if base == "" {
		base = defaultName
	}
	return truncateName(base, ext, maxLength)
}

// removeInvalidChars removes all invalid characters.
func removeInvalidChars(filename string, allowed string) string {
	var result strings.Builder
	for _, char := range filename {
		if unicode.IsDigit(char) ||
			unicode.IsLetter(char) ||
			isAllowedChar(char, allowed) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// truncateName trims a filename string to max size,
// preserving reasonably-sized original file extensions.
func truncateName(base string, ext string, maxLength int) string {
	const maxExtensionLength = 5
	if len(ext) > maxExtensionLength {
		ext = ext[:maxExtensionLength]
	}
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

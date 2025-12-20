package storage

import (
	"path/filepath"
	"strings"
	"unicode"
)

// SanitizeName validates strings for use as filename.
func SanitizeName(input string, maxLength int, specialChars string) string {
	f := filepath.Base(input)
	f = removeInvalidChars(f, specialChars)
	f = truncateName(f, maxLength)
	return f
}

// isAllowedChar returns true if a character is allowed.
func isAllowedChar(r rune, allowed string) bool {
	return strings.ContainsRune(allowed, r)
}

// removeInvalidChars removes all invalid characters.
func removeInvalidChars(filename string, allowed string) string {
	var result strings.Builder
	for _, char := range filename {
		if unicode.IsLetter(char) ||
			unicode.IsDigit(char) ||
			isAllowedChar(char, allowed) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// truncateName trims a filename string to max size,
// preserving the original file extension.
func truncateName(filename string, length int) string {
	if len(filename) <= length {
		return filename
	}
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	maxBaseLength := length - len(ext)
	if len(base) > maxBaseLength {
		base = base[:maxBaseLength]
	}
	return base + ext
}

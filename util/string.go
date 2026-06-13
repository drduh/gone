package util

import "unicode"

// upperFirst capitalizes the first letter of s.
func upperFirst(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// IsNumeric returns true if the string contains only numbers.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

package handlers

import (
	"unicode/utf8"

	"github.com/drduh/gone/util"
)

// getMask returns a masked address string.
func getMask(addr string) string {
	return util.GetMaskAddr(addr, false)
}

// refreshMask sets a new masked address.
func refreshMask(addr string) {
	util.GetMaskAddr(addr, true)
}

// charCount returns the number of characters in s.
func charCount(s string) int {
	return utf8.RuneCountInString(s)
}

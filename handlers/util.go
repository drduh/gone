package handlers

import "github.com/drduh/gone/util"

// getMask returns a masked address string.
func getMask(addr string) string {
	return util.GetMaskAddr(addr, false)
}

// refreshMask sets a new masked address.
func refreshMask(addr string) {
	util.GetMaskAddr(addr, true)
}

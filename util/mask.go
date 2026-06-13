package util

import (
	"net"
	"sync"
)

var masks sync.Map

// GetMask creates or loads a masked string, with
// an option to refresh an assigned mask.
func GetMask(s string, refresh bool) string {
	if !refresh {
		if name, ok := masks.Load(s); ok {
			str, ok := name.(string)
			if !ok {
				return "maskError"
			}

			return str
		}
	}
	mask := GetRandom("mask")
	masks.Store(s, mask)

	return mask
}

// GetMaskAddr creates or loads a masked addr string.
func GetMaskAddr(s string, refresh bool) string {
	addr, _, err := net.SplitHostPort(s)
	if err != nil {
		addr = "unknown"
	}

	return GetMask(addr, refresh)
}

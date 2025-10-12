package util

import "sync"

var masks sync.Map

// Mask creates or loads a replacement address string.
func Mask(address string) string {
	if name, ok := masks.Load(address); ok {
		return name.(string)
	}
	mask, _ := masks.LoadOrStore(address, GetRandom(""))
	return mask.(string)
}

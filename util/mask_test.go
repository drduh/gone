package util

import (
	"sync"
	"testing"
)

// TestMask tests assignment of unique address masks.
func TestMask(t *testing.T) {
	masks = sync.Map{}
	addr1 := "127.0.0.1"
	addr2 := "127.0.0.2"
	mask1 := Mask(addr1)
	mask2 := Mask(addr2)
	remask1 := Mask(addr1)
	if mask1 != remask1 {
		t.Errorf("%q should be %q", mask1, remask1)
	}
	if mask1 == mask2 {
		t.Errorf("masks should not be %q", mask1)
	}
}

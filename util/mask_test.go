package util

import "testing"

// TestMask tests assignment of unique address masks.
func TestGetMask(t *testing.T) {
	cases := []struct {
		addr         string
		addrWithPort string
	}{
		{"192.168.1.10", "192.168.1.10:8080"},
		{"192.168.1.8", "192.168.1.8:8080"},
		{"2001:db8::1", "[2001:db8::1]:8080"},
		{"127.0.0.1", "127.0.0.1:8080"},
		{"0.0.0.0", "0.0.0.0:8080"},
		{"255.255.255.255", "255.255.255.255:8080"},
		{"::1", "[::1]:8080"},
	}

	seen := make(map[string]string)
	for _, tc := range cases {
		mask1 := GetMask(tc.addr, false)
		if mask1 != GetMask(tc.addr, false) {
			t.Errorf("%q: expected same mask", tc.addr)
		}

		mask2 := GetMask(tc.addr, true)
		if mask1 == mask2 {
			t.Errorf("%q: expected new mask", tc.addr)
		}

		if GetMaskAddr(tc.addrWithPort, false) != mask2 {
			t.Errorf("%q: expected same mask", tc.addr)
		}

		for _, prev := range seen {
			if mask2 == prev {
				t.Errorf("%q: expected new mask, got %q",
					tc.addr, mask2)
			}
		}
		seen[tc.addr] = mask2
	}
}

// TestGetMaskInvalid tests assignment of masks to
// invalid input.
func TestGetMaskInvalid(t *testing.T) {
	cases := []string{
		"",
		"notanaddress",
		"127.o.O.1",
	}

	unknownMask := ""
	for _, addr := range cases {
		mask := GetMaskAddr(addr, false)
		if unknownMask == "" {
			unknownMask = mask
		} else if mask != unknownMask {
			t.Errorf("%q: expected same mask, got %q/%q",
				addr, mask, unknownMask)
		}
	}
}

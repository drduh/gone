package util

import (
	"fmt"
	"testing"
)

// TestFormatSize tests size integer conversion to readable string.
func TestFormatSize(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0 bytes"},
		{200, "200 bytes"},
		{1024, "1 kb"},
		{5000, "4.88 kb"},
		{1048576, "1 mb"},
		{5242880, "5 mb"},
		{100000000, "95.37 mb"},
		{50000 * 50000, "2.33 gb"},
		{50000 * 50000 * 50000, "113.69 tb"},
		{-1000, "0 bytes"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("input: %d", tt.input), func(t *testing.T) {
			if got := FormatSize(tt.input); got != tt.expect {
				t.Errorf("size %d: got %v, expected %v",
					tt.input, got, tt.expect)
			}
		})
	}
}

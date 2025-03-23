package util

import (
	"fmt"
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0 Bytes"},
		{200, "200.00 Bytes"},
		{1024, "1.00 KB"},
		{5000, "4.88 KB"},
		{1048576, "1.00 MB"},
		{5242880, "5.00 MB"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input: %d", tt.input), func(t *testing.T) {
			if got := FormatSize(tt.input); got != tt.expect {
				t.Errorf("%d=%v; expect %v", tt.input, got, tt.expect)
			}
		})
	}
}

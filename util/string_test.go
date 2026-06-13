package util

import "testing"

// TestUpperFirst tests first character is capitalized.
func TestUpperFirst(t *testing.T) {
	t.Parallel()

	if got := upperFirst(""); got != "" {
		t.Fatalf("got %q, want empty string", got)
	}

	if got := upperFirst("oak"); got != "Oak" {
		t.Fatalf("got %q, want %q", got, "Oak")
	}
}

// TestIsNumeric tests strings for numeric characters.
func TestIsNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input  string
		expect bool
	}{
		{"", false},
		{"10", true},
		{" 10", false},
		{"10s", false},
		{"10 1", false},
		{"10.1", false},
		{"001234", true},
		{"1000000000000", true},
		{"100000000000a", false},
		{"１２３", false},
		{"４2", false},
		{"²³", false},
		{"1_000", false},
		{"1,000", false},
		{"0x10", false},
		{"1e10", false},
		{"-10", false},
		{"+10", false},
		{"--10", false},
		{"10\t1", false},
		{"10\x00", false},
		{"\u200b10", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			if got := IsNumeric(tt.input); got != tt.expect {
				t.Errorf("%q = %v; expect %v",
					tt.input, got, tt.expect)
			}
		})
	}
}

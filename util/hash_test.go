package util

import "testing"

// TestSums tests Sum returns correct SHA-256 hash sums.
func TestSum(t *testing.T) {
	tests := []struct {
		name string
		hash string
		data []byte
	}{
		{
			name: "nil",
			hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			data: nil,
		},
		{
			name: "empty",
			hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			data: []byte(""),
		},
		{
			name: "hello",
			hash: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			data: []byte("hello"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.data)
			if got != tt.hash{
				t.Fatalf("Sum(%q) = %q; expect %q", tt.data, got, tt.hash)
			}
		})
	}
}

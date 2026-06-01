package util

import "testing"

// TestSums tests Sum returns correct SHA-256 hash sums.
func TestSum(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		hash string
	}{
		{
			name: "nil",
			data: nil,
			hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "empty",
			data: []byte(""),
			hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "string",
			data: []byte("hello, world!\n"),
			hash: "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.data)
			if got != tt.hash {
				t.Fatalf("Sum(%q) = %q; expect %q",
					tt.data, got, tt.hash)
			}
		})
	}
}

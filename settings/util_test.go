package settings

import "testing"

// TestGetAddr tests server addr and port format.
func TestGetAddr(t *testing.T) {
	tests := []struct {
		name string
		s    Settings
		want string
	}{
		{
			name: "empty host and port",
			s: Settings{
				ServerAddr: "",
				ServerPort: 8080,
			},
			want: ":8080",
		},
		{
			name: "hostname and port",
			s: Settings{
				ServerAddr: "localhost",
				ServerPort: 8080,
			},
			want: "localhost:8080",
		},
		{
			name: "ipv4 and port",
			s: Settings{
				ServerAddr: "127.0.0.1",
				ServerPort: 8080,
			},
			want: "127.0.0.1:8080",
		},
		{
			name: "ipv6 and port",
			s: Settings{
				ServerAddr: "::1",
				ServerPort: 8080,
			},
			want: "[::1]:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetAddr()
			if got != tt.want {
				t.Fatalf("GetAddr returned %q, want %q",
					got, tt.want)
			}
		})
	}
}

// TestMegabytes tests MB and byte conversions.
func TestMegabytes(t *testing.T) {
	if got := MegabytesToBytes(3); got != 3_000_000 {
		t.Fatalf("MegabytesToBytes(3) = %d", got)
	}
	if got := BytesToMegabytes(3_999_999); got != 3 {
		t.Fatalf("BytesToMegabytes = %d; want 3", got)
	}

	var l Limit
	l.FileLimits.SizeEachMb = 2
	l.FileLimits.SizeTotalMb = 5
	if got := l.GetMaxFileBytes(); got != 2_000_000 {
		t.Fatalf("GetMaxFileBytes = %d", got)
	}
	if got := l.GetMaxTotalFilesBytes(); got != 5_000_000 {
		t.Fatalf("GetMaxTotalFilesBytes = %d", got)
	}
}

package util

import "testing"

// TestGetBasePath tests strings for path trims.
func TestGetBasePath(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"", ""},
		{"/", "/"},
		{"foo", "foo"},
		{"/foo", "/foo"},
		{"/foo/bar", "/foo/"},
		{"foo/bar", "foo/"},
		{"foo/bar/zoo", "foo/"},
		{"foo///bar", "foo/"},
		{"foo/bar  ", "foo/"},
	}

	for _, tt := range tests {
		got := GetBasePath(tt.input)
		if got != tt.expect {
			t.Errorf("%q = %q; expect %q",
				tt.input, got, tt.expect)
		}
	}
}

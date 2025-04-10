package util

import "testing"

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"10", true},
		{"10s", false},
		{"10 1", false},
		{"10.1", false},
		{"00010", true},
		{"1000000000000", true},
		{"100000000000a", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsNumeric(tt.input); got != tt.expect {
				t.Errorf("%q=%v; expect %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestGetBasePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/", "/"},
		{"/foo", "/foo"},
		{"/foo/bar", "/foo/"},
		{"foo/bar", "foo/"},
		{"foo", "foo"},
		{"foo/bar/zoo", "foo/"},
		{"", ""},
		{"foo///bar", "foo/"},
		{"foo/bar  ", "foo/"},
	}

	for _, tt := range tests {
		got := GetBasePath(tt.input)
		if got != tt.expected {
			t.Errorf("%q=%q; expect %q", tt.input, got, tt.expected)
		}
	}
}

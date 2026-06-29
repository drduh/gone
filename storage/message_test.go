package storage

import (
	"reflect"
	"testing"
)

// TestMessageParts tests Message data is split into
// plain-text and URL parts.
func TestMessageParts(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []MessageParts
	}{
		{
			name: "no urls",
			data: "hello, world!",
			want: []MessageParts{
				{Text: "hello, world!"},
			},
		},
		{
			name: "one url",
			data: "visit https://example.com now",
			want: []MessageParts{
				{Text: "visit "},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: " now"},
			},
		},
		{
			name: "multiple urls",
			data: "a https://1.example.com b http://two.example.com/path c",
			want: []MessageParts{
				{Text: "a "},
				{Text: "https://1.example.com",
					URL: "https://1.example.com", HasURL: true},
				{Text: " b "},
				{Text: "http://two.example.com/path",
					URL: "http://two.example.com/path", HasURL: true},
				{Text: " c"},
			},
		},
		{
			name: "javascript scheme",
			data: "javascript:alert(1)",
			want: []MessageParts{
				{Text: "javascript:alert(1)"},
			},
		},
		{
			name: "data scheme",
			data: "data:text/html,<script>alert(1)</script>",
			want: []MessageParts{
				{Text: "data:text/html,<script>alert(1)</script>"},
			},
		},
		{
			name: "html script tag and url",
			data: "<script>alert(\"xss\")</script> https://example.com",
			want: []MessageParts{
				{Text: "<script>alert(\"xss\")</script> "},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
			},
		},
		{
			name: "mailto not linkified",
			data: "contact mailto:user@example.com",
			want: []MessageParts{
				{Text: "contact mailto:user@example.com"},
			},
		},
		{
			name: "ftp not linkified",
			data: "download ftp://example.com/file.txt",
			want: []MessageParts{
				{Text: "download ftp://example.com/file.txt"},
			},
		},
		{
			name: "url at start",
			data: "https://example.com is first",
			want: []MessageParts{
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: " is first"},
			},
		},
		{
			name: "url at end",
			data: "last is https://example.com",
			want: []MessageParts{
				{Text: "last is "},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
			},
		},
		{
			name: "invalid url",
			data: "broken https:///example.com link",
			want: []MessageParts{
				{Text: "broken https:///example.com link"},
			},
		},
		{
			name: "parentheses",
			data: "(https://example.com)",
			want: []MessageParts{
				{Text: "("},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: ")"},
			},
		},
		{
			name: "trailing period",
			data: "see https://example.com.",
			want: []MessageParts{
				{Text: "see "},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: "."},
			},
		},
		{
			name: "trailing comma",
			data: "https://example.com, and more",
			want: []MessageParts{
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: ", and more"},
			},
		},
		{
			name: "trailing punctuation",
			data: "is this it?! https://example.com?!",
			want: []MessageParts{
				{Text: "is this it?! "},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: "?!"},
			},
		},
		{
			name: "digit query string",
			data: "https://example.com/?page=2",
			want: []MessageParts{
				{Text: "https://example.com/?page=2",
					URL: "https://example.com/?page=2", HasURL: true},
			},
		},
		{
			name: "path and query",
			data: "https://example.com/foo/bar?x=1&y=2",
			want: []MessageParts{
				{Text: "https://example.com/foo/bar?x=1&y=2",
					URL: "https://example.com/foo/bar?x=1&y=2", HasURL: true},
			},
		},
		{
			name: "only trailing dot",
			data: "http://.",
			want: []MessageParts{
				{Text: "http://."},
			},
		},
		{
			name: "trailing slash",
			data: "https://example.com/path/",
			want: []MessageParts{
				{Text: "https://example.com/path/",
					URL: "https://example.com/path/", HasURL: true},
			},
		},
		{
			name: "separated by comma",
			data: "https://a.com,https://b.com",
			want: []MessageParts{
				{Text: "https://a.com",
					URL: "https://a.com", HasURL: true},
				{Text: ","},
				{Text: "https://b.com",
					URL: "https://b.com", HasURL: true},
			},
		},
		{
			name: "two urls with no separator",
			data: "https://a.comhttps://b.com",
			want: []MessageParts{
				{Text: "https://a.com",
					URL: "https://a.com", HasURL: true},
				{Text: "https://b.com",
					URL: "https://b.com", HasURL: true},
			},
		},
		{
			name: "non-ascii text",
			data: "café https://example.com 日本語", //nolint:gosmopolitan
			want: []MessageParts{
				{Text: "café "},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: " 日本語"}, //nolint:gosmopolitan
			},
		},
		{
			name: "period lookalike",
			data: "https://example.com．",
			want: []MessageParts{
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: "．"},
			},
		},
		{
			name: "uppercase scheme",
			data: "HTTPS://example.com", // allowed per rfc3986
			want: []MessageParts{
				{Text: "HTTPS://example.com",
					URL: "HTTPS://example.com", HasURL: true},
			},
		},
		{
			name: "with port number",
			data: "https://example.com:8080/path",
			want: []MessageParts{
				{Text: "https://example.com:8080/path",
					URL: "https://example.com:8080/path", HasURL: true},
			},
		},
		{
			name: "ipv4 host",
			data: "http://192.168.1.1/admin",
			want: []MessageParts{
				{Text: "http://192.168.1.1/admin",
					URL: "http://192.168.1.1/admin", HasURL: true},
			},
		},
		{
			name: "ipv6 host",
			data: "http://[::1]:8080/",
			want: []MessageParts{
				{Text: "http://[::1]:8080/",
					URL: "http://[::1]:8080/", HasURL: true},
			},
		},
		{
			name: "markdown syntax",
			data: "[click here](https://example.com)",
			want: []MessageParts{
				{Text: "[click here]("},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: ")"},
			},
		},
		{
			name: "angle brackets",
			data: "<https://example.com>",
			want: []MessageParts{
				{Text: "<"},
				{Text: "https://example.com",
					URL: "https://example.com", HasURL: true},
				{Text: ">"},
			},
		},
		{
			name: "only cutset chars",
			data: "http://...,;:!?",
			want: []MessageParts{
				{Text: "http://...,;:!?"},
			},
		},
		{
			name: "empty string",
			data: "",
			want: []MessageParts{
				{Text: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{Data: tt.data}
			got := msg.GetParts()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s failed:\ngot  %+v\nwant %+v",
					tt.name, got, tt.want)
			}
		})
	}
}

package storage

import (
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
			want: []MessageParts{{Text: "hello, world!"}},
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
			data: "a https://one.example.com b http://two.example.com/path c",
			want: []MessageParts{
				{Text: "a "},
				{Text: "https://one.example.com",
					URL: "https://one.example.com", HasURL: true},
				{Text: " b "},
				{Text: "http://two.example.com/path",
					URL: "http://two.example.com/path", HasURL: true},
				{Text: " c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{Data: tt.data}
			got := msg.GetParts()

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("%s failed:\ngot  %+v\nwant %+v",
						tt.name, got[i], tt.want[i])
				}
			}
		})
	}
}

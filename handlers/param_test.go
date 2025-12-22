package handlers

import (
	"net/http"
	"net/url"
	"testing"
)

// TestGetRequestParameter validates the parameter value
// is read from the URL or form.
func TestGetRequestParameter(t *testing.T) {
	tests := []struct {
		name    string
		request *http.Request
		pathLen int
		field   string
		want    string
	}{
		{
			name: "URL parameter",
			request: &http.Request{
				URL: &url.URL{Path: "/download/id1"},
			},
			pathLen: 10,
			field:   "",
			want:    "id1",
		},
		{
			name: "Query parameter",
			request: &http.Request{
				URL: &url.URL{
					Path:     "/download/",
					RawQuery: "field=id2",
				},
			},
			pathLen: 10,
			field:   "field",
			want:    "id2",
		},
		{
			name: "Form parameter",
			request: &http.Request{
				Body: http.NoBody,
				Form: url.Values{"field": {"id3"}},
				URL:  &url.URL{Path: "/download/"},
			},
			pathLen: 10,
			field:   "field",
			want:    "id3",
		},
		{
			name: "No parameter",
			request: &http.Request{
				Body: http.NoBody,
				URL:  &url.URL{Path: "/download/"},
			},
			pathLen: 10,
			field:   "",
			want:    "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getRequestParameter(test.request, test.pathLen, test.field)
			if result != test.want {
				t.Fatalf("%s: expect '%q'; got '%q'", test.name, test.want, result)
			}
		})
	}
}

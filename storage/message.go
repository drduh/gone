package storage

import (
	"net/url"
	"regexp"
	"strings"
)

const cutSet = ".,;:!?)]}'\">"

var (
	schemeRe = regexp.MustCompile(`(?i)https?://`)
	urlRe    = regexp.MustCompile(
		`(?i)https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=%]+`)
)

func splitEmbeddedURLs(s string) []string {
	locs := schemeRe.FindAllStringIndex(s, -1)
	if len(locs) <= 1 {
		return []string{s}
	}

	var out []string

	start := 0
	for i := 1; i < len(locs); i++ {
		out = append(out, s[start:locs[i][0]])
		start = locs[i][0]
	}

	out = append(out, s[start:])

	return out
}

func trimCutset(data string, start, end int) int {
	for end > start &&
		strings.IndexByte(cutSet, data[end-1]) >= 0 {
		end--
	}
	return end
}

func isValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}

// GetParts extracts Message plain-text and URL parts.
func (m Message) GetParts() []MessageParts {
	matches := urlRe.FindAllStringIndex(m.Data, -1)
	if len(matches) == 0 {
		return []MessageParts{
			{Text: m.Data},
		}
	}

	parts := make([]MessageParts, 0, len(matches)*2+1)

	last := 0
	for _, match := range matches {
		start, end := match[0], match[1]

		pieces := splitEmbeddedURLs(m.Data[start:end])
		offset := start
		for _, piece := range pieces {
			pStart, pEnd := offset, offset+len(piece)
			offset = pEnd

			pEnd = trimCutset(m.Data, pStart, pEnd)
			candidate := m.Data[pStart:pEnd]

			if !isValidURL(candidate) {
				continue
			}

			if pStart > last {
				parts = append(parts, MessageParts{
					Text: m.Data[last:pStart],
				})
			}

			parts = append(parts, MessageParts{
				Text:   candidate,
				URL:    candidate,
				HasURL: true,
			})

			last = pEnd
		}
	}

	if last < len(m.Data) {
		parts = append(parts, MessageParts{
			Text: m.Data[last:],
		})
	}

	return parts
}

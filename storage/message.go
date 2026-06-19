package storage

import "regexp"

var urlRe = regexp.MustCompile(`https?://[^\s<]+`)

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

		if start > last {
			parts = append(parts, MessageParts{
				Text: m.Data[last:start],
			})
		}

		url := m.Data[start:end]

		parts = append(parts, MessageParts{
			Text:   url,
			URL:    url,
			HasURL: true,
		})

		last = end
	}

	if last < len(m.Data) {
		parts = append(parts, MessageParts{
			Text: m.Data[last:],
		})
	}

	return parts
}

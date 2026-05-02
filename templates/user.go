package templates

import "github.com/drduh/gone/storage"

// User represents user request state.
type User struct {

	// Request details
	storage.Owner `json:"user,omitempty"`

	// Whether the request originated from a browser
	IsBrowser bool `json:"isBrowser"`
}

package templates

import "github.com/drduh/gone/storage"

// User represents user request state.
type User struct {

	// Request details
	storage.Owner `json:"user,omitempty"`
}

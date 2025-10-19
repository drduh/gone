// Package handlers defines how the server handles configured routes
package handlers

// Request contains relevant HTTP request attributes to log.
type Request struct {

	// Handler path ("/")
	Action string `json:"action"`

	// IP address including port ("127.0.0.1:12345")
	Address string `json:"address"`

	// Masked address ("User123")
	Mask string `json:"mask"`

	// User agent ("Mozilla/5.0 ...")
	Agent string `json:"agent"`
}

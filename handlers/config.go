// Package handlers defines how the server handles configured routes
package handlers

// Request contains relevant HTTP request attributes to log.
type Request struct {

	// Handler path ("/")
	Action string `json:"action"`

	// User IP and port address ("127.0.0.1:12345")
	Address string `json:"address"`

	// User agent ("Mozilla/5.0 ...")
	Agent string `json:"agent"`
}

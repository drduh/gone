package server

import "crypto/tls"

// newTLSConfig returns the configuration
// and policy used by the HTTPS server.
func newTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}
}

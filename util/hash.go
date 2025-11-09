package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sum returns the SHA-256 hash sum.
func Sum(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// RandomNumber returns a zero-padded 3-digit string,
// like "007", "123", "999", or "000" on error.
func RandomNumber() string {
	num, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "000"
	}
	return fmt.Sprintf("%03d", num.Int64())
}

package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var names = []string{
	"Alice", "Bob", "Charlie", "Diana", "Eve",
	"Frank", "Grace", "Henry", "Ivan", "Judy",
	"Ken", "Laura", "Mallory", "Nancy", "Olivia",
	"Peggy", "Quentin", "Rupert", "Sam", "Trent",
	"Uma", "Victor", "Wendy", "Xavier", "Yvonne", "Zack",
}

// randomInt returns a random int64 up to max, or -1 on error.
func randomInt(max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return -1
	}
	return n.Int64()
}

// RandomNumber returns a zero-padded 3-digit string,
// like "007", "123", "999", or "000" on error.
func RandomNumber() string {
	n := randomInt(1000)
	if n < 0 {
		return "000"
	}
	return fmt.Sprintf("%03d", n)
}

// RandomName returns a random string from the names list,
// like "Alice", "Zack", or "Bob" on error.
func RandomName() string {
	i := randomInt(int64(len(names)))
	if i < 0 {
		return "Bob"
	}
	return names[i]
}

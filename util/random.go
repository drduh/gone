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

// RandomName returns a random string from the names list,
// like "Alice", "Zack", or "Bob" on error.
func RandomName() string {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(names))))
	if err != nil {
		return "Bob"
	}
	return names[index.Int64()]
}

// RandomNumber returns a zero-padded 3-digit string,
// like "007", "123", "999", or "000" on error.
func RandomNumber() string {
	num, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "000"
	}
	return fmt.Sprintf("%03d", num.Int64())
}

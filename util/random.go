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

var nato = []string{
	"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot",
	"Golf", "Hotel", "India", "Juliett", "Kilo", "Lima",
	"Mike", "November", "Oscar", "Papa", "Quebec", "Romeo",
	"Sierra", "Tango", "Uniform", "Victor", "Whiskey", "X-ray",
	"Yankee", "Zulu",
}

// randomInt returns a random int64 up to max, or -1 on error.
func randomInt(max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return -1
	}
	return n.Int64()
}

// pickRandom returns a random string from list, or fallback on error.
func pickRandom(list []string, fallback string) string {
	if len(list) == 0 {
		return fallback
	}
	i := randomInt(int64(len(list)))
	if i < 0 {
		return fallback
	}
	return list[i]
}

// FlipCoin returns "heads" or "tails" at random.
func FlipCoin() string {
	if f := randomInt(2); f == 0 {
		return "heads"
	}
	return "tails"
}

// RandomName returns a random string from the names list,
// like "Alice", "Zack", or "Bob" on error.
func RandomName() string {
	return pickRandom(names, "Bob")
}

// RandomNato returns a random string from the nato list,
// like "Alpha", "Zulu", or "Bravo" on error.
func RandomNato() string {
	return pickRandom(nato, "Bravo")
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

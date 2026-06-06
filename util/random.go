package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

var nato = []string{
	"Alfa", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot",
	"Golf", "Hotel", "India", "Juliett", "Kilo", "Lima",
	"Mike", "November", "Oscar", "Papa", "Quebec", "Romeo",
	"Sierra", "Tango", "Uniform", "Victor", "Whiskey", "X-ray",
	"Yankee", "Zulu",
}

// randomInt returns a random int64 up to max; or -1 on error.
func randomInt(maximum int64) int64 {
	if maximum <= 0 {
		return -1
	}
	n, err := rand.Int(rand.Reader, big.NewInt(maximum))
	if err != nil {
		return -1
	}
	return n.Int64()
}

// pickRandom returns a random string from list; or fallback on error.
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
	const coinSides = 2
	if randomInt(coinSides) == 0 {
		return "heads"
	}
	return "tails"
}

// RandomHex returns a random string with hexadecimal
// characters only; or "0" on error.
func RandomHex(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return strings.Repeat("0", length)
	}
	result := hex.EncodeToString(bytes)
	return result[:length]
}

// RandomName returns a random string from the names list,
// like "Alice", "Zack"; or "Bob" on error.
func RandomName() string {
	return pickRandom(loadedNames, "Bob")
}

// RandomNato returns a random string from the nato list,
// like "Alpha", "Zulu"; or "Bravo" on error.
func RandomNato() string {
	return pickRandom(nato, "Bravo")
}

// RandomNumber returns a zero-padded 3-digit string,
// like "007", "123", "999"; or "000" on error.
func RandomNumber() string {
	const maxNumber = 999
	n := randomInt(maxNumber)
	if n < 0 {
		return "000"
	}
	return fmt.Sprintf("%03d", n)
}

// RandomID returns a 32-byte URL-encoded random string;
// or "unknown" on error.
func RandomID() string {
	const randomTokenBytes = 32
	bytes := make([]byte, randomTokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "unknown"
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// RandomMask returns a name and number combination.
func RandomMask() string {
	return RandomName() + RandomNumber()
}

// Random returns a random string of a given length.
func Random(length int) string {
	const charset = `ABCDEFGHJKLMNPQRTVWXYZ` +
		`-_2346789` + `abcdefghijkmnpqrtvwxyz`
	bytes := make([]byte, length)
	for i := range bytes {
		n := randomInt(int64(len(charset)))
		if n < 0 {
			bytes[i] = 'a'
		} else {
			bytes[i] = charset[n]
		}
	}
	return string(bytes)
}

// GetRandom returns a random string by requested path.
func GetRandom(path string) string {
	const defaultLength = 20

	var response string

	switch path {
	case "coin":
		response = FlipCoin()
	case "hex":
		response = RandomHex(defaultLength)
	case "id":
		response = RandomID()
	case "mask":
		response = RandomMask()
	case "name":
		response = RandomName()
	case "nato":
		response = RandomNato()
	case "number":
		response = RandomNumber()
	case "pass":
		response = Random(defaultLength)
	default:
		response = RandomMask()
	}
	return response
}

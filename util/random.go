package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
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

// RandomHex returns a hexadecimal string of byte size.
func RandomHex(numBytes int) string {
	b := make([]byte, numBytes)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// RandomName returns a random string from the names list,
// like "ancientLilac", "zestyWillow"; or "blueBay" on error.
func RandomName() string {
	return pickRandom(loadedNames, "blueBay")
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

// RandomID returns a 32-byte URL-encoded random string.
func RandomID() string {
	const randomTokenBytes = 32
	b := make([]byte, randomTokenBytes)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// RandomMask returns a name and number combination.
func RandomMask() string {
	return RandomName() + RandomNumber()
}

// Random returns a random string of a given length.
func Random(length int) string {
	const charset = `ABCDEFGHJKLMNPQRTVWXYZ` +
		`-_2346789` + `abcdefghijkmnpqrtvwxyz`
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomInt(int64(len(charset)))]
	}
	return string(b)
}

// GetRandom returns a random string by requested path.
func GetRandom(path string) string {
	const defaultLength = 32

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

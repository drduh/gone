package util

import (
	"slices"
	"strings"
	"testing"
)

// TestRandomInt tests for a value within expected range.
func TestRandomInt(t *testing.T) {
	result := randomInt(9001)
	if result < 0 || result >= 9001 {
		t.Errorf("randomInt returned %d", result)
	}
}

// TestRandom tests for value length and uniqueness.
func TestRandom(t *testing.T) {
	lengths := []int{0, 1, 8, 16, 32, 64}
	for _, l := range lengths {
		pass := Random(l)
		if len(pass) != l {
			t.Errorf("Random(%d) returned %d, want %d",
				l, len(pass), l)
		}
	}
	pass1 := Random(20)
	pass2 := Random(20)
	if pass1 == pass2 {
		t.Error("expected different pass")
	}
}

// TestPickRandom tests for a value from list or fallback.
func TestPickRandom(t *testing.T) {
	list := []string{"foo", "bar", "zoo"}
	fallback := "bar"

	result := pickRandom(list, fallback)
	found := slices.Contains(list, result)
	if !found && result != fallback {
		t.Errorf("pickRandom returned %s", result)
	}
	result = pickRandom(nil, fallback)
	if result != fallback {
		t.Errorf("pickRandom did not return fallback (%s)", result)
	}
}

// TestFlipCoin tests for a "heads" or "tails".
func TestFlipCoin(t *testing.T) {
	result := FlipCoin()
	if result != "heads" && result != "tails" {
		t.Errorf("FlipCoin returned %s", result)
	}
}

// TestRandomName tests for a value from names list or "Bob".
func TestRandomName(t *testing.T) {
	result := RandomName()
	found := slices.Contains(loadedNames, result)
	if !found && result != "Bob" {
		t.Errorf("RandomName returned %s", result)
	}
}

// TestRandomNato tests for a value from nato list or "Bravo".
func TestRandomNato(t *testing.T) {
	result := RandomNato()
	found := slices.Contains(nato, result)
	if !found && result != "Bravo" {
		t.Errorf("RandomNato returned %s", result)
	}
}

// TestRandomNumber tests for a valid zero-padded 3-digit string.
func TestRandomNumber(t *testing.T) {
	result := RandomNumber()
	if len(result) != 3 {
		t.Errorf("RandomNumber returned %s", result)
	}
}

// TestRandomHex tests for a valid hex value.
func TestRandomHex(t *testing.T) {
	length := 32
	result := RandomHex(length)
	if len(result) != length*2 {
		t.Errorf("expected length %d, got %d",
			length*2, len(result))
	}

	for _, c := range result {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("unexpected '%c' in result", c)
		}
	}
}

// TestGetRandom tests each random type produces output.
func TestGetRandom(t *testing.T) {
	for _, p := range []string{"coin", "hex", "id", "mask",
		"name", "nato", "number", "pass", "", "unknown"} {
		if got := GetRandom(p); got == "" {
			t.Fatalf("GetRandom(%q) empty", p)
		}
	}

	if got := GetRandom("hex"); len(got) != 64 {
		t.Fatalf("hex len = %d; want 64", len(got))
	}

	if got := GetRandom("pass"); len(got) != 32 {
		t.Fatalf("pass len = %d; want 32", len(got))
	}
}

// TestRandomID tests id length and URL-safe alphabet.
func TestRandomID(t *testing.T) {
	id := RandomID()
	if len(id) != 43 {
		t.Fatalf("len = %d; want 43", len(id))
	}

	for _, c := range id {
		if !strings.ContainsRune(
			"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", c) {
			t.Fatalf("unexpected char %q in %q", c, id)
		}
	}
}

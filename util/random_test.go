package util

import "testing"

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
	found := false
	for _, v := range list {
		if result == v {
			found = true
			break
		}
	}
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
	found := false
	for _, v := range names {
		if result == v {
			found = true
			break
		}
	}
	if !found && result != "Bob" {
		t.Errorf("RandomName returned %s", result)
	}
}

// TestRandomNato tests for a value from nato list or "Bravo".
func TestRandomNato(t *testing.T) {
	result := RandomNato()
	found := false
	for _, v := range nato {
		if result == v {
			found = true
			break
		}
	}
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
	length := 16
	result := RandomHex(length)
	if len(result) != length {
		t.Errorf("expected length %d, got %d",
			length, len(result))
	}
	for _, c := range result {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("unexpected '%c' in result", c)
		}
	}
}

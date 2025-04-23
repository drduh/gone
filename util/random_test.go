package util

import "testing"

// TestRandomInt tests randomInt returns values within expected range.
func TestRandomInt(t *testing.T) {
	result := randomInt(9001)
	if result < 0 || result >= 9001 {
		t.Errorf("randomInt returned %d", result)
	}
}

// TestPickRandom tests pickRandom returns value from list or fallback.
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

	result = pickRandom([]string{}, fallback)
	if result != fallback {
		t.Errorf("pickRandom did not return fallback (%s)", result)
	}
}

// TestFlipCoin tests FlipCoin returns "heads" or "tails".
func TestFlipCoin(t *testing.T) {
	result := FlipCoin()
	if result != "heads" && result != "tails" {
		t.Errorf("FlipCoin returned %s", result)
	}
}

// TestRandomName tests RandomName returns value from names list or "Bob".
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

// TestRandomNato tests RandomNato returns value from nato list or "Bravo".
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

// TestRandomNumber tests RandomNumber returns valid zero-padded 3-digit string.
func TestRandomNumber(t *testing.T) {
	result := RandomNumber()
	if len(result) != 3 {
		t.Errorf("RandomNumber returned %s", result)
	}
}

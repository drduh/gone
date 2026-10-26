package version

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

// TestGet tests full and redacted build detail.
func TestGet(t *testing.T) {
	full := Get(Full)
	if full["appID"] != ID || full["buildGoVers"] != Go ||
		len(full) != 10 {
		t.Fatalf("full = %v", full)
	}

	redacted := Get(Redacted)
	if len(redacted) != 1 || redacted["appID"] != ID {
		t.Fatalf("redacted = %v", redacted)
	}
}

// TestPrint tests JSON version output.
func TestPrint(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	Print()
	os.Stdout = orig
	_ = w.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("decode %q: %v", out, err)
	}
	if got["appID"] != ID {
		t.Fatalf("appID = %q; want %q", got["appID"], ID)
	}
}

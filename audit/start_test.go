package audit

import (
	"os"
	"path/filepath"
	"testing"
)

// Temporary test directory cleanup.
func cleanupTempDir(t *testing.T, tempDir string) {
	t.Helper()
	if err := os.RemoveAll(tempDir); err != nil {
		t.Fatalf("failed to remove temp dir: %v", err)
	}
}

// TestStartValidFilenames tests Auditor starts
// with a valid filename.
func TestStartValidFilenames(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-auditor-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer cleanupTempDir(t, tempDir)

	cases := []struct {
		name     string
		filename string
		debug    bool
	}{
		{"file", "test.log", false},
		{"file with debug", "debug.log", true},
		{"question mark", "foo?bar.log", false},
		{"double quote", `foo"bar.log`, false},
		{"single quote", "foo'bar.log", false},
		{"backslash", `foo\bar.log`, false},
		{"colon", "foo:bar.log", false},
		{"space", "foo bar.log", false},
		{"star", "foo*bar.log", false},
		{"tab", "foo\tbar.log", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fname := filepath.Join(tempDir, tc.filename)
			cfg := &Config{
				Filename: fname,
				Debug:    tc.debug,
			}
			auditor, err := Start(cfg)
			if err != nil {
				t.Errorf("got error: %v", err)
			}
			if auditor == nil {
				t.Errorf("expected auditor, got nil")
			}
			if _, err := os.Stat(fname); err != nil {
				t.Errorf("log not created: %v", err)
			}
		})
	}
}

// TestStartInvalidFilenames tests Auditor produces error
// when an invalid filename is used.
func TestStartInvalidFilenames(t *testing.T) {
	cases := []struct {
		name     string
		filename string
	}{
		{"directory path", "."},
		{"forbidden path", "/test.log"},
		{"invalid nul byte", string(
			[]byte{'f', 0, 'l', 'e'})},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				Filename: tc.filename,
				Debug:    false,
			}
			auditor, err := Start(cfg)
			if err == nil {
				t.Errorf("%q: expected error, got nil",
					tc.filename)
			}
			if auditor != nil {
				t.Errorf("%q: expected nil auditor",
					tc.filename)
			}
		})
	}
}

// TestStartDebug tests Auditor starts with debug flag.
func TestStartDebug(t *testing.T) {
	cases := []struct {
		name  string
		debug bool
	}{
		{"debug on", true},
		{"debug off", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				Filename: "",
				Debug:    tc.debug,
			}
			auditor, err := Start(cfg)
			if err != nil {
				t.Errorf("%v: unexpected error: %v",
					tc.debug, err)
			}
			if auditor == nil {
				t.Errorf("%v: expected auditor",
					tc.debug)
			}
		})
	}
}

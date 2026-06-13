package util

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func setupNames(t *testing.T, dir, filename, content string) {
	t.Helper()

	path := filepath.Join(dir, filename)

	if content == "" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			t.Fatalf("setupNames remove %s: %v", path, err)
		}

		return
	}

	err := os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("setupNames write %s: %v", path, err)
	}
}

// TestLoadNames tests loading names from various files.
func TestLoadNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		content  string
		filename string
		result   []string
	}{
		{
			name:    "empty file",
			content: "",
			result:  defaultNames(),
		},
		{
			name:    "newlines only",
			content: "\n\n",
			result:  defaultNames(),
		},
		{
			name:    "basic",
			content: "foo\nbar123\nZoo\n",
			result:  []string{"foo", "bar123", "Zoo"},
		},
		{
			name:    "mix empty lines",
			content: "foo\n\nbar\n\n",
			result:  []string{"foo", "bar"},
		},
		{
			name:    "spaces trimmed",
			content: "  foo  \n bar\n  zoo ",
			result:  []string{"foo", "bar", "zoo"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			filename := "names.txt"
			setupNames(t, dir, filename, tc.content)

			got := loadNames(dir, filename)
			if !slices.Equal(got, tc.result) {
				t.Errorf("%s: got %v, want %v",
					tc.name, got, tc.result)
			}
		})
	}
}

// TestLoadNamesMissingFile tests loading default names.
func TestLoadNamesMissingFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	got := loadNames(dir, "does-not-exist.txt")
	if !slices.Equal(got, defaultNames()) {
		t.Errorf("expected default names, got %v", got)
	}
}

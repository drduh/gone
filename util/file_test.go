package util

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var namesFile = "names.txt"

func setupNames(filename, content string) {
	if content == "" {
		_ = os.Remove(filename)
	} else {
		f, _ := os.Create(filename)
		if f != nil {
			_, _ = f.WriteString(content)
			_ = f.Close()
		}
	}
}

// TestFormatSize tests size integer conversion to readable string.
func TestFormatSize(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0 bytes"},
		{200, "200 bytes"},
		{1024, "1 kb"},
		{5000, "4.88 kb"},
		{1048576, "1 mb"},
		{5242880, "5 mb"},
		{100000000, "95.37 mb"},
		{50000 * 50000, "2.33 gb"},
		{50000 * 50000 * 50000, "113.69 tb"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("input: %d", tt.input), func(t *testing.T) {
			if got := FormatSize(tt.input); got != tt.expect {
				t.Errorf("%d = %v; expect %v", tt.input, got, tt.expect)
			}
		})
	}
}

// TestGetOutput tests output destinations and errors.
func TestGetOutput(t *testing.T) {
	writer, err := GetOutput("")
	if err != nil {
		t.Fatalf("empty filename error: %v", err)
	}
	if writer != os.Stdout {
		t.Errorf("expected stdout writer, got %v", writer)
	}

	testFile := "test.txt"
	writer, err = GetOutput(testFile)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	file, ok := writer.(*os.File)
	if !ok {
		t.Errorf("expected file writer, got %v", writer)
	} else {
		if err := file.Close(); err != nil {
			t.Errorf("error closing file: %v", err)
		}
		if err := os.Remove(testFile); err != nil {
			t.Errorf("error removing file: %v", err)
		}
	}

	badFile := string([]byte{0})
	writer, err = GetOutput(badFile)
	if err == nil {
		t.Errorf("did not receive error")
	}
	if writer != nil {
		t.Errorf("did not expect writer")
	}
}

// TestGetOutputAppends tests existing output file is appended.
func TestGetOutputAppends(t *testing.T) {
	testFile := "test_append.txt"
	defer func() { _ = os.Remove(testFile) }()

	first, err := GetOutput(testFile)
	if err != nil {
		t.Fatalf("first open error: %v", err)
	}
	if _, err := fmt.Fprint(first, "hello"); err != nil {
		t.Fatalf("first write error: %v", err)
	}
	if err := first.(*os.File).Close(); err != nil {
		t.Fatalf("first close error: %v", err)
	}

	second, err := GetOutput(testFile)
	if err != nil {
		t.Fatalf("second open error: %v", err)
	}
	if _, err := fmt.Fprint(second, " world"); err != nil {
		t.Fatalf("second write error: %v", err)
	}
	if err := second.(*os.File).Close(); err != nil {
		t.Fatalf("second close error: %v", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if got := string(data); got != "hello world" {
		t.Errorf("expected content 'hello world', got %q", got)
	}
}

// TestLoadNamesEmpty tests a file containing only blank lines.
func TestLoadNamesEmpty(t *testing.T) {
	setupNames(namesFile, "\n\n")
	defer func() { _ = os.Remove(namesFile) }()
	names := loadNames(namesFile)
	if len(names) != len(defaultNames) {
		t.Errorf("expected %d default names, got %d",
			len(defaultNames), len(names))
	}
}

// TestLoadNamesExist test a file with three names is loaded.
func TestLoadNamesExist(t *testing.T) {
	setupNames(namesFile, "Sun\nEarth\nMoon\n")
	defer func() { _ = os.Remove(namesFile) }()
	names := loadNames(namesFile)
	expected := []string{"Sun", "Earth", "Moon"}
	if len(names) != len(expected) {
		t.Errorf("expected %d names, got %d", len(expected), len(names))
	}
	for i := range expected {
		if names[i] != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], names[i])
		}
	}
}

// TestLoadNamesMissing tests setting up names with no file.
func TestLoadNamesMissing(t *testing.T) {
	setupNames(namesFile, "")
	names := loadNames(namesFile)
	if len(names) != len(defaultNames) {
		t.Errorf("expected %d names, got %d", len(defaultNames), len(names))
	}
}

// TestLoadNamesMix tests a file with blank lines and valid names.
func TestLoadNamesMix(t *testing.T) {
	setupNames(namesFile, "\nSun\n\nMoon\n\n")
	defer func() { _ = os.Remove(namesFile) }()
	names := loadNames(namesFile)
	expected := []string{"Sun", "Moon"}
	if len(names) != len(expected) {
		t.Errorf("expected %d names, got %d", len(expected), len(names))
	}
	for i := range expected {
		if names[i] != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], names[i])
		}
	}
}

// TestLoadNamesSpaces tests file with extra spaces to trim.
func TestLoadNamesSpaces(t *testing.T) {
	setupNames(namesFile, "  Sun  \n\tEarth\t\n  Moon\n")
	defer func() { _ = os.Remove(namesFile) }()
	names := loadNames(namesFile)
	expected := []string{"Sun", "Earth", "Moon"}
	if len(names) != len(expected) {
		t.Errorf("expected %d names, got %d", len(expected), len(names))
	}
	for i := range expected {
		if names[i] != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], names[i])
		}
		if names[i] != strings.TrimSpace(names[i]) {
			t.Errorf("index %d: name %q was not trimmed", i, names[i])
		}
	}
}

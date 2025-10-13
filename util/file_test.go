package util

import (
	"fmt"
	"os"
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
		{0, "0 Bytes"},
		{200, "200.00 Bytes"},
		{1024, "1.00 KB"},
		{5000, "4.88 KB"},
		{1048576, "1.00 MB"},
		{5242880, "5.00 MB"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("input: %d", tt.input), func(t *testing.T) {
			if got := FormatSize(tt.input); got != tt.expect {
				t.Errorf("%d=%v; expect %v", tt.input, got, tt.expect)
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

func TestLoadNamesEmpty(t *testing.T) {
	setupNames(namesFile, "\n\n")
	defer func() { _ = os.Remove(namesFile) }()
	names := loadNames(namesFile)
	if len(names) != len(defaultNames) {
		t.Errorf("expected names, got %d empty lines", len(names))
	}
}

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

func TestLoadNamesMissing(t *testing.T) {
	setupNames(namesFile, "")
	names := loadNames(namesFile)
	if len(names) != len(defaultNames) {
		t.Errorf("expected %d names, got %d", len(defaultNames), len(names))
	}
}

package util

import (
	"fmt"
	"os"
	"testing"
)

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
	f1, ok := first.(*os.File)
	if !ok {
		t.Fatalf("expected *os.File, got %T", first)
	}
	if err := f1.Close(); err != nil {
		t.Fatalf("first close error: %v", err)
	}

	second, err := GetOutput(testFile)
	if err != nil {
		t.Fatalf("second open error: %v", err)
	}
	if _, err := fmt.Fprint(second, " world"); err != nil {
		t.Fatalf("second write error: %v", err)
	}
	f2, ok := second.(*os.File)
	if !ok {
		t.Fatalf("expected *os.File, got %T", second)
	}
	if err := f2.Close(); err != nil {
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

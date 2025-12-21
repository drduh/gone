package storage

import "testing"

const extraChars = "_.-"

// TestRemoveInvalidChars validates invalid characters
// are removed from filename strings.
func TestRemoveInvalidChars(t *testing.T) {
	tests := []struct {
		input string
		allow string
		want  string
	}{
		{"filename.txt", extraChars, "filename.txt"},
		{"file@name!.txt", extraChars, "filename.txt"},
		{"_file-name.txt", extraChars, "_file-name.txt"},
		{"123-abc_ABC.txt", extraChars, "123-abc_ABC.txt"},
		{"chars@#$name.txt", extraChars, "charsname.txt"},
		{"afile", extraChars, "afile"},
		{".....", extraChars, "....."},
		{"!@#$%^&()", extraChars, ""},
		{"", extraChars, ""},
	}
	for _, test := range tests {
		result := removeInvalidChars(test.input, test.allow)
		if result != test.want {
			t.Fatalf("name: '%s' ('%s' allowed): '%s'; want: '%s'",
				test.input, test.allow, result, test.want)
		}
	}
}

// TestTruncateName validates filenames are truncated
// to the desired length.
func TestTruncateName(t *testing.T) {
	tests := []struct {
		base   string
		ext    string
		length int
		want   string
	}{
		{"shrt", ".exceedinglylongext", 10, "shrt.exce"},
		{"base", ".longextension", 12, "base.long"},
		{"exactfit", ".jpeg", 13, "exactfit.jpeg"},
		{"longfilename", ".txt", 10, "longfi.txt"},
		{"truncatebase", ".png", 8, "trun.png"},
		{"onlybase", ".toolong", 8, "onl.tool"},
		{"short", ".dat", 9, "short.dat"},
		{"example", "", 5, "examp"},
		{"", ".zip", 5, ".zip"},
		{"", "", 0, ""},
	}
	for _, test := range tests {
		result := truncateName(test.base, test.ext, test.length)
		if result != test.want {
			t.Fatalf("base: '%s', ext: '%s' (%d); got: '%s', want: '%s'",
				test.base, test.ext, test.length, result, test.want)
		}
	}
}

// TestSanitizeName validates filenames do not exceed
// length nor contain disallowed special characters.
func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input      string
		maxLength  int
		extraChars string
		want       string
	}{
		{"my@invalid#name?.txt", 20, extraChars, "myinvalidname.txt"},
		{"averylongfilename.png", 15, extraChars, "averylongfi.png"},
		{"my.file.name.txt", 20, extraChars, "my.file.name.txt"},
		{".hiddenfile", 15, extraChars, defaultName + ".hidd"},
		{"myfilename.png", 20, extraChars, "myfilename.png"},
		{"/path/to/file.txt", 20, extraChars, "file.txt"},
		{"myfilenames", 10, extraChars, "myfilename"},
		{".@#$%^&*.png", 15, extraChars, "..png"},
		{"@#$%^&*", 10, extraChars, defaultName},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := SanitizeName(test.input, test.maxLength, test.extraChars)
			if result != test.want {
				t.Fatalf("length: %d, special: '%s', got: '%s', want: '%s'",
					test.maxLength, test.extraChars, result, test.want)
			}
		})
	}
}

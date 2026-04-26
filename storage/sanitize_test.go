package storage

import (
	"strings"
	"testing"
)

const extraChars = "_.- "

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
		{"file", ". . .   long", 12, "file...lo"},
		{"truncatebase", ".png", 8, "trun.png"},
		{"onlybase", ".toolong", 8, "onl.tool"},
		{"test", "test test", 18, "testtestt"},
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
		extraChars string
		maxLength  int
		want       string
	}{
		{"my@invalid#name?.txt", extraChars, 20, "myinvalidname.txt"},
		{"averylongfilename.png", extraChars, 15, "averylongfi.png"},
		{"my.file.name.txt", extraChars, 20, "my.file.name.txt"},
		{".hiddenfile", extraChars, 15, defaultName + ".hidd"},
		{"myfilename.png", extraChars, 20, "myfilename.png"},
		{"/path/to/file.txt", extraChars, 20, "file.txt"},
		{"myfilenames", extraChars, 10, "myfilename"},
		{"@#$%^&*.png", extraChars, 20, defaultName},
		{"!@#$%^&*()[]{}<>", extraChars, 20, defaultName},
		{"!@# a copy.alongext", extraChars, 20, "a copy.alon"},
		{"<>$ a copy.m", extraChars, 10, "a copy.m"},
		{"a copy.my copy", extraChars, 15, "a copy.myco"},
		{"filename.", extraChars, 15, "filename."},
		{"/etc/passwd", extraChars, 15, "passwd"},
		{"name\u0000.txt", extraChars, 20, "name.txt"},
		{"control\ttest.txt", extraChars, 20, "controltest.txt"},
		{"/path/../file.txt", extraChars, 20, "file.txt"},
		{"#⃣#⃣#⃣™ℹ↔", extraChars, 20, defaultName},
		{"#⃣#⃣#⃣™ℹ↔.txt", extraChars, 20, defaultName + ".txt"},
		{" #⃣ #⃣ #⃣ ™ ℹ ↔ ↕ ↖ ↗ ↘ ↙ ↩ ↪ ⌚ ⌛ ▪ ▫ ▶ ◀ ◻ ◼ ◽ ◾ ♠ ♣ ♥ ♦ ♨ ⤴ ⤵ ⬅ ⬆ ⬇ ⬛ ⬜ ⭕ 〰 〽 ㊗ ㊙ ⏏ 🟰 ‼ ⁉ 〰 ⭕ 〽 © ® ™ *⃣ ℹ Ⓜ ㊗ ㊙ 󠁧󠁢󠁥 ♨ ♟ ⌨", extraChars, 20, defaultName},
		{"<script>alert('xss')</script>", extraChars, 25, "script"},
		{" ", extraChars, 15, defaultName},
		{"percent%encoded%name.doc", extraChars, 20, defaultName},
		{"filename%20with%20spaces.txt", extraChars, 20,
			"filename with sp.txt"},
		{"filename with spaces.txt", extraChars, 25,
			"filename with spaces.txt"},
		{"my%2Fcool%2Bdoc%26about%2Cstuff.md", extraChars + "/+&,", 40,
			"cool+doc&about,stuff.md"},
		{"example." + strings.Repeat("x", 1000), extraChars, 20,
			"example.xxxx"},
		{strings.Repeat("a", 1000) + ".txt", extraChars, 50,
			strings.Repeat("a", 46) + ".txt"},
		{strings.Repeat(".", 200), extraChars, 80,
			strings.Repeat(".", 80)},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := SanitizeName(test.input, test.extraChars, test.maxLength)
			if result != test.want {
				t.Fatalf("length: %d, special: '%s', got: '%s', want: '%s'",
					test.maxLength, test.extraChars, result, test.want)
			}
		})
	}
}

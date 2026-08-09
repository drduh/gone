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
		{".golangci.yml", extraChars, ".golangci.yml"},
		{`foo"bar.txt`, extraChars + `"`, `foo"bar.txt`},
		{`foo"bar.txt`, extraChars, "foobar.txt"},
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
		{"test", "test ext", 18, "testteste"},
		{"short", ".dat", 9, "short.dat"},
		{"example", "", 5, "examp"},
		{"", ".zip", 5, ".zip"},
		{"", "", 0, defaultName},
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
// length nor contain unapproved special characters.
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
		{"/etc/passwd", extraChars, 20, "passwd"},
		{"./etc/passwd", extraChars, 20, "passwd"},
		{"．．/passwd", extraChars, 20, "passwd"},
		{"../../../etc/passwd", extraChars, 20, "passwd"},
		{`..\..\..\etc\passwd`, extraChars, 20, "passwd"},
		{"....//....//etc/passwd", extraChars, 20, "passwd"},
		{"..%c0%af..%c0%afetc/passwd", extraChars, 20, "passwd"},
		{"..%2F..%2Fetc%2Fpasswd", extraChars, 20, "passwd"},
		{"safe.txt/../../etc/passwd", extraChars, 20, "passwd"},
		{"file.txt\x00../../etc/passwd", extraChars, 20, "passwd"},
		{"..\\..\\..\\etc\\passwd", extraChars, 20, "passwd"},
		{"..%2F..%2F..%2Fetc%2Fpasswd", extraChars, 20, "passwd"},
		{"..%252F..%252F..%252Fetc%252Fpasswd", extraChars, 20,
			"..2F..2F..2Fet"},
		{"payload.ph%p.txt", extraChars, 20, defaultName},
		{"payload.ph%61.txt", extraChars, 20, "payload.pha.txt"},
		{"...........", extraChars, 20, "..........."},
		{"test\x00name.txt", extraChars, 20, "testname.txt"},
		{"name\u0000.txt", extraChars, 20, "name.txt"},
		{"\u003cscript\u003e", extraChars, 20, "script"},
		{"control\ttest.txt", extraChars, 20, "controltest.txt"},
		{"/path/../file.txt", extraChars, 20, "file.txt"},
		{"$(rm -rf /)", extraChars, 20, defaultName},
		{"#⃣#⃣#⃣™ℹ↔", extraChars, 20, defaultName},
		{"#⃣#⃣#⃣™ℹ↔.txt", extraChars, 20, defaultName + ".txt"},
		{` #⃣ #⃣ #⃣ ™ ℹ ↔ ↕ ↖ ↗ ↘ ↙ ↩ ↪ ⌚ ⌛ ▪ ▫ ▶ ◀ ◻ ◼ ◽ ◾ ♠ ♣ ♥ ♦ ♨ ` +
			`⤴ ⤵ ⬅ ⬆ ⬇ ⬛ ⬜ ⭕ 〰 〽 ㊗ ㊙ ⏏ 🟰 ‼ ⁉ 〰 ⭕ 〽 © ® ™ *⃣ ℹ` +
			`Ⓜ ㊗ ㊙ 󠁧󠁢󠁥 ♨ ♟ ⌨`, extraChars, 20, defaultName},
		{"j%61vascript:alert(1)", extraChars, 25, "javascriptalert1"},
		{"<script>alert('xss')</script>", extraChars, 25, "script"},
		{"<img src=x onerror=alert(1)>", extraChars + "(" + ")", 25,
			"img srcx onerroralert(1)"},
		{"'; DROP TABLE all;--", extraChars, 20, "DROP TABLE all--"},
		{" ", extraChars, 15, defaultName},
		{"   .txt", extraChars, 15, defaultName + ".txt"},
		{"%20%20%20", extraChars, 15, defaultName},
		{"percent%encoded%name.doc", extraChars, 20, defaultName},
		{"filename%20with%20spaces.txt", extraChars, 20,
			"filename with sp.txt"},
		{"filename with spaces.txt", extraChars, 25,
			"filename with spaces.txt"},
		{"my%2Fcool%2Bdoc%26about%2Cstuff.md", extraChars + "/+&,", 40,
			"cool+doc&about,stuff.md"},
		{"unicøde.txt", extraChars, 20, "unicde.txt"},
		{`""foo.txt`, extraChars, 20, "foo.txt"},
		{"example." + strings.Repeat("x", 1000), extraChars, 20,
			"example.xxxx"},
		{strings.Repeat("a", 1000) + ".txt", extraChars, 50,
			strings.Repeat("a", 46) + ".txt"},
		{strings.Repeat(".", 200), extraChars, 80,
			strings.Repeat(".", 80)},
		{"DATA:IMAGE/PNG;BASE64,AAAA",
			extraChars, 40, defaultName},
		{"data%3Aimage%2Fpng%3Bbase64%2CAAAA",
			extraChars, 40, defaultName},
		{"data:image/svg+xml,%3Csvg%20onload%3Dalert(1)%3E",
			extraChars, 40, defaultName},
		{`![normal alt text](https://example.invalid/image.png)`,
			extraChars, 40, "image.png"},
		{`<img alt="profile" src="data:image/png;base64,AAAA">`,
			extraChars, 40, defaultName},
		{`![preview](DATA:text/html;base64,PHNjcmlwdD4=)`,
			extraChars, 40, defaultName},
		{`![encoded](data%3Aimage%2Fpng%3Bbase64%2CAAAA)`,
			extraChars, 40, defaultName},
		{`C:\Windows\System32\drivers\etc\hosts`,
			extraChars, 20, "hosts"},
		{"report+final.txt", extraChars, 30, "report final.txt"},
		{"report%2Bfinal.txt", extraChars, 30, "reportfinal.txt"},
		{"report%2Bfinal.txt", extraChars + "+", 30, "report+final.txt"},
		{"report%252Bfinal.txt", extraChars, 30, "report2Bfinal.txt"},
		{"100%25-real.txt", extraChars, 30, "100-real.txt"},
		{"100%real.txt", extraChars, 30, defaultName},
		{`\\server\share\secret.txt`, extraChars, 20, "secret.txt"},
		{"..%5C..%5Cetc%5Cpasswd", extraChars, 20, "passwd"},
		{"..%255C..%255Cetc%255Cpasswd", extraChars, 20, "..5C..5Cet"},
		{"file.txt\r\nSet-Cookie: sess=evil", extraChars, 40, "file.txtS"},
		{"file.txt\r\nX-Injected: true", extraChars, 40, "file.txtX"},
		{"file.txt\nLocation: https://evil.com", extraChars, 50, "evil.com"},
		{"café-menu_2026.pdf", extraChars, 50, "caf-menu_2026.pdf"},
		{"cafe\u0301.pdf", extraChars, 50, "cafe.pdf"},
		{"file\xff\xfe.txt", extraChars, 50, "file.txt"},
		{"sales_data:2024.csv", extraChars, 40, defaultName},
		{"\u202Erun.exe", extraChars, 40, "run.exe"},
		{"da%74a:xxx", extraChars, 40, defaultName},
		{"DaTa:xxx", extraChars, 40, defaultName},
		{"dir∕subdir∕file.txt", extraChars, 20, "dirsubdirfile.txt"},
		{"dir／subdir／file.txt", extraChars, 20, "dirsubdirfile.txt"},
		{"archive.tar.gz", extraChars, 30, "archive.tar.gz"},
		{"UPPER.PNG", extraChars, 30, "UPPER.PNG"},
		{"file%2Etxt", extraChars, 30, "file.txt"},
		{"/", extraChars, 20, defaultName},
		{`\\`, extraChars, 20, defaultName},
		{"", extraChars, 20, defaultName},
		{"   ", extraChars, 20, defaultName},
		{"%20%20%20", extraChars, 20, defaultName},
		{".", extraChars, 20, defaultName},
		{"..", extraChars, 20, defaultName},
		{".txt", extraChars, 20, defaultName + ".txt"},
		{"file.", extraChars, 30, "file."},
		{"file.txt   ", extraChars, 30, "file.txt"},
		{"file..txt", extraChars, 30, "file..txt"},
		{" file . txt ", extraChars, 30, "file.txt"},
		{"file%0D%0Ainjected.txt", extraChars, 30, "fileinjected.txt"},
		{"file%09tab.txt", extraChars, 30, "filetab.txt"},
		{"file%FF.txt", extraChars, 30, "file.txt"},
		{"file.txt", extraChars, -1, defaultName},
		{"file.txt", extraChars, 0, defaultName},
		{"file.txt", extraChars, 1, "."},
		{"file.txt", extraChars, 4, ".txt"},
		{"file.txt", extraChars, 5, "f.txt"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := SanitizeName(
				test.input, test.extraChars, test.maxLength)
			if result != test.want {
				t.Fatalf("'%s' (length %d, special '%s')\ngot  '%s'\nwant '%s'",
					test.input, test.maxLength, test.extraChars, result, test.want)
			}
		})
	}
}

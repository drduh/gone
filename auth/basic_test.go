package auth

import "testing"

// TestBasic test basic token auth.
func TestBasic(t *testing.T) {
	secret := []byte("correct-token")

	cases := []struct {
		name  string
		token []byte
		want  bool
	}{
		{"correct token",
			[]byte("correct-token"), true},
		{"incorrect token",
			[]byte("incorrect-token"), false},
		{"correct token with trailing byte",
			[]byte("correct-tokenz"), false},
		{"partial correct token",
			[]byte("correct"), false},
		{"empty token",
			[]byte{}, false},
		{"nil token",
			nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Basic(secret, tc.token)
			if got != tc.want {
				t.Errorf("Basic(%q, %q) should be %v",
					secret, tc.token, tc.want)
			}
		})
	}
}

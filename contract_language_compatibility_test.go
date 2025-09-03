package synnergy

import "testing"

func TestIsLanguageSupported(t *testing.T) {
	cases := []struct {
		lang string
		want bool
	}{
		{"wasm", true},
		{"golang", true},
		{"javascript", true},
		{"solidity", true},
		{"rust", true},
		{"python", true},
		{"yul", true},
		{"haskell", false},
	}
	for _, c := range cases {
		if got := IsLanguageSupported(c.lang); got != c.want {
			t.Errorf("%s: expected %v got %v", c.lang, c.want, got)
		}
	}
}

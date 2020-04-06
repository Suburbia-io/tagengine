package tagengine

import (
	"log"
	"testing"
)

func TestNGramLength(t *testing.T) {
	type Case struct {
		Input  string
		Length int
	}

	cases := []Case{
		{"a b c", 3},
		{"  xyz\nlkj  dflaj a", 4},
		{"a", 1},
		{" a", 1},
		{"a", 1},
		{" a\n", 1},
		{" a ", 1},
		{"\tx\ny\nz q ", 4},
	}

	for _, tc := range cases {
		length := ngramLength(tc.Input)
		if length != tc.Length {
			log.Fatalf("%s: %d != %d", tc.Input, length, tc.Length)
		}
	}
}

package tagengine

import "testing"

func TestSanitize(t *testing.T) {
	sanitize := newSanitizer()

	type Case struct {
		In  string
		Out string
	}

	cases := []Case{
		{"", ""},
		{"123abc", "123 abc"},
		{"abc123", "abc 123"},
		{"abc123xyz", "abc 123 xyz"},
		{"1f2", "1 f 2"},
		{" abc", "abc"},
		{" ; KitKat/m&m's (bÖttle) @ ", "; kitkat / m & ms ( bottle ) @"},
		{" Pott`s gin   königs beer;SOJU  ", "potts gin konigs beer ; soju"},
		{"brot & brötchen", "brot & brotchen"},
		{"Gâteau au fromage blanc, Stück", "gateau au fromage blanc , stuck"},
		{"Maisels Weisse Weißbier 0,5l", "maisels weisse weissbier 0 , 5 l"},
		{"Maisels´s Weisse - Hefeweizen 0,5l", "maiselss weisse - hefeweizen 0 , 5 l"},
		{"€", "€"},
	}

	for _, tc := range cases {
		out := sanitize(tc.In)
		if out != tc.Out {
			t.Fatalf("%v != %v", out, tc.Out)
		}
	}
}

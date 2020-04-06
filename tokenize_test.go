package tagengine

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	type Case struct {
		Input    string
		MaxNgram int
		Output   []string
	}

	cases := []Case{
		{
			Input:    "a bb c d",
			MaxNgram: 3,
			Output: []string{
				"a", "c", "d", "bb",
				"c d", "a bb", "bb c",
				"a bb c", "bb c d",
			},
		}, {
			Input:    "a b",
			MaxNgram: 3,
			Output: []string{
				"a", "b", "a b",
			},
		}, {
			Input:    "- b c d",
			MaxNgram: 3,
			Output: []string{
				"b", "c", "d",
				"- b", "b c", "c d",
				"- b c", "b c d",
			},
		}, {
			Input:    "a a b c d c d",
			MaxNgram: 3,
			Output: []string{
				"a", "b", "c", "d",
				"a a", "a b", "b c", "c d", "d c",
				"a a b", "a b c", "b c d", "c d c", "d c d",
			},
		},
	}

	for _, tc := range cases {
		output := Tokenize(tc.Input, tc.MaxNgram)
		if !reflect.DeepEqual(output, tc.Output) {
			t.Fatalf("%s: %#v", tc.Input, output)
		}
	}
}

package tagengine

import (
	"reflect"
	"testing"
)

func TestRulesSet(t *testing.T) {
	rs := NewRuleSet()
	rs.AddRule(Rule{
		Tag:      "cc/2",
		Includes: []string{"cola", "coca"},
	})
	rs.AddRule(Rule{
		Tag:      "cc/0",
		Includes: []string{"coca cola"},
	})
	rs.AddRule(Rule{
		Tag:      "cz/2",
		Includes: []string{"coca", "zero"},
	})
	rs.AddRule(Rule{
		Tag:      "cc0/3",
		Includes: []string{"zero", "coca", "cola"},
	})
	rs.AddRule(Rule{
		Tag:      "cc0/3.1",
		Includes: []string{"coca", "cola", "zero"},
		Excludes: []string{"pepsi"},
	})
	rs.AddRule(Rule{
		Tag:      "spa",
		Includes: []string{"spa"},
		Blocks:   []string{"cc/0", "cc0/3", "cc0/3.1"},
	})

	type TestCase struct {
		Input   string
		Matches []Match
	}

	cases := []TestCase{
		{
			Input: "coca-cola zero",
			Matches: []Match{
				{"cc0/3.1", 0.3},
				{"cc0/3", 0.3},
				{"cz/2", 0.2},
				{"cc/2", 0.2},
			},
		}, {
			Input: "coca cola",
			Matches: []Match{
				{"cc/0", 0.6},
				{"cc/2", 0.4},
			},
		}, {
			Input: "coca cola zero pepsi",
			Matches: []Match{
				{"cc0/3", 0.3},
				{"cc/0", 0.3},
				{"cz/2", 0.2},
				{"cc/2", 0.2},
			},
		}, {
			Input:   "fanta orange",
			Matches: []Match{},
		}, {
			Input: "coca-cola zero / fanta / spa",
			Matches: []Match{
				{"cz/2", 0.4},
				{"cc/2", 0.4},
				{"spa", 0.2},
			},
		},
	}

	for _, tc := range cases {
		matches := rs.Match(tc.Input)
		if !reflect.DeepEqual(matches, tc.Matches) {
			t.Fatalf("%v != %v", matches, tc.Matches)
		}
	}
}

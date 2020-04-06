package tagengine

import (
	"sort"
)

type RuleSet struct {
	root     *node
	maxNgram int
	sanitize func(...string) string
	rules    []*Rule
}

func NewRuleSet() *RuleSet {
	return &RuleSet{
		root: &node{
			Token:    "/",
			Children: map[string]*node{},
		},
		sanitize: newSanitizer(),
		rules:    []*Rule{},
	}
}

func NewRuleSetFromList(rules []Rule) *RuleSet {
	rs := NewRuleSet()
	rs.AddRule(rules...)
	return rs
}

func (t *RuleSet) Add(ruleOrGroup ...interface{}) {
	for _, ix := range ruleOrGroup {
		switch x := ix.(type) {
		case Rule:
			t.AddRule(x)
		case RuleGroup:
			t.AddRuleGroup(x)
		default:
			panic("Add expects either Rule or RuleGroup objects.")
		}
	}
}

func (t *RuleSet) AddRule(rules ...Rule) {
	for _, rule := range rules {
		rule := rule

		// Make sure rule is well-formed.
		rule.normalize()

		// Update maxNgram.
		N := rule.maxNGram()
		if N > t.maxNgram {
			t.maxNgram = N
		}

		t.rules = append(t.rules, &rule)
		t.root.AddRule(&rule)
	}
}

func (t *RuleSet) AddRuleGroup(ruleGroups ...RuleGroup) {
	for _, rg := range ruleGroups {
		t.AddRule(rg.ToList()...)
	}
}

// MatchRules will return a list of all matching rules. The rules are sorted by
// the match's "score". The best match will be first.
func (t *RuleSet) MatchRules(input string) (rules []*Rule) {
	input = t.sanitize(input)
	tokens := Tokenize(input, t.maxNgram)

	rules = t.root.Match(tokens)
	if len(rules) == 0 {
		return rules
	}

	// Check excludes.
	l := rules[:0]
	for _, r := range rules {
		if !r.isExcluded(tokens) {
			l = append(l, r)
		}
	}

	rules = l

	// Sort rules descending.
	sort.Slice(rules, func(i, j int) bool {
		return ruleLess(rules[j], rules[i])
	})

	// Update rule stats.
	if len(rules) > 0 {
		rules[0].FirstCount++
		for _, r := range rules {
			r.MatchCount++
		}
	}

	return rules
}

type Match struct {
	Tag        string
	Confidence float64 // In the range (0,1].
}

// Return a list of matches with confidence.
func (t *RuleSet) Match(input string) []Match {
	rules := t.MatchRules(input)
	if len(rules) == 0 {
		return []Match{}
	}
	if len(rules) == 1 {
		return []Match{{
			Tag:        rules[0].Tag,
			Confidence: 1,
		}}
	}

	// Create list of blocked tags.
	blocks := map[string]struct{}{}
	for _, rule := range rules {
		for _, tag := range rule.Blocks {
			blocks[tag] = struct{}{}
		}
	}

	// Remove rules for blocked tags.
	iOut := 0
	for _, rule := range rules {
		if _, ok := blocks[rule.Tag]; ok {
			continue
		}
		rules[iOut] = rule
		iOut++
	}
	rules = rules[:iOut]

	// Matches by index.
	matches := map[string]int{}
	out := []Match{}
	sum := float64(0)

	for _, rule := range rules {
		idx, ok := matches[rule.Tag]
		if !ok {
			idx = len(matches)
			matches[rule.Tag] = idx
			out = append(out, Match{Tag: rule.Tag})
		}
		out[idx].Confidence += float64(rule.score)
		sum += float64(rule.score)
	}

	for i := range out {
		out[i].Confidence /= sum
	}

	return out
}

// ListRules returns rules used in the ruleset sorted by the rules'
// FirstCount. This is the number of times the given rule was the best match to
// an input.
func (t *RuleSet) ListRules() []*Rule {
	sort.Slice(t.rules, func(i, j int) bool {
		if t.rules[j].FirstCount != t.rules[i].FirstCount {
			return t.rules[j].FirstCount < t.rules[i].FirstCount
		}

		if t.rules[j].MatchCount != t.rules[i].MatchCount {
			return t.rules[j].MatchCount < t.rules[i].MatchCount
		}

		return t.rules[j].Tag < t.rules[i].Tag
	})
	return t.rules
}

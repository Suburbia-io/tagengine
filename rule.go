package tagengine

type Rule struct {
	Tag      string
	Includes []string
	Excludes []string
	Blocks   []string // List of blocked tags.

	MatchCount int
	FirstCount int

	score int

	excludes map[string]struct{}
}

func NewRule(tag string) Rule {
	return Rule{Tag: tag}
}

func (r Rule) Inc(l ...string) Rule {
	return Rule{
		Tag:      r.Tag,
		Includes: l,
		Excludes: r.Excludes,
		Blocks:   r.Blocks,
	}
}

func (r Rule) Exc(l ...string) Rule {
	return Rule{
		Tag:      r.Tag,
		Includes: r.Includes,
		Excludes: l,
		Blocks:   r.Blocks,
	}
}

func (r Rule) Block(l ...string) Rule {
	return Rule{
		Tag:      r.Tag,
		Includes: r.Includes,
		Excludes: r.Excludes,
		Blocks:   l,
	}
}

func (rule *Rule) normalize() {
	sanitize := newSanitizer()

	for i, token := range rule.Includes {
		rule.Includes[i] = sanitize(token)
	}
	for i, token := range rule.Excludes {
		rule.Excludes[i] = sanitize(token)
	}

	sortTokens(rule.Includes)
	sortTokens(rule.Excludes)

	rule.excludes = map[string]struct{}{}
	for _, s := range rule.Excludes {
		rule.excludes[s] = struct{}{}
	}

	rule.score = rule.computeScore()
}

func (r Rule) maxNGram() int {
	max := 0
	for _, s := range r.Includes {
		n := ngramLength(s)
		if n > max {
			max = n
		}
	}
	for _, s := range r.Excludes {
		n := ngramLength(s)
		if n > max {
			max = n
		}
	}

	return max
}

func (r Rule) isExcluded(tokens []string) bool {
	// This is most often the case.
	if len(r.excludes) == 0 {
		return false
	}

	for _, s := range tokens {
		if _, ok := r.excludes[s]; ok {
			return true
		}
	}
	return false
}

func (r Rule) computeScore() (score int) {
	for _, token := range r.Includes {
		n := ngramLength(token)
		score += n * (n + 1) / 2
	}
	return score
}

func ruleLess(lhs, rhs *Rule) bool {
	// If scores differ, sort by score.
	if lhs.score != rhs.score {
		return lhs.score < rhs.score
	}

	// If include depth differs, sort by depth.
	lDepth := len(lhs.Includes)
	rDepth := len(rhs.Includes)

	if lDepth != rDepth {
		return lDepth < rDepth
	}

	// If exclude depth differs, sort by depth.
	lDepth = len(lhs.Excludes)
	rDepth = len(rhs.Excludes)

	if lDepth != rDepth {
		return lDepth < rDepth
	}

	// Sort alphabetically by includes.
	for i := range lhs.Includes {
		if lhs.Includes[i] != rhs.Includes[i] {
			return lhs.Includes[i] < rhs.Includes[i]
		}
	}

	// Sort by alphabetically by excludes.
	for i := range lhs.Excludes {
		if lhs.Excludes[i] != rhs.Excludes[i] {
			return lhs.Excludes[i] < rhs.Excludes[i]
		}
	}

	// Sort by tag.
	if lhs.Tag != rhs.Tag {
		return lhs.Tag < rhs.Tag
	}

	return false
}

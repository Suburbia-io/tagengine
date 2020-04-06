package tagengine

// A RuleGroup can be converted into a list of rules. Each rule will point to
// the same tag, and have the same exclude set.
type RuleGroup struct {
	Tag      string
	Includes [][]string
	Excludes []string
	Blocks   []string
}

func NewRuleGroup(tag string) RuleGroup {
	return RuleGroup{
		Tag:      tag,
		Includes: [][]string{},
		Excludes: []string{},
		Blocks:   []string{},
	}
}

func (g RuleGroup) Inc(l ...string) RuleGroup {
	return RuleGroup{
		Tag:      g.Tag,
		Includes: append(g.Includes, l),
		Excludes: g.Excludes,
		Blocks:   g.Blocks,
	}
}

func (g RuleGroup) Exc(l ...string) RuleGroup {
	return RuleGroup{
		Tag:      g.Tag,
		Includes: g.Includes,
		Excludes: l,
		Blocks:   g.Blocks,
	}
}

func (g RuleGroup) Block(l ...string) RuleGroup {
	return RuleGroup{
		Tag:      g.Tag,
		Includes: g.Includes,
		Excludes: g.Excludes,
		Blocks:   l,
	}
}

func (rg RuleGroup) ToList() (l []Rule) {
	for _, includes := range rg.Includes {
		l = append(l, Rule{
			Tag:      rg.Tag,
			Excludes: rg.Excludes,
			Includes: includes,
		})
	}
	return
}

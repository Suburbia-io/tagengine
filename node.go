package tagengine

import (
	"fmt"
	"strings"
)

type node struct {
	Token    string
	Matches  []*Rule
	Children map[string]*node
}

func (n *node) AddRule(r *Rule) {
	n.addRule(r, 0)
}

func (n *node) addRule(r *Rule, idx int) {
	if len(r.Includes) == idx {
		n.Matches = append(n.Matches, r)
		return
	}

	token := r.Includes[idx]

	child, ok := n.Children[token]
	if !ok {
		child = &node{
			Token:    token,
			Children: map[string]*node{},
		}
		n.Children[token] = child
	}

	child.addRule(r, idx+1)
}

// Note that tokens must be sorted. This is the case for tokens created from
// the tokenize function.
func (n *node) Match(tokens []string) (rules []*Rule) {
	return n.match(tokens, rules)
}

func (n *node) match(tokens []string, rules []*Rule) []*Rule {
	// Check for a match.
	if n.Matches != nil {
		rules = append(rules, n.Matches...)
	}

	if len(tokens) == 0 {
		return rules
	}

	// Attempt to match children.
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if child, ok := n.Children[token]; ok {
			rules = child.match(tokens[i+1:], rules)
		}
	}

	return rules
}

func (n *node) Dump() {
	n.dump(0)
}

func (n *node) dump(depth int) {
	indent := strings.Repeat(" ", 2*depth)
	tag := ""
	for _, m := range n.Matches {
		tag += " " + m.Tag
	}
	fmt.Printf("%s%s%s\n", indent, n.Token, tag)
	for _, child := range n.Children {
		child.dump(depth + 1)
	}
}

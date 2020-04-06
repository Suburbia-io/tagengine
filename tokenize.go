package tagengine

import (
	"sort"
	"strings"
)

var ignoreTokens = map[string]struct{}{}

func init() {
	// These on their own are ignored.
	tokens := []string{
		"`", `~`, `!`, `@`, `#`, `%`, `^`, `&`, `*`, `(`, `)`,
		`-`, `_`, `+`, `=`, `[`, `{`, `]`, `}`, `\`, `|`,
		`:`, `;`, `"`, `'`, `,`, `<`, `.`, `>`, `?`, `/`,
	}
	for _, s := range tokens {
		ignoreTokens[s] = struct{}{}
	}
}

func Tokenize(
	input string,
	maxNgram int,
) (
	tokens []string,
) {
	// Avoid duplicate ngrams.
	ignored := map[string]bool{}

	fields := strings.Fields(input)

	if len(fields) < maxNgram {
		maxNgram = len(fields)
	}

	for i := 1; i < maxNgram+1; i++ {
		jMax := len(fields) - i + 1

		for j := 0; j < jMax; j++ {
			ngram := strings.Join(fields[j:i+j], " ")
			if _, ok := ignoreTokens[ngram]; !ok {
				if _, ok := ignored[ngram]; !ok {
					tokens = append(tokens, ngram)
					ignored[ngram] = true
				}
			}
		}
	}

	sortTokens(tokens)

	return tokens
}

func sortTokens(tokens []string) {
	sort.Slice(tokens, func(i, j int) bool {
		if len(tokens[i]) != len(tokens[j]) {
			return len(tokens[i]) < len(tokens[j])
		}
		return tokens[i] < tokens[j]
	})
}

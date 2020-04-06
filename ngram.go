package tagengine

import "unicode"

func ngramLength(s string) int {
	N := len(s)
	i := 0
	count := 0

	for {
		// Eat spaces.
		for i < N && unicode.IsSpace(rune(s[i])) {
			i++
		}

		// Done?
		if i == N {
			break
		}

		// Non-space!
		count++

		// Eat non-spaces.
		for i < N && !unicode.IsSpace(rune(s[i])) {
			i++
		}
	}
	return count
}

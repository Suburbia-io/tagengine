package tagengine

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func newSanitizer() func(...string) string {
	diactricsFix := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	return func(l ...string) string {

		s := strings.Join(l, " ")

		// Lowercase.
		s = strings.ToLower(s)

		// Remove apostrophes.
		s = strings.ReplaceAll(s, "ß", "ss")
		s = strings.ReplaceAll(s, "'s", "s")
		s = strings.ReplaceAll(s, "`s", "s")
		s = strings.ReplaceAll(s, "´s", "s")

		// Remove diacritics.
		if out, _, err := transform.String(diactricsFix, s); err == nil {
			s = out
		}

		// Clean spaces.
		s = spaceNumbers(s)
		s = addSpaces(s)
		s = collapseSpaces(s)

		return s
	}
}

func spaceNumbers(s string) string {
	if len(s) == 0 {
		return s
	}

	isDigit := func(b rune) bool {
		switch b {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return true
		}
		return false
	}

	b := strings.Builder{}

	var first rune
	for _, c := range s {
		first = c
		break
	}

	digit := isDigit(first)

	// Range over runes.
	for _, c := range s {
		thisDigit := isDigit(c)
		if thisDigit != digit {
			b.WriteByte(' ')
			digit = thisDigit
		}
		b.WriteRune(c)
	}

	return b.String()
}

func addSpaces(s string) string {
	needsSpace := func(r rune) bool {
		switch r {
		case '`', '~', '!', '@', '#', '%', '^', '&', '*', '(', ')',
			'-', '_', '+', '=', '[', '{', ']', '}', '\\', '|',
			':', ';', '"', '\'', ',', '<', '.', '>', '?', '/':
			return true
		}
		return false
	}

	b := strings.Builder{}

	// Range over runes.
	for _, r := range s {
		if needsSpace(r) {
			b.WriteRune(' ')
			b.WriteRune(r)
			b.WriteRune(' ')
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func collapseSpaces(s string) string {
	// Trim leading and trailing spaces.
	s = strings.TrimSpace(s)

	b := strings.Builder{}
	wasSpace := false

	// Range over runes.
	for _, c := range s {
		if unicode.IsSpace(c) {
			wasSpace = true
			continue
		} else if wasSpace {
			wasSpace = false
			b.WriteRune(' ')
		}
		b.WriteRune(c)
	}

	return b.String()
}

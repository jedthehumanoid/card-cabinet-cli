package slug

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

// From returns "this format" from "this-format"
// Note that slugging is lossy, returning from slug
// might differ from original
func From(s string) string {
	return strings.Replace(s, "-", " ", -1)
}

// To returns a slug out of string
// Makes lowercase, removes international accents, turns spaces and
// other non alphanumeric characters to dashes, but keeps slashes
func To(s string) string {
	var re = regexp.MustCompile("[^a-z0-9/]+")

	s = strings.ToLower(s)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	s, _, _ = transform.String(t, s)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

package main

import (
	"regexp"
)

// (?ms) Flags=multiline, dot matches newline
// ^---$ Three dashes as only thing on line
// .* Everything in between
// ^---$ Three dashes as only thing on line again
var yamlregexp = regexp.MustCompile("(?ms)^---$.*^---$")

var tomlregexp = regexp.MustCompile("(?ms)^+++$.*^+++$")

// HasYAMLFrontmatter returns true frontmatter is present
func HasYAMLFrontmatter(b []byte) bool {
	return yamlregexp.Match(b)
}

// GetYAMLFrontmatter returns everthing in between and including ---
func GetYAMLFrontmatter(b []byte) []byte {
	fm := yamlregexp.Find(b)
	if fm == nil {
		return []byte{}
	}
	return fm
}

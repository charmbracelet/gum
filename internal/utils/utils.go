package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LipglossTruncate truncates a given line based on its lipgloss width.
func LipglossTruncate(s string, width int) string {
	var i int
	for i = 0; i < len(s) && lipgloss.Width(s[:i]) < width; i++ {
	} //revive:disable-line:empty-block
	return s[:i]
}

// LipglossLengthPadding calculated calculates how much padding a string is given by a style.
func LipglossLengthPadding(s string, style lipgloss.Style) (int, int) {
	render := style.Render(s)
	before := strings.Index(render, s)
	after := len(render) - len(s) - before
	return before, after
}

// UniqueStrings returns a list of unique strings given a list of strings.
func UniqueStrings(strings []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, s := range strings {
		if _, uniq := keys[s]; !uniq {
			keys[s] = true
			list = append(list, s)
		}
	}
	return list
}

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
func LipglossPadding(style lipgloss.Style) (int, int) {
	render := style.Render(" ")
	before := strings.Index(render, " ")
	after := len(render) - len(" ") - before
	return before, after
}

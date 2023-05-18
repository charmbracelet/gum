package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LipglossPadding calculates how much padding a string is given by a style.
func LipglossPadding(style lipgloss.Style) (int, int) {
	render := style.Render(" ")
	before := strings.Index(render, " ")
	after := len(render) - len(" ") - before
	return before, after
}

package style

import (
	"github.com/charmbracelet/gum/internal/decode"
	"github.com/charmbracelet/lipgloss"
)

// ToLipgloss takes a Styles flag set and returns the corresponding
// lipgloss.Style.
func (s Styles) ToLipgloss() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Background)).
		Foreground(lipgloss.Color(s.Foreground)).
		BorderBackground(lipgloss.Color(s.BorderBackground)).
		BorderForeground(lipgloss.Color(s.BorderForeground)).
		Align(decode.Align[s.Align]).
		Border(border[s.Border]).
		Height(s.Height).
		Width(s.Width).
		Margin(parseMargin(s.Margin)).
		Padding(parsePadding(s.Padding)).
		Bold(s.Bold).
		Faint(s.Faint).
		Italic(s.Italic).
		Strikethrough(s.Strikethrough).
		Underline(s.Underline)
}

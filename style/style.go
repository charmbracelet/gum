package style

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() {
	text := strings.Join(o.Text, "\n")

	fmt.Println(lipgloss.NewStyle().
		Foreground(lipgloss.Color(o.Foreground)).
		Background(lipgloss.Color(o.Background)).
		BorderBackground(lipgloss.Color(o.BorderBackground)).
		BorderForeground(lipgloss.Color(o.BorderForeground)).
		Align(align[o.Align]).
		Bold(o.Bold).
		Border(border[o.Border]).
		Margin(parseMargin(o.Margin)).
		Padding(parsePadding(o.Padding)).
		Height(o.Height).
		Width(o.Width).
		Faint(o.Faint).
		Italic(o.Italic).
		Strikethrough(o.Strikethrough).
		Render(text))
}

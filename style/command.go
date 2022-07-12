package style

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/gum/internal/decode"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() {
	text := strings.Join(o.Text, "\n")

	fmt.Println(lipgloss.NewStyle().
		Foreground(lipgloss.Color(o.Style.Foreground)).
		Background(lipgloss.Color(o.Style.Background)).
		BorderBackground(lipgloss.Color(o.Style.BorderBackground)).
		BorderForeground(lipgloss.Color(o.Style.BorderForeground)).
		Align(decode.Align[o.Style.Align]).
		Bold(o.Style.Bold).
		Border(border[o.Style.Border]).
		Margin(parseMargin(o.Style.Margin)).
		Padding(parsePadding(o.Style.Padding)).
		Height(o.Style.Height).
		Width(o.Style.Width).
		Faint(o.Style.Faint).
		Italic(o.Style.Italic).
		Strikethrough(o.Style.Strikethrough).
		Render(text))
}

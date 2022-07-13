package style

import (
	"github.com/charmbracelet/gum/internal/decode"
	"github.com/charmbracelet/lipgloss"
)

// Styles is a flag set of possible styles.
// It corresponds to the available options in the lipgloss.Style struct.
type Styles struct {
	// Colors
	Background string `help:"Background color of the ${name=element}" default:"${defaultBackground}" hidden:""`
	Foreground string `help:"color of the ${name=element}" default:"${defaultForeground}"`

	// Border
	Border           string `help:"Border style to apply" enum:"none,hidden,normal,rounded,thick,double" default:"none" hidden:""`
	BorderBackground string `help:"Border background color" hidden:""`
	BorderForeground string `help:"Border foreground color" hidden:""`

	// Layout
	Align   string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left" hidden:""`
	Height  int    `help:"Height of output" hidden:""`
	Width   int    `help:"Width of output" hidden:""`
	Margin  string `help:"Margin to apply around the text." default:"0 0" hidden:""`
	Padding string `help:"Padding to apply around the text." default:"0 0" hidden:""`

	// Format
	Bold          bool `help:"Apply bold formatting" hidden:""`
	Faint         bool `help:"Apply faint formatting" hidden:""`
	Italic        bool `help:"Apply italic formatting" hidden:""`
	Strikethrough bool `help:"Apply strikethrough formatting" hidden:""`
}

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
		Strikethrough(s.Strikethrough)
}

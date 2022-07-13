package style

import (
	"github.com/charmbracelet/gum/internal/decode"
	"github.com/charmbracelet/lipgloss"
)

// Styles is a flag set of possible styles.
// It corresponds to the available options in the lipgloss.Style struct.
type Styles struct {
	// Colors
	Background string `help:"Background color of the ${name=element}" default:"${defaultBackground}" hidden:"" group:"Style Flags"`
	Foreground string `help:"color of the ${name=element}" default:"${defaultForeground}" group:"Style Flags"`

	// Border
	Border           string `help:"Border style to apply" enum:"none,hidden,normal,rounded,thick,double" default:"none" hidden:"" group:"Style Flags"`
	BorderBackground string `help:"Border background color" hidden:"" group:"Style Flags"`
	BorderForeground string `help:"Border foreground color" hidden:"" group:"Style Flags"`

	// Layout
	Align   string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left" hidden:"" group:"Style Flags"`
	Height  int    `help:"Height of output" hidden:"" group:"Style Flags"`
	Width   int    `help:"Width of output" hidden:"" group:"Style Flags"`
	Margin  string `help:"Margin to apply around the text." default:"0 0" hidden:"" group:"Style Flags"`
	Padding string `help:"Padding to apply around the text." default:"0 0" hidden:""`

	// Format
	Bold          bool `help:"Apply bold formatting" hidden:"" group:"Style Flags"`
	Faint         bool `help:"Apply faint formatting" hidden:"" group:"Style Flags"`
	Italic        bool `help:"Apply italic formatting" hidden:"" group:"Style Flags"`
	Strikethrough bool `help:"Apply strikethrough formatting" hidden:"" group:"Style Flags"`
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

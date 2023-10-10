package style

import "github.com/charmbracelet/lipgloss"

// Border maps strings to `lipgloss.Border`s.
var Border map[string]lipgloss.Border = map[string]lipgloss.Border{
	"double":  lipgloss.DoubleBorder(),
	"hidden":  lipgloss.HiddenBorder(),
	"none":    {},
	"normal":  lipgloss.NormalBorder(),
	"rounded": lipgloss.RoundedBorder(),
	"thick":   lipgloss.ThickBorder(),
}

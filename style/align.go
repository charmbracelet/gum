package style

import "github.com/charmbracelet/lipgloss"

// align maps strings to `lipgloss.Position`s
var align = map[string]lipgloss.Position{
	"center": lipgloss.Center,
	"left":   lipgloss.Left,
	"top":    lipgloss.Top,
	"bottom": lipgloss.Bottom,
	"right":  lipgloss.Right,
}

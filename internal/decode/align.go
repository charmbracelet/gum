package decode

import "github.com/charmbracelet/lipgloss"

// Align maps strings to `lipgloss.Position`s
var Align = map[string]lipgloss.Position{
	"center": lipgloss.Center,
	"left":   lipgloss.Left,
	"top":    lipgloss.Top,
	"bottom": lipgloss.Bottom,
	"right":  lipgloss.Right,
}

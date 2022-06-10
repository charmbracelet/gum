package style

import "github.com/charmbracelet/lipgloss"

var alignment = map[string]lipgloss.Position{
	"center": lipgloss.Center,
	"left":   lipgloss.Left,
	"top":    lipgloss.Top,
	"bottom": lipgloss.Bottom,
	"right":  lipgloss.Right,
}

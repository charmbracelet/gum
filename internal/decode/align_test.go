package decode

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAlign(t *testing.T) {
	for name, want := range map[string]lipgloss.Position{
		"center": lipgloss.Center,
		"left":   lipgloss.Left,
		"top":    lipgloss.Top,
		"bottom": lipgloss.Bottom,
		"right":  lipgloss.Right,
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := Align[name]
			if !ok {
				t.Fatalf("Align map missing key %q", name)
			}
			if got != want {
				t.Errorf("Align[%q] = %v, want %v", name, got, want)
			}
		})
	}
}

func TestAlignSize(t *testing.T) {
	if got, want := len(Align), 5; got != want {
		t.Errorf("Align map size = %d, want %d", got, want)
	}
}

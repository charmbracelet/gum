package input

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTextinputWidth(t *testing.T) {
	t.Parallel()
	padding := []int{0, 2, 0, 2}
	prompt := "> "
	want := 80 - 1 - lipgloss.Width(prompt) - padding[1] - padding[3]
	got := textinputWidth(80, prompt, padding)
	if got != want {
		t.Fatalf("textinputWidth() = %d, want %d", got, want)
	}
	if textinputWidth(3, prompt, padding) != 0 {
		t.Fatal("expected non-negative width clamped to 0")
	}
}

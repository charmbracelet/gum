package input

import (
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

func textinputWidth(termWidth int, prompt string, padding []int) int {
	w := termWidth - 1 - lipgloss.Width(prompt) - padding[1] - padding[3]
	if w < 0 {
		return 0
	}
	return w
}

func applyTerminalWidth(ti *textinput.Model, padding []int) bool {
	w, _, err := term.GetSize(os.Stderr.Fd())
	if err != nil || w < 1 {
		return false
	}
	ti.Width = textinputWidth(w, ti.Prompt, padding)
	return true
}

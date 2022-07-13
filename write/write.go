// Package write provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/tree/master/textarea
//
// It can be used to ask the user to write some long form of text
// (multi-line) input. The text the user entered will be sent to stdout.
//
//   $ gum write > output.text
//
package write

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	quitting bool
	textarea textarea.Model
}

func (m model) Init() tea.Cmd { return textarea.Blink }
func (m model) View() string {
	if m.quitting {
		return ""
	}
	return m.textarea.View()
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC, tea.KeyCtrlD:
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

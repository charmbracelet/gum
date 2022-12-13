// Package write provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/tree/master/textarea
//
// It can be used to ask the user to write some long form of text (multi-line)
// input. The text the user entered will be sent to stdout.
// Text entry is completed with CTRL+D and aborted with CTRL+C or Escape.
//
// $ gum write > output.text
package write

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	aborted     bool
	header      string
	headerStyle lipgloss.Style
	quitting    bool
	textarea    textarea.Model
}

func (m model) Init() tea.Cmd { return textarea.Blink }
func (m model) View() string {
	if m.quitting {
		return ""
	}

	// Display the header above the text area if it is not empty.
	if m.header != "" {
		header := m.headerStyle.Render(m.header)
		return lipgloss.JoinVertical(lipgloss.Left, header, m.textarea.View())
	}

	return m.textarea.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case "esc", "ctrl+d":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

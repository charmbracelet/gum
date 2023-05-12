// Package input provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/tree/master/textinput
//
// It can be used to prompt the user for some input. The text the user entered
// will be sent to stdout.
//
// $ gum input --placeholder "What's your favorite gum?" > answer.text
package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	autoWidth   bool
	header      string
	headerStyle lipgloss.Style
	textinput   textinput.Model
	quitting    bool
	aborted     bool
}

func (m model) Init() tea.Cmd { return textinput.Blink }
func (m model) View() string {
	if m.quitting {
		return ""
	}

	if m.header != "" {
		header := m.headerStyle.Render(m.header)
		return lipgloss.JoinVertical(lipgloss.Left, header, m.textinput.View())
	}

	return m.textinput.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.autoWidth {
			m.textinput.Width = msg.Width - lipgloss.Width(m.textinput.Prompt) - 1
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			m.aborted = true
			return m, tea.Quit
		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

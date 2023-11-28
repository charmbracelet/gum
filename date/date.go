// Package date provides a shell script interface for picking a date.
//
// The date the user selected will be sent to stdout in ISO-8601 format:
// YYYY-MM-DD.
//
// $ gum date --value 2023-11-28 > date.text
package date

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/timeout"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	header      string
	headerStyle lipgloss.Style
	picker      *picker
	quitting    bool
	aborted     bool
	timeout     time.Duration
	hasTimeout  bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		timeout.Init(m.timeout, nil),
	)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.header != "" {
		header := m.headerStyle.Render(m.header)
		return lipgloss.JoinVertical(lipgloss.Left, header, m.picker.View())
	}

	return m.picker.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.quitting = true
			m.aborted = true
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	// case tea.WindowSizeMsg:
	// 	if m.autoWidth {
	// 		m.textinput.Width = msg.Width - lipgloss.Width(m.textinput.Prompt) - 1
	// 	}
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
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

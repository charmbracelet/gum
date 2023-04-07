// Package file provides an interface to pick a file from a folder (tree).
// The user is provided a file manager-like interface to navigate, to
// select a file.
//
// Let's pick a file from the current directory:
//
// $ gum file
// $ gum file .
//
// Let's pick a file from the home directory:
//
// $ gum file $HOME
package file

import (
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/timeout"
	"time"
)

type model struct {
	filepicker   filepicker.Model
	selectedPath string
	aborted      bool
	quitting     bool
	timeout      time.Duration
	hasTimeout   bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		timeout.Init(m.timeout, nil),
		m.filepicker.Init(),
	)

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		}
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.quitting = true
			m.aborted = true
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	}
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedPath = path
		m.quitting = true
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return m.filepicker.View()
}

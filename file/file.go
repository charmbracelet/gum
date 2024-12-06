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
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/timeout"
	"github.com/charmbracelet/lipgloss"
)

type keymap filepicker.KeyMap

var keyQuit = key.NewBinding(
	key.WithKeys("esc", "q"),
	key.WithHelp("esc", "close"),
)

var keyAbort = key.NewBinding(
	key.WithKeys("ctrl+c"),
	key.WithHelp("ctrl+c", "abort"),
)

func defaultKeymap() keymap {
	km := filepicker.DefaultKeyMap()
	km.Down.SetHelp("↓", "down")
	km.Up.SetHelp("↑", "up")
	return keymap(km)
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		keyQuit,
		k.Select,
	}
}

type model struct {
	filepicker   filepicker.Model
	selectedPath string
	aborted      bool
	timedOut     bool
	quitting     bool
	timeout      time.Duration
	hasTimeout   bool
	showHelp     bool
	help         help.Model
	keymap       keymap
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		timeout.Init(m.timeout, nil),
		m.filepicker.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.showHelp {
			m.filepicker.Height -= lipgloss.Height(m.helpView())
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyAbort):
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keyQuit):
			m.quitting = true
			return m, tea.Quit
		}
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.quitting = true
			m.timedOut = true
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
	if !m.showHelp {
		return m.filepicker.View()
	}
	return m.filepicker.View() + m.helpView()
}

func (m model) helpView() string {
	return "\n" + m.help.View(m.keymap)
}

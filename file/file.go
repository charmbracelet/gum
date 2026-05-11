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
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	return keymap(km)
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↓↑", "navigate"),
		),
		keyQuit,
		k.Select,
	}
}

type model struct {
	header       string
	headerStyle  lipgloss.Style
	filepicker   filepicker.Model
	selectedPath string
	quitting     bool
	showHelp     bool
	padding      []int
	help         help.Model
	keymap       keymap
}

func (m model) Init() tea.Cmd { return m.filepicker.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		height := msg.Height - m.padding[0] - m.padding[2]
		if m.showHelp {
			height -= lipgloss.Height(m.helpView())
		}
		m.filepicker.SetHeight(height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyAbort):
			m.quitting = true
			return m, tea.Interrupt
		case key.Matches(msg, keyQuit):
			m.quitting = true
			return m, tea.Quit
		}
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
	var parts []string
	if m.header != "" {
		parts = append(parts, m.headerStyle.Render(m.header))
	}
	parts = append(parts, m.filepicker.View())
	if m.showHelp {
		parts = append(parts, m.helpView())
	}
	return lipgloss.NewStyle().
		Padding(m.padding...).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			parts...,
		))
}

func (m model) helpView() string {
	return m.help.View(m.keymap)
}

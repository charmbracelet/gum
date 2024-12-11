// Package table provides a shell script interface for the table bubble.
// https://github.com/charmbracelet/bubbles/tree/master/table
//
// It is useful to render tabular (CSV) data in a terminal and allows
// the user to select a row from the table.
//
// Let's render a table of gum flavors:
//
// $ gum table <<< "Flavor,Price\nStrawberry,$0.50\nBanana,$0.99\nCherry,$0.75"
//
//	Flavor      Price
//	Strawberry  $0.50
//	Banana      $0.99
//	Cherry      $0.75
package table

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type keymap struct {
	Navigate,
	Select,
	Quit,
	Abort key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Navigate,
		k.Select,
		k.Quit,
	}
}

func defaultKeymap() keymap {
	return keymap{
		Navigate: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↓↑", "navigate"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+q", "q"),
			key.WithHelp("esc", "quit"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
	}
}

type model struct {
	table    table.Model
	selected table.Row
	quitting bool
	showHelp bool
	help     help.Model
	keymap   keymap
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		km := m.keymap
		switch {
		case key.Matches(msg, km.Select):
			m.selected = m.table.SelectedRow()
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Abort):
			m.quitting = true
			return m, tea.Interrupt
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	s := m.table.View()
	if m.showHelp {
		s += "\n" + m.help.View(m.keymap)
	}
	return s
}

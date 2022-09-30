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
//  Flavor      Price
//  Strawberry  $0.50
//  Banana      $0.99
//  Cherry      $0.75
//
package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	table    table.Model
	selected table.Row
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.selected = m.table.SelectedRow()
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return m.table.View()
}

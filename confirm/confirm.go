// Package confirm provides an interface to ask a user to confirm an action.
// The user is provided with an interface to choose an affirmative or negative
// answer, which is then reflected in the exit code for use in scripting.
//
// If the user selects the affirmative answer, the program exits with 0. If the
// user selects the negative answer, the program exits with 1.
//
// I.e. confirm if the user wants to delete a file
//
//   $ gum confirm "Are you sure?" && rm file.txt
//
package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	prompt      string
	affirmative string
	negative    string
	vertical    bool
	quitting    bool
	selected    int

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "ctrl+p":
			m.selected = m.selected - 1
			if m.selected < 0 {
				m.selected = 1
			}
		case "right", "l", "ctrl+n":
			m.selected = m.selected + 1
			if m.selected > 1 {
				m.selected = 0
			}
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c", "q":
			m.selected = 1
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	joinFunc := lipgloss.JoinHorizontal
	if m.vertical {
		joinFunc = lipgloss.JoinVertical
	}

	var aff, neg string

	if m.selected == 0 {
		aff = m.selectedStyle.Render(m.affirmative)
		neg = m.unselectedStyle.Render(m.negative)
	} else {
		aff = m.unselectedStyle.Render(m.affirmative)
		neg = m.selectedStyle.Render(m.negative)
	}

	return lipgloss.JoinVertical(lipgloss.Center, m.promptStyle.Render(m.prompt), joinFunc(lipgloss.Left, aff, neg))
}

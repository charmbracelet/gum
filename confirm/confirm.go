// Package confirm provides an interface to ask a user to confirm an action.
// The user is provided with an interface to choose an affirmative or negative
// answer, which is then reflected in the exit code for use in scripting.
//
// If the user selects the affirmative answer, the program exits with 0. If the
// user selects the negative answer, the program exits with 1.
//
// I.e. confirm if the user wants to delete a file
//
// $ gum confirm "Are you sure?" && rm file.txt
package confirm

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	prompt      string
	affirmative string
	negative    string
	quitting    bool
	timeout     int

	confirmation bool

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q", "n", "N":
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		}
	}

	return updateChoices(msg, m)
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "ctrl+p", "tab",
			"right", "l", "ctrl+n", "shift+tab":
			m.confirmation = !m.confirmation
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case "y", "Y":
			m.quitting = true
			m.confirmation = true
			return m, tea.Quit
		}
	case tickMsg:
		if m.timeout == 0 {
			m.quitting = true
			m.confirmation = false
			return m, tea.Quit
		}
		m.timeout -= 1
		return m, tick()
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var aff, neg, negativeMessage string

	negativeMessage = fmt.Sprintf("%s(%d)", m.negative, m.timeout)

	if m.confirmation {
		aff = m.selectedStyle.Render(m.affirmative)
		neg = m.unselectedStyle.Render(negativeMessage)
	} else {
		aff = m.unselectedStyle.Render(m.affirmative)
		neg = m.selectedStyle.Render(negativeMessage)
	}

	return lipgloss.JoinVertical(lipgloss.Center, m.promptStyle.Render(m.prompt), lipgloss.JoinHorizontal(lipgloss.Left, aff, neg))
}

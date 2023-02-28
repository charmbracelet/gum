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
	aborted     bool
	hasTimeout  bool
	timeout     time.Duration

	confirmation bool

	defaultSelection bool

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

const tickInterval = time.Second

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	if m.timeout > 0 {
		return tick()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			fallthrough
		case "esc":
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		case "q", "n", "N":
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		case "left", "h", "ctrl+p", "tab",
			"right", "l", "ctrl+n", "shift+tab":
			if m.negative == "" {
				break
			}
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
		if m.timeout <= 0 {
			m.quitting = true
			m.confirmation = m.defaultSelection
			return m, tea.Quit
		}
		m.timeout -= tickInterval
		return m, tick()
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var aff, neg, timeout, affirmativeTimeout, negativeTimeout string

	if m.hasTimeout {
		timeout = fmt.Sprintf(" (%d)", max(0, int(m.timeout.Seconds())))
	}

	// set timer based on defaultSelection
	if m.defaultSelection {
		affirmativeTimeout = m.affirmative + timeout
		negativeTimeout = m.negative
	} else {
		affirmativeTimeout = m.affirmative
		negativeTimeout = m.negative + timeout
	}

	if m.confirmation {
		aff = m.selectedStyle.Render(affirmativeTimeout)
		neg = m.unselectedStyle.Render(negativeTimeout)
	} else {
		aff = m.unselectedStyle.Render(affirmativeTimeout)
		neg = m.selectedStyle.Render(negativeTimeout)
	}

	// If the option is intentionally empty, do not show it.
	if m.negative == "" {
		neg = ""
	}

	return lipgloss.JoinVertical(lipgloss.Center, m.promptStyle.Render(m.prompt), lipgloss.JoinHorizontal(lipgloss.Left, aff, neg))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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

type selectionType int

// Type of user selection.:w
const (
	Confirmed selectionType = iota
	Negative
	Cancel
)

type model struct {
	prompt     string
	options    []string
	quitting   bool
	hasTimeout bool
	timeout    time.Duration

	selectedOption   selectionType
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

func (m *model) NextOption() {
	for {
		m.selectedOption++
		if (int)(m.selectedOption) >= len(m.options) {
			m.selectedOption = 0
		}
		if m.options[m.selectedOption] != "" {
			break
		}
	}
}

func (m *model) PrevOption() {
	for {
		m.selectedOption--
		if (int)(m.selectedOption) < 0 {
			m.selectedOption = selectionType(len(m.options) - 1)
		}
		if m.options[m.selectedOption] != "" {
			break
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "n", "N":
			if m.options[Negative] != "" {
				m.selectedOption = Negative
				m.quitting = true
			}
		case "ctrl+c":
			m.selectedOption = Cancel
			m.quitting = true
			return m, tea.Quit
		case "right", "l", "ctrl+p", "tab":
			m.NextOption()
		case
			"left", "h", "ctrl+n", "shift+tab":
			m.PrevOption()
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case "y", "Y":
			m.quitting = true
			m.selectedOption = Confirmed
			return m, tea.Quit
		}
	case tickMsg:
		if m.timeout <= 0 {
			m.quitting = true
			if !m.defaultSelection && m.options[Negative] != "" {
				m.selectedOption = Negative
			} else {
				m.selectedOption = Confirmed
			}
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

	var timeout string

	if m.hasTimeout {
		timeout = fmt.Sprintf(" (%d)", max(0, int(m.timeout.Seconds())))
	}

	var renderedOptions []string
	for i, value := range m.options {
		if m.options[i] == "" {
			continue
		}
		if m.hasTimeout {
			if (m.defaultSelection && i == (int)(Confirmed)) || (!m.defaultSelection && i == (int)(Negative)) {
				value = value + timeout
			}
		}
		if (int)(m.selectedOption) == i {
			renderedOptions = append(renderedOptions, m.selectedStyle.Render(value))
		} else {
			renderedOptions = append(renderedOptions, m.unselectedStyle.Render(value))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, m.promptStyle.Render(m.prompt), lipgloss.JoinHorizontal(lipgloss.Left, renderedOptions...))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

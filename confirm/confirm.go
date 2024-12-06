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
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/gum/timeout"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func defaultKeymap(affirmative, negative string) keymap {
	return keymap{
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Negative: key.NewBinding(
			key.WithKeys("n", "N", "q"),
			key.WithHelp("n", negative),
		),
		Affirmative: key.NewBinding(
			key.WithKeys("y", "Y"),
			key.WithHelp("y", affirmative),
		),
		Toggle: key.NewBinding(
			key.WithKeys(
				"left",
				"h",
				"ctrl+n",
				"shift+tab",
				"right",
				"l",
				"ctrl+p",
				"tab",
			),
			key.WithHelp("←/→", "toggle"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
	}
}

type keymap struct {
	Abort       key.Binding
	Quit        key.Binding
	Negative    key.Binding
	Affirmative key.Binding
	Toggle      key.Binding
	Submit      key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle, k.Submit, k.Affirmative, k.Negative}
}

type model struct {
	prompt      string
	affirmative string
	negative    string
	quitting    bool
	aborted     bool
	hasTimeout  bool
	showHelp    bool
	help        help.Model
	keys        keymap
	timeout     time.Duration

	confirmation bool
	timedOut     bool

	defaultSelection bool

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return timeout.Init(m.timeout, m.defaultSelection)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Abort):
			m.confirmation = false
			m.aborted = true
			fallthrough
		case key.Matches(msg, m.keys.Quit):
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Negative):
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Toggle):
			if m.negative == "" {
				break
			}
			m.confirmation = !m.confirmation
		case key.Matches(msg, m.keys.Submit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Affirmative):
			m.quitting = true
			m.confirmation = true
			return m, tea.Quit
		}
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.quitting = true
			m.confirmation = m.defaultSelection
			m.timedOut = true
			return m, tea.Quit
		}

		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var aff, neg, timeoutStrYes, timeoutStrNo string
	timeoutStrNo = ""
	timeoutStrYes = ""
	if m.hasTimeout {
		if m.defaultSelection {
			timeoutStrYes = timeout.Str(m.timeout)
		} else {
			timeoutStrNo = timeout.Str(m.timeout)
		}
	}

	if m.confirmation {
		aff = m.selectedStyle.Render(m.affirmative + timeoutStrYes)
		neg = m.unselectedStyle.Render(m.negative + timeoutStrNo)
	} else {
		aff = m.unselectedStyle.Render(m.affirmative + timeoutStrYes)
		neg = m.selectedStyle.Render(m.negative + timeoutStrNo)
	}

	// If the option is intentionally empty, do not show it.
	if m.negative == "" {
		neg = ""
	}

	if m.showHelp {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.promptStyle.Render(m.prompt)+"\n",
			lipgloss.JoinHorizontal(lipgloss.Left, aff, neg),
			"\n"+m.help.View(m.keys),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.promptStyle.Render(m.prompt)+"\n",
		lipgloss.JoinHorizontal(lipgloss.Left, aff, neg),
	)
}

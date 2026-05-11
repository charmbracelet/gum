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
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

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
			key.WithHelp("←→", "toggle"),
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
	showHelp    bool
	help        help.Model
	keys        keymap

	showOutput   bool
	confirmation bool

	defaultSelection bool

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	padding         []int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Abort):
			m.confirmation = false
			return m, tea.Interrupt
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
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var aff, neg string

	if m.confirmation {
		aff = m.selectedStyle.Render(m.affirmative)
		neg = m.unselectedStyle.Render(m.negative)
	} else {
		aff = m.unselectedStyle.Render(m.affirmative)
		neg = m.selectedStyle.Render(m.negative)
	}

	// If the option is intentionally empty, do not show it.
	if m.negative == "" {
		neg = ""
	}

	parts := []string{
		m.promptStyle.Render(m.prompt) + "\n",
		lipgloss.JoinHorizontal(lipgloss.Left, aff, neg),
	}

	if m.showHelp {
		parts = append(parts, "", m.help.View(m.keys))
	}

	return lipgloss.NewStyle().
		Padding(m.padding...).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			parts...,
		))
}

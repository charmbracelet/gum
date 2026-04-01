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
	extra       []string
	quitting    bool
	showHelp    bool
	help        help.Model
	keys        keymap

	showOutput   bool
	confirmation bool
	selected     int // 0 = affirmative, 1 = negative, 2+ = extra options

	defaultSelection bool

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	padding         []int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) totalOptions() int {
	n := 1 // affirmative
	if m.negative != "" {
		n++
	}
	n += len(m.extra)
	return n
}

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
			if len(m.extra) == 0 {
				m.confirmation = false
				m.selected = 1
				m.quitting = true
				return m, tea.Quit
			}
			// When extra options exist, n/N just toggles like other keys
			m.confirmation = !m.confirmation
		case key.Matches(msg, m.keys.Toggle):
			if m.negative == "" && len(m.extra) == 0 {
				break
			}
			total := m.totalOptions()
			m.selected = (m.selected + 1) % total
			m.confirmation = m.selected == 0
		case key.Matches(msg, m.keys.Submit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Affirmative):
			if len(m.extra) == 0 {
				m.quitting = true
				m.confirmation = true
				m.selected = 0
				return m, tea.Quit
			}
			// When extra options exist, y/Y just toggles
			m.confirmation = !m.confirmation
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	// Build the list of all options with their rendered styles.
	var options []string

	// Affirmative (index 0)
	if m.selected == 0 {
		options = append(options, m.selectedStyle.Render(m.affirmative))
	} else {
		options = append(options, m.unselectedStyle.Render(m.affirmative))
	}

	// Negative (index 1, if not empty)
	if m.negative != "" {
		if m.selected == 1 {
			options = append(options, m.selectedStyle.Render(m.negative))
		} else {
			options = append(options, m.unselectedStyle.Render(m.negative))
		}
	}

	// Extra options (index 2+)
	baseIdx := 1
	if m.negative != "" {
		baseIdx = 2
	}
	for i, extra := range m.extra {
		if m.selected == baseIdx+i {
			options = append(options, m.selectedStyle.Render(extra))
		} else {
			options = append(options, m.unselectedStyle.Render(extra))
		}
	}

	parts := []string{
		m.promptStyle.Render(m.prompt) + "\n",
		lipgloss.JoinHorizontal(lipgloss.Left, options...),
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

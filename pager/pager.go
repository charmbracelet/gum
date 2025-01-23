// Package pager provides a pager (similar to less) for the terminal.
//
// $ cat file.txt | gum pager
package pager

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	Home,
	End,
	Search,
	NextMatch,
	PrevMatch,
	Abort,
	Quit,
	ConfirmSearch,
	CancelSearch key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding {
	return nil
}

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("left", "down", "up", "rigth"),
			key.WithHelp("←↓↑→", "navigate"),
		),
		k.Quit,
		k.Search,
		k.NextMatch,
		k.PrevMatch,
	}
}

func defaultKeymap() keymap {
	return keymap{
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("h", "home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G", "end"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		PrevMatch: key.NewBinding(
			key.WithKeys("p", "N"),
			key.WithHelp("N", "previous match"),
		),
		NextMatch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next match"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("esc", "quit"),
		),
		ConfirmSearch: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		CancelSearch: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+d", "esc"),
			key.WithHelp("ctrl+c", "cancel"),
		),
	}
}

type model struct {
	viewport        viewport.Model
	help            help.Model
	showLineNumbers bool
	lineNumberStyle lipgloss.Style
	search          search
	maxWidth        int
	keymap          keymap
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.processText(msg)
		m.search.input.Width = msg.Width
	case tea.KeyMsg:
		return m.keyHandler(msg)
	}

	m.keymap.PrevMatch.SetEnabled(m.search.navigating)
	m.keymap.NextMatch.SetEnabled(m.search.navigating)

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.search.input, cmd = m.search.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) helpView() string {
	return m.help.View(m.keymap)
}

func (m *model) processText(msg tea.WindowSizeMsg) {
	m.viewport.Height = msg.Height - lipgloss.Height(m.helpView())
	m.viewport.Width = msg.Width

	// Determine max width of a line.
	m.maxWidth = m.viewport.Width
}

const heightOffset = 2

func (m model) keyHandler(msg tea.KeyMsg) (model, tea.Cmd) {
	km := m.keymap
	var cmd tea.Cmd
	if m.search.visible {
		switch {
		case key.Matches(msg, km.ConfirmSearch):
			if m.search.input.Value() != "" {
				m.viewport.SetHighlights(m.search.Execute(m.viewport.GetContent()))
			} else {
				m.search.Done()
				m.viewport.ClearHighlights()
			}
		case key.Matches(msg, km.CancelSearch):
			m.search.Done()
			m.viewport.ClearHighlights()
		default:
			m.search.input, cmd = m.search.input.Update(msg)
		}
	} else {
		switch {
		case key.Matches(msg, km.Home):
			m.viewport.GotoTop()
		case key.Matches(msg, km.End):
			m.viewport.GotoBottom()
		case key.Matches(msg, km.Search):
			m.search.Show(m.viewport.Width)
			return m, textinput.Blink
		case key.Matches(msg, km.PrevMatch):
			m.viewport.HighlightPrevious()
		case key.Matches(msg, km.NextMatch):
			m.viewport.HighlightNext()
		case key.Matches(msg, km.Quit):
			return m, tea.Quit
		case key.Matches(msg, km.Abort):
			return m, tea.Interrupt
		}
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.search.visible {
		return m.viewport.View() + "\n " + m.search.input.View()
	}

	return m.viewport.View() + "\n" + m.helpView()
}

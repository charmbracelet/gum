// Package pager provides a pager (similar to less) for the terminal.
//
// $ cat file.txt | gum pager
package pager

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
	content             string
	origContent         string
	viewport            viewport.Model
	help                help.Model
	showLineNumbers     bool
	lineNumberStyle     lipgloss.Style
	softWrap            bool
	search              search
	matchStyle          lipgloss.Style
	matchHighlightStyle lipgloss.Style
	maxWidth            int
	keymap              keymap
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.processText(msg)
	case tea.KeyMsg:
		return m.keyHandler(msg)
	}

	m.keymap.PrevMatch.SetEnabled(m.search.query != nil)
	m.keymap.NextMatch.SetEnabled(m.search.query != nil)

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.search.input, cmd = m.search.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) helpView() string {
	return "\n" + m.help.View(m.keymap)
}

func (m *model) processText(msg tea.WindowSizeMsg) {
	m.viewport.Height = msg.Height - lipgloss.Height(m.helpView())
	m.viewport.Width = msg.Width

	var text strings.Builder

	// Determine max width of a line.
	m.maxWidth = m.viewport.Width
	if m.softWrap {
		vpStyle := m.viewport.Style
		m.maxWidth -= vpStyle.GetHorizontalBorderSize() + vpStyle.GetHorizontalMargins() + vpStyle.GetHorizontalPadding()
		if m.showLineNumbers {
			m.maxWidth -= lipgloss.Width("     │ ")
		}
	}

	for i, line := range strings.Split(m.content, "\n") {
		line = strings.ReplaceAll(line, "\t", "    ")
		if m.showLineNumbers {
			text.WriteString(m.lineNumberStyle.Render(fmt.Sprintf("%4d │ ", i+1)))
		}
		if m.softWrap {
			idx := 0
			for ansi.StringWidth(line) >= idx {
				truncatedLine := ansi.Cut(line, idx, m.maxWidth+idx)
				if m.showLineNumbers && idx > 0 {
					text.WriteString(m.lineNumberStyle.Render("     │ "))
				}
				text.WriteString(truncatedLine + "\n")
				idx += m.maxWidth
			}
		} else {
			text.WriteString(line + "\n")
		}
	}

	diffHeight := m.viewport.Height - lipgloss.Height(text.String())
	if diffHeight > 0 && m.showLineNumbers {
		remainingLines := "   ~ │ " + strings.Repeat("\n   ~ │ ", diffHeight-1)
		text.WriteString(m.lineNumberStyle.Render(remainingLines))
	}
	m.viewport.SetContent(text.String())
}

const heightOffset = 2

func (m model) keyHandler(msg tea.KeyMsg) (model, tea.Cmd) {
	km := m.keymap
	var cmd tea.Cmd
	if m.search.active {
		switch {
		case key.Matches(msg, km.ConfirmSearch):
			if m.search.input.Value() != "" {
				m.content = m.origContent
				m.search.Execute(&m)

				// Trigger a view update to highlight the found matches.
				m.search.NextMatch(&m)
				m.processText(tea.WindowSizeMsg{Height: m.viewport.Height + heightOffset, Width: m.viewport.Width})
			} else {
				m.search.Done()
			}
		case key.Matches(msg, km.CancelSearch):
			m.search.Done()
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
			m.search.Begin()
			return m, textinput.Blink
		case key.Matches(msg, km.PrevMatch):
			m.search.PrevMatch(&m)
			m.processText(tea.WindowSizeMsg{Height: m.viewport.Height + heightOffset, Width: m.viewport.Width})
		case key.Matches(msg, km.NextMatch):
			m.search.NextMatch(&m)
			m.processText(tea.WindowSizeMsg{Height: m.viewport.Height + heightOffset, Width: m.viewport.Width})
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
	if m.search.active {
		return m.viewport.View() + "\n " + m.search.input.View()
	}

	return m.viewport.View() + m.helpView()
}

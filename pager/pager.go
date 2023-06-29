// Package pager provides a pager (similar to less) for the terminal.
//
// $ cat file.txt | gum page
package pager

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/gum/timeout"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type model struct {
	content             string
	origContent         string
	viewport            viewport.Model
	helpStyle           lipgloss.Style
	showLineNumbers     bool
	lineNumberStyle     lipgloss.Style
	softWrap            bool
	search              search
	matchStyle          lipgloss.Style
	matchHighlightStyle lipgloss.Style
	maxWidth            int
	timeout             time.Duration
	hasTimeout          bool
}

func (m model) Init() tea.Cmd {
	return timeout.Init(m.timeout, nil)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)

	case tea.WindowSizeMsg:
		m.ProcessText(msg)
	case tea.KeyMsg:
		return m.KeyHandler(msg)
	}

	return m, nil
}

func (m *model) ProcessText(msg tea.WindowSizeMsg) {
	m.viewport.Height = msg.Height - lipgloss.Height(m.helpStyle.Render("?")) - 1
	m.viewport.Width = msg.Width
	textStyle := lipgloss.NewStyle().Width(m.viewport.Width)
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
		for m.softWrap && lipgloss.Width(line) > m.maxWidth {
			truncatedLine := truncate.String(line, uint(m.maxWidth))
			text.WriteString(textStyle.Render(truncatedLine))
			text.WriteString("\n")
			if m.showLineNumbers {
				text.WriteString(m.lineNumberStyle.Render("     │ "))
			}
			line = strings.Replace(line, truncatedLine, "", 1)
		}
		text.WriteString(textStyle.Render(truncate.String(line, uint(m.maxWidth))))
		text.WriteString("\n")
	}

	diffHeight := m.viewport.Height - lipgloss.Height(text.String())
	if diffHeight > 0 && m.showLineNumbers {
		remainingLines := "   ~ │ " + strings.Repeat("\n   ~ │ ", diffHeight-1)
		text.WriteString(m.lineNumberStyle.Render(remainingLines))
	}
	m.viewport.SetContent(text.String())
}

func (m model) KeyHandler(key tea.KeyMsg) (model, func() tea.Msg) {
	var cmd tea.Cmd
	const HeightOffset = 2
	if m.search.active {
		switch key.String() {
		case "enter":
			if m.search.input.Value() != "" {
				m.content = m.origContent
				m.search.Execute(&m)

				// Trigger a view update to highlight the found matches.
				m.search.NextMatch(&m)
				m.ProcessText(tea.WindowSizeMsg{Height: m.viewport.Height + HeightOffset, Width: m.viewport.Width})
			} else {
				m.search.Done()
			}
		case "ctrl+d", "ctrl+c", "esc":
			m.search.Done()
		default:
			m.search.input, cmd = m.search.input.Update(key)
		}
	} else {
		switch key.String() {
		case "g":
			m.viewport.GotoTop()
		case "G":
			m.viewport.GotoBottom()
		case "/":
			m.search.Begin()
		case "p", "N":
			m.search.PrevMatch(&m)
			m.ProcessText(tea.WindowSizeMsg{Height: m.viewport.Height + HeightOffset, Width: m.viewport.Width})
		case "n":
			m.search.NextMatch(&m)
			m.ProcessText(tea.WindowSizeMsg{Height: m.viewport.Height + HeightOffset, Width: m.viewport.Width})
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
		m.viewport, cmd = m.viewport.Update(key)
	}

	return m, cmd
}

func (m model) View() string {
	var timeoutStr string
	if m.hasTimeout {
		timeoutStr = timeout.Str(m.timeout) + " "
	}
	helpMsg := "\n"+timeoutStr+" ↑/↓: Navigate • q: Quit • /: Search "
	if m.search.query != nil {
		helpMsg += "• n: Next Match "
		helpMsg += "• N: Prev Match "
	}
	if m.search.active {
		return m.viewport.View() + "\n"+timeoutStr+ " "+ m.search.input.View()
	}

	return m.viewport.View() + m.helpStyle.Render(helpMsg)
}

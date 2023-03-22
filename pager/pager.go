// Package pager provides a pager (similar to less) for the terminal.
//
// $ cat file.txt | gum page
package pager

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	content         string
	origContent     string
	viewport        viewport.Model
	helpStyle       lipgloss.Style
	showLineNumbers bool
	lineNumberStyle lipgloss.Style
	softWrap        bool
	search          search
	matchStyle      lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

	// Determine max width of a line
	maxLineWidth := m.viewport.Width
	if m.softWrap {
		vpStyle := m.viewport.Style
		maxLineWidth -= vpStyle.GetHorizontalBorderSize() + vpStyle.GetHorizontalMargins() + vpStyle.GetHorizontalPadding()
		if m.showLineNumbers {
			maxLineWidth -= lipgloss.Width("     │ ")
		}
	}

	for i, line := range strings.Split(m.content, "\n") {
		line = strings.ReplaceAll(line, "\t", "    ")
		if m.showLineNumbers {
			text.WriteString(m.lineNumberStyle.Render(fmt.Sprintf("%4d │ ", i+1)))
		}
		for m.softWrap && lipgloss.Width(line) > maxLineWidth {
			truncatedLine := lipglossTruncate(line, maxLineWidth)
			text.WriteString(textStyle.Render(truncatedLine))
			text.WriteString("\n")
			if m.showLineNumbers {
				text.WriteString(m.lineNumberStyle.Render("     │ "))
			}
			line = strings.Replace(line, truncatedLine, "", 1)
		}
		text.WriteString(textStyle.Render(lipglossTruncate(line, maxLineWidth)))
		text.WriteString("\n")
	}

	diffHeight := m.viewport.Height - lipgloss.Height(text.String())
	if diffHeight > 0 && m.showLineNumbers {
		remainingLines := "   ~ │ " + strings.Repeat("\n   ~ │ ", diffHeight-1)
		text.WriteString(m.lineNumberStyle.Render(remainingLines))
	}
	m.viewport.SetContent(text.String())
}

func lipglossTruncate(s string, width int) string {
	var i int
	for i = 0; i < len(s) && lipgloss.Width(s[:i]) < width; i++ {
	} //revive:disable-line:empty-block
	return s[:i]
}

func (m model) KeyHandler(key tea.KeyMsg) (model, func() tea.Msg) {
	var cmd tea.Cmd
	if m.search.active {
		switch key.String() {
		case "enter":
			if m.search.input.Value() != "" {
				m.content = m.origContent
				m.search.Execute(&m)

				// Trigger a view update to highlight the found matches
				m.ProcessText(tea.WindowSizeMsg{Height: m.viewport.Height + 2, Width: m.viewport.Width})
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
		case "p":
			m.search.PrevMatch(&m)
		case "n":
			m.search.NextMatch(&m)
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
		m.viewport, cmd = m.viewport.Update(key)
	}

	return m, cmd
}

func (m model) View() string {
	helpMsg := "\n ↑/↓: Navigate • q: Quit • /: Search • n: Next Match • p: Previous Match"
	if m.search.active {
		return m.viewport.View() + m.helpStyle.Render(helpMsg) + "\n" + m.search.input.View()
	}

	return m.viewport.View() + m.helpStyle.Render(helpMsg)
}

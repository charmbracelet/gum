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
	"github.com/mattn/go-runewidth"
)

type model struct {
	content         string
	viewport        viewport.Model
	helpStyle       lipgloss.Style
	showLineNumbers bool
	lineNumberStyle lipgloss.Style
	softWrap        bool
	search          Search
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
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
				maxLineWidth -= len("     │ ")
			}
		}

		for i, line := range strings.Split(m.content, "\n") {
			line = strings.ReplaceAll(line, "\t", "    ")
			if m.showLineNumbers {
				text.WriteString(m.lineNumberStyle.Render(fmt.Sprintf("%4d │ ", i+1)))
			}
			for m.softWrap && len(line) > maxLineWidth {
				truncatedLine := runewidth.Truncate(line, maxLineWidth, "")
				text.WriteString(textStyle.Render(truncatedLine))
				text.WriteString("\n")
				if m.showLineNumbers {
					text.WriteString(m.lineNumberStyle.Render("     │ "))
				}
				line = strings.Replace(line, truncatedLine, "", 1)
			}
			text.WriteString(textStyle.Render(runewidth.Truncate(line, maxLineWidth, "")))
			text.WriteString("\n")
		}

		diffHeight := m.viewport.Height - lipgloss.Height(text.String())
		if diffHeight > 0 && m.showLineNumbers {
			remainingLines := "   ~ │ " + strings.Repeat("\n   ~ │ ", diffHeight-1)
			text.WriteString(m.lineNumberStyle.Render(remainingLines))
		}
		m.viewport.SetContent(text.String())
	case tea.KeyMsg:
		return m.KeyHandler(msg)
	}

	return m, nil
}

func (m model) KeyHandler(key tea.KeyMsg) (model, func() tea.Msg) {
	var cmd tea.Cmd
	if m.search.Active {
		switch key.String() {
		case "enter":
			if m.search.Input.Value() != "" {
				m.search.Execute(&m)
			} else {
				m.search.Done()
			}
		case "ctrl+d", "ctrl+c", "esc":
			m.search.Done()
		default:
			m.search.Input, cmd = m.search.Input.Update(key)
		}
	} else {
		switch key.String() {
		case "g":
			m.viewport.GotoTop()
		case "G":
			m.viewport.GotoBottom()
		case "/":
			m.search.Begin()
		case "n":
			m.search.NextMatch(&m)
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(key)
	return m, cmd
}

func (m model) View() string {
	helpMsg := "\n ↑/↓: Navigate • q: Quit"
	if m.search.Active {
		return m.viewport.View() + m.helpStyle.Render(helpMsg) + "\n" + m.search.Input.View()
	}

	return m.viewport.View() + m.helpStyle.Render(helpMsg)
}

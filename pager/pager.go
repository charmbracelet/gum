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
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 2
		m.viewport.Width = msg.Width
		textStyle := lipgloss.NewStyle().Width(m.viewport.Width)
		var text strings.Builder
		for i, line := range strings.Split(m.content, "\n") {
			line = strings.ReplaceAll(line, "\t", "    ")
			if m.showLineNumbers {
				text.WriteString(m.lineNumberStyle.Render(fmt.Sprintf("%4d │ ", i+1)))
			}
			text.WriteString(textStyle.Render(runewidth.Truncate(line, m.viewport.Width, "")))
			text.WriteString("\n")
		}
		m.viewport.SetContent(text.String())
	case tea.KeyMsg:
		switch msg.String() {
<<<<<<< HEAD
		case "g":
			m.viewport.GotoTop()
		case "G":
			m.viewport.GotoBottom()
=======
>>>>>>> next
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.viewport.View() + m.helpStyle.Render("\n ↑/↓: Navigate • q: Quit")
}

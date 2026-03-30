// Package progress provides a shell script interface for the progress bar
// bubble. https://github.com/charmbracelet/bubbles/tree/master/progress
//
// It is useful to display progress of a task by reading percentage values
// from standard input.
//
// For example, pipe progress values to display a progress bar:
//
//	for i in $(seq 0 10 100); do echo "$i"; sleep 0.2; done | gum progress --title "Working..."
package progress

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg float64

type doneMsg struct{}

type model struct {
	progress progress.Model
	title    string
	percent  float64
	quitting bool
	done     bool
	padding  []int

	titleStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case tickMsg:
		m.percent = float64(msg)
		if m.percent >= 1.0 {
			m.done = true
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	case doneMsg:
		m.done = true
		m.quitting = true
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - m.padding[1] - m.padding[3]
		return m, nil
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting && m.done {
		return ""
	}

	var s string
	if m.title != "" {
		s = m.titleStyle.Render(m.title) + "\n"
	}
	s += m.progress.ViewAs(m.percent)

	return lipgloss.NewStyle().
		Padding(m.padding...).
		Render(s)
}

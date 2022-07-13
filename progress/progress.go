// Package progress provides a shell script interface for the progress bubble.
// https://github.com/charmbracelet/bubbles/tree/master/progress
//
// It's useful for indicating that something is happening in the background
// that will end after some set time. It can be passed an increment value which
// specifies how much the progress bar should move every interval, which can
// also be configured.
//
//   $ gum progress --increment 0.1 --interval 250ms
//
package progress

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	interval  time.Duration
	progress  progress.Model
	increment float64
}

type tickMsg time.Time

func tickCmd(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd { return tickCmd(m.interval) }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tickMsg:
		if m.progress.Percent() >= 1 {
			return m, tea.Quit
		}
		cmd := m.progress.IncrPercent(m.increment)
		return m, tea.Batch(tickCmd(m.interval), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string { return m.progress.View() }

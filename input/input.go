package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct{ textinput textinput.Model }

func (m model) Init() tea.Cmd { return textinput.Blink }
func (m model) View() string  { return m.textinput.View() }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC, tea.KeyEnter:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

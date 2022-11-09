package choose

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputType int

const (
	Select InputType = iota
	Input
)

type InputModels struct {
	inputState InputStyle
	paginator  paginator.Model
	input      textinput.Model
}

func (m *InputModels) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.inputState {
	case SELECT:
		m.paginator, cmd = m.paginator.Update(msg)
	case INPUT:
		m.input, cmd = m.input.Update(msg)
	}
	return cmd
}

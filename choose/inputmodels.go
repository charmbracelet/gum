package choose

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputStyle int64

// SELECT User is selecting one of the options.
// INPUT User enters into input field.
const (
	SELECT inputStyle = iota
	INPUT
)

type inputModels struct {
	inputState inputStyle
	paginator  paginator.Model
	input      textinput.Model
}

func (m *inputModels) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.inputState {
	case SELECT:
		m.paginator, cmd = m.paginator.Update(msg)
	case INPUT:
		m.input, cmd = m.input.Update(msg)
	}
	return cmd
}

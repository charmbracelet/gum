package input

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() {
	i := textinput.New()
	i.Focus()

	i.SetValue(o.Value)
	i.Prompt = o.Prompt
	i.Placeholder = o.Placeholder
	i.Width = o.Width
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.CursorStyle = o.CursorStyle.ToLipgloss()

	p := tea.NewProgram(model{i}, tea.WithOutput(os.Stderr))
	m, _ := p.StartReturningModel()
	fmt.Println(m.(model).textinput.Value())
}

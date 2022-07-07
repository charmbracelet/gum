package input

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.Placeholder = o.Placeholder
	i.Width = o.Width
	i.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(o.PromptColor))
	i.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(o.CursorColor))

	p := tea.NewProgram(model{i}, tea.WithOutput(os.Stderr))
	m, _ := p.StartReturningModel()
	fmt.Println(m.(model).textinput.Value())
}

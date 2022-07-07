package write

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/textarea
func (o Options) Run() {
	a := textarea.New()
	a.Focus()

	a.Prompt = o.Prompt
	a.Placeholder = o.Placeholder
	a.ShowLineNumbers = o.ShowLineNumbers
	if !o.ShowCursorLine {
		a.FocusedStyle.CursorLine = lipgloss.NewStyle()
		a.BlurredStyle.CursorLine = lipgloss.NewStyle()
	}

	a.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(o.CursorColor))

	a.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(o.PromptColor))
	a.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(o.PromptColor))

	a.SetWidth(o.Width)

	p := tea.NewProgram(model{a}, tea.WithOutput(os.Stderr))
	m, _ := p.StartReturningModel()
	fmt.Println(m.(model).textarea.Value())
}

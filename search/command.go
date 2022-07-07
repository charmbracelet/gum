package search

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/log"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for the search bubble.
// https://github.com/charmbracelet/bubbles/search
func (o Options) Run() {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(o.AccentColor))
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	input, err := stdin.Read()
	if err != nil || input == "" {
		log.Error("No input provided.")
		return
	}
	choices := strings.Split(string(input), "\n")

	p := tea.NewProgram(model{
		textinput: i,
		choices:   choices,
		matches:   matchAll(choices),
		indicator: o.Indicator,
	}, tea.WithOutput(os.Stderr))

	tm, _ := p.StartReturningModel()
	m := tm.(model)

	if len(m.matches) > m.selected && m.selected >= 0 {
		fmt.Println(m.matches[m.selected].Str)
	}
}

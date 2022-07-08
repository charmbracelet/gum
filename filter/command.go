package filter

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/log"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(o.PromptColor))
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	input, err := stdin.Read()
	if err != nil {
		log.Error("Could not read stdin.")
		return
	}

	var choices []string
	if input != "" {
		choices = strings.Split(string(input), "\n")
	} else {
		choices = files.List()
	}

	p := tea.NewProgram(model{
		textinput:      i,
		choices:        choices,
		matches:        matchAll(choices),
		indicator:      o.Indicator,
		highlightStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(o.HighlightColor)),
		indicatorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(o.IndicatorColor)),
	}, tea.WithOutput(os.Stderr))

	tm, _ := p.StartReturningModel()
	m := tm.(model)

	if len(m.matches) > m.selected && m.selected >= 0 {
		fmt.Println(m.matches[m.selected].Str)
	}
}

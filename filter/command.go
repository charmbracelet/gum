package filter

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() error {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	v := viewport.New(o.Width, o.Height)

	var choices []string
	if input, _ := stdin.Read(); input != "" {
		choices = strings.Split(strings.TrimSpace(input), "\n")
	} else {
		choices = files.List()
	}

	options := []tea.ProgramOption{tea.WithOutput(os.Stderr)}
	if o.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	p := tea.NewProgram(model{
		choices:        choices,
		indicator:      o.Indicator,
		matches:        matchAll(choices),
		textinput:      i,
		viewport:       &v,
		indicatorStyle: o.IndicatorStyle.ToLipgloss(),
		matchStyle:     o.MatchStyle.ToLipgloss(),
		textStyle:      o.TextStyle.ToLipgloss(),
		height:         o.Height,
	}, options...)

	tm, err := p.StartReturningModel()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}
	m := tm.(model)

	if m.aborted {
		return exit.ErrAborted
	}
	if len(m.matches) > m.selected && m.selected >= 0 {
		fmt.Println(m.matches[m.selected].Str)
	}

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

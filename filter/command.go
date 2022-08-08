package filter

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"

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
		input = strings.TrimSpace(input)
		if input != "" {
			choices = strings.Split(input, "\n")
		}
	} else {
		choices = files.List()
	}

	if len(choices) == 0 {
		return errors.New("no options provided, see `gum filter --help`")
	}

	options := []tea.ProgramOption{tea.WithOutput(os.Stderr)}
	if o.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	var matches []fuzzy.Match
	if o.Value != "" {
		i.SetValue(o.Value)
		matches = fuzzy.Find(o.Value, choices)
	} else {
		matches = matchAll(choices)
	}

	p := tea.NewProgram(model{
		choices:        choices,
		indicator:      o.Indicator,
		matches:        matches,
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

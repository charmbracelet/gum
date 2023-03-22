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
		input = strings.TrimSuffix(input, "\n")
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
	}
	switch {
	case o.Value != "" && o.Fuzzy:
		matches = fuzzy.Find(o.Value, choices)
	case o.Value != "" && !o.Fuzzy:
		matches = exactMatches(o.Value, choices)
	default:
		matches = matchAll(choices)
	}

	if o.NoLimit {
		o.Limit = len(choices)
	}

	p := tea.NewProgram(model{
		choices:               choices,
		indicator:             o.Indicator,
		matches:               matches,
		header:                o.Header,
		textinput:             i,
		viewport:              &v,
		indicatorStyle:        o.IndicatorStyle.ToLipgloss(),
		selectedPrefixStyle:   o.SelectedPrefixStyle.ToLipgloss(),
		selectedPrefix:        o.SelectedPrefix,
		unselectedPrefixStyle: o.UnselectedPrefixStyle.ToLipgloss(),
		unselectedPrefix:      o.UnselectedPrefix,
		matchStyle:            o.MatchStyle.ToLipgloss(),
		headerStyle:           o.HeaderStyle.ToLipgloss(),
		textStyle:             o.TextStyle.ToLipgloss(),
		height:                o.Height,
		selected:              make(map[string]struct{}),
		limit:                 o.Limit,
		reverse:               o.Reverse,
		fuzzy:                 o.Fuzzy,
	}, options...)

	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}
	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}

	// allSelections contains values only if limit is greater
	// than 1 or if flag --no-limit is passed, hence there is
	// no need to further checks
	if len(m.selected) > 0 {
		for k := range m.selected {
			fmt.Println(k)
		}
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		fmt.Println(m.matches[m.cursor].Str)
	}

	if !o.Strict && len(m.textinput.Value()) != 0 && len(m.matches) == 0 {
		fmt.Println(m.textinput.Value())
	}
	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

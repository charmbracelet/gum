package filter

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/sahilm/fuzzy"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() error {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.PlaceholderStyle = o.PlaceholderStyle.ToLipgloss()
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	v := viewport.New(o.Width, o.Height)

	if len(o.Options) == 0 {
		if input, _ := stdin.Read(); input != "" {
			o.Options = strings.Split(input, "\n")
		} else {
			o.Options = files.List()
		}
	}

	if len(o.Options) == 0 {
		return errors.New("no options provided, see `gum filter --help`")
	}

	if o.SelectIfOne && len(o.Options) == 1 {
		fmt.Println(o.Options[0])
		return nil
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
		matches = fuzzy.Find(o.Value, o.Options)
	case o.Value != "" && !o.Fuzzy:
		matches = exactMatches(o.Value, o.Options)
	default:
		matches = matchAll(o.Options)
	}

	if o.NoLimit {
		o.Limit = len(o.Options)
	}

	p := tea.NewProgram(model{
		choices:               o.Options,
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
		cursorTextStyle:       o.CursorTextStyle.ToLipgloss(),
		height:                o.Height,
		selected:              make(map[string]struct{}),
		limit:                 o.Limit,
		reverse:               o.Reverse,
		fuzzy:                 o.Fuzzy,
		timeout:               o.Timeout,
		hasTimeout:            o.Timeout > 0,
		sort:                  o.Sort && o.FuzzySort,
		strict:                o.Strict,
		showHelp:              o.ShowHelp,
		keymap:                defaultKeymap(),
		help:                  help.New(),
	}, options...)

	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}
	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}
	if m.timedOut {
		return exit.ErrTimeout
	}

	isTTY := term.IsTerminal(os.Stdout.Fd())

	// allSelections contains values only if limit is greater
	// than 1 or if flag --no-limit is passed, hence there is
	// no need to further checks
	if len(m.selected) > 0 {
		o.checkSelected(m, isTTY)
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		if isTTY {
			fmt.Println(m.matches[m.cursor].Str)
		} else {
			fmt.Println(ansi.Strip(m.matches[m.cursor].Str))
		}
	}

	return nil
}

func (o Options) checkSelected(m model, isTTY bool) {
	for k := range m.selected {
		if isTTY {
			fmt.Println(k)
		} else {
			fmt.Println(ansi.Strip(k))
		}
	}
}

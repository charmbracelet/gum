package filter

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/internal/tty"
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/x/ansi"
	"github.com/sahilm/fuzzy"
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
		if input, _ := stdin.Read(stdin.StripANSI(o.StripANSI)); input != "" {
			o.Options = strings.Split(input, o.InputDelimiter)
		} else {
			o.Options = files.List()
		}
	}

	if len(o.Options) == 0 {
		return errors.New("no options provided, see `gum filter --help`")
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	options := []tea.ProgramOption{
		tea.WithOutput(os.Stderr),
		tea.WithReportFocus(),
		tea.WithContext(ctx),
	}
	if o.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	var matches []fuzzy.Match
	if o.Value != "" {
		i.SetValue(o.Value)
	}

	choices := map[string]string{}
	filteringChoices := []string{}
	for _, opt := range o.Options {
		s := ansi.Strip(opt)
		choices[s] = opt
		filteringChoices = append(filteringChoices, s)
	}
	switch {
	case o.Value != "" && o.Fuzzy:
		matches = fuzzy.Find(o.Value, filteringChoices)
	case o.Value != "" && !o.Fuzzy:
		matches = exactMatches(o.Value, filteringChoices)
	default:
		matches = matchAll(filteringChoices)
	}

	if o.NoLimit {
		o.Limit = len(o.Options)
	}

	if o.SelectIfOne && len(matches) == 1 {
		tty.Println(matches[0].Str)
		return nil
	}

	km := defaultKeymap()
	if o.NoLimit || o.Limit > 1 {
		km.Toggle.SetEnabled(true)
		km.ToggleAndPrevious.SetEnabled(true)
		km.ToggleAndNext.SetEnabled(true)
		km.ToggleAll.SetEnabled(true)
	}
	top, right, bottom, left := style.ParsePadding(o.Padding)
	m := model{
		choices:               choices,
		filteringChoices:      filteringChoices,
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
		padding:               []int{top, right, bottom, left},
		selected:              make(map[string]struct{}),
		limit:                 o.Limit,
		reverse:               o.Reverse,
		fuzzy:                 o.Fuzzy,
		sort:                  o.Sort && o.FuzzySort,
		strict:                o.Strict,
		showHelp:              o.ShowHelp,
		keymap:                km,
		help:                  help.New(),
	}

	isSelectAll := len(o.Selected) == 1 && o.Selected[0] == "*"
	currentSelected := 0
	if len(o.Selected) > 0 {
		for i, option := range matches {
			if currentSelected >= o.Limit || (!isSelectAll && !slices.Contains(o.Selected, option.Str)) {
				continue
			}
			if o.Limit == 1 {
				m.cursor = i
				m.selected[option.Str] = struct{}{}
			} else {
				currentSelected++
				m.selected[option.Str] = struct{}{}
			}
		}
	}

	tm, err := tea.NewProgram(m, options...).Run()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}

	m = tm.(model)
	if !m.submitted {
		return errors.New("nothing selected")
	}

	// allSelections contains values only if limit is greater
	// than 1 or if flag --no-limit is passed, hence there is
	// no need to further checks
	if len(m.selected) > 0 {
		o.checkSelected(m)
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		tty.Println(m.matches[m.cursor].Str)
	}

	return nil
}

func (o Options) checkSelected(m model) {
	out := []string{}
	for k := range m.selected {
		out = append(out, k)
	}
	tty.Println(strings.Join(out, o.OutputDelimiter))
}

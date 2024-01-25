package choose

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	if len(o.Options) <= 0 {
		input, _ := stdin.Read()
		if input == "" {
			return errors.New("no options provided, see `gum choose --help`")
		}
		o.Options = strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	}

	theme := huh.ThemeCharm()
	options := huh.NewOptions(o.Options...)

	theme.Focused.Base.Border(lipgloss.Border{})
	theme.Focused.SelectSelector = o.CursorStyle.ToLipgloss().SetString(o.Cursor)
	theme.Focused.SelectedOption = o.SelectedItemStyle.ToLipgloss()
	theme.Focused.UnselectedOption = o.ItemStyle.ToLipgloss()
	theme.Focused.SelectedPrefix = o.SelectedItemStyle.ToLipgloss().SetString(o.SelectedPrefix)
	theme.Focused.UnselectedPrefix = o.ItemStyle.ToLipgloss().SetString(o.UnselectedPrefix)

	for _, s := range o.Selected {
		for i, opt := range options {
			if s == opt.Key || s == opt.Value {
				options[i] = opt.Selected(true)
			}
		}
	}

	if o.NoLimit {
		o.Limit = len(o.Options)
	}

	if o.Limit > 1 {
		var choices []string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Options(options...).
					Title(o.Header).
					Filterable(false).
					Height(o.Height).
					Limit(o.Limit).
					Value(&choices),
			),
		).
			WithTheme(theme).
			WithShowHelp(false).
			Run()
		if err != nil {
			return err
		}
		if len(choices) > 0 {
			fmt.Println(strings.Join(choices, "\n"))
		}
		return nil
	}

	var choice string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(options...).
				Title(o.Header).
				Height(o.Height).
				Value(&choice),
		),
	).
		WithTheme(theme).
		WithShowHelp(false).
		Run()
	if err != nil {
		return err
	}
	fmt.Println(choice)
	return nil
}

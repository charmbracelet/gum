package choose

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"

	"github.com/charmbracelet/gum/ansi"
	"github.com/charmbracelet/gum/internal/stdin"
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

	if o.SelectIfOne && len(o.Options) == 1 {
		fmt.Println(o.Options[0])
		return nil
	}

	theme := huh.ThemeCharm()
	options := huh.NewOptions(o.Options...)

	theme.Focused.Base.Border(lipgloss.Border{})
	theme.Focused.Title = o.HeaderStyle.ToLipgloss()
	theme.Focused.SelectSelector = o.CursorStyle.ToLipgloss().SetString(o.Cursor)
	theme.Focused.MultiSelectSelector = o.CursorStyle.ToLipgloss().SetString(o.Cursor)
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

	width := max(widest(o.Options)+
		max(lipgloss.Width(o.SelectedPrefix), lipgloss.Width(o.UnselectedPrefix))+
		lipgloss.Width(o.Cursor)+1, lipgloss.Width(o.Header)+1)

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
			WithWidth(width).
			WithShowHelp(false).
			WithTheme(theme).
			Run()
		if err != nil {
			return err
		}
		if len(choices) > 0 {
			s := strings.Join(choices, "\n")
			if isatty.IsTerminal(os.Stdout.Fd()) {
				fmt.Println(s)
			} else {
				fmt.Print(ansi.Strip(s))
			}
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
		WithWidth(width).
		WithTheme(theme).
		WithShowHelp(false).
		Run()

	if err != nil {
		return err
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println(choice)
	} else {
		fmt.Print(ansi.Strip(choice))
	}

	return nil
}

func widest(options []string) int {
	var max int
	for _, o := range options {
		w := lipgloss.Width(o)
		if w > max {
			max = w
		}
	}
	return max
}

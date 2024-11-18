package choose

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
)

const widthBuffer = 2

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	if len(o.Options) <= 0 {
		input, _ := stdin.Read()
		if input == "" {
			return errors.New("no options provided, see `gum choose --help`")
		}
		o.Options = strings.Split(input, "\n")
	}

	if o.SelectIfOne && len(o.Options) == 1 {
		fmt.Println(o.Options[0])
		return nil
	}

	theme := huh.ThemeCharm()
	options := huh.NewOptions(o.Options...)

	theme.Focused.Base = lipgloss.NewStyle()
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

	width := max(widest(o.Options)+
		max(lipgloss.Width(o.SelectedPrefix)+lipgloss.Width(o.UnselectedPrefix))+
		lipgloss.Width(o.Cursor)+1, lipgloss.Width(o.Header)+widthBuffer)

	if o.NoLimit {
		o.Limit = 0
	}

	if o.Limit > 1 || o.NoLimit {
		var choices []string

		field := huh.NewMultiSelect[string]().
			Options(options...).
			Title(o.Header).
			Height(o.Height).
			Limit(o.Limit).
			Value(&choices)

		form := huh.NewForm(huh.NewGroup(field))

		err := form.
			WithWidth(width).
			WithShowHelp(o.ShowHelp).
			WithTheme(theme).
			WithTimeout(o.Timeout).
			Run()
		if err != nil {
			return exit.Handle(err, o.Timeout)
		}
		if len(choices) > 0 {
			s := strings.Join(choices, "\n")
			ansiprint(s)
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
		WithTimeout(o.Timeout).
		WithShowHelp(o.ShowHelp).
		Run()
	if err != nil {
		return exit.Handle(err, o.Timeout)
	}

	if term.IsTerminal(os.Stdout.Fd()) {
		fmt.Println(choice)
	} else {
		fmt.Print(ansi.Strip(choice))
	}

	return nil
}

func widest(options []string) int {
	var maxw int
	for _, o := range options {
		w := lipgloss.Width(o)
		if w > maxw {
			maxw = w
		}
	}
	return maxw
}

func ansiprint(s string) {
	if term.IsTerminal(os.Stdout.Fd()) {
		fmt.Println(s)
	} else {
		fmt.Print(ansi.Strip(s))
	}
}

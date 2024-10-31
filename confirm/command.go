package confirm

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	theme := huh.ThemeCharm()
	theme.Focused.Title = o.PromptStyle.ToLipgloss()
	theme.Focused.FocusedButton = o.SelectedStyle.ToLipgloss()
	theme.Focused.BlurredButton = o.UnselectedStyle.ToLipgloss()

	choice := o.Default

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Affirmative(o.Affirmative).
				Negative(o.Negative).
				Title(o.Prompt).
				Value(&choice),
		),
	).
		WithTimeout(o.Timeout).
		WithTheme(theme).
		WithShowHelp(o.ShowHelp).
		Run()

	if err != nil {
		allowErr := o.errIsValidTimeout(err)

		if !allowErr {
			return fmt.Errorf("unable to run confirm: %w", err)
		}
	}

	if !choice {
		os.Exit(1)
	}

	return nil
}

// errIsValidTimeout returns false unless 1) the user has specified a nonzero timeout and 2) the error is a huh.ErrTimeout.
func (o Options) errIsValidTimeout(err error) bool {
	errWasTimeout := err.Error() == huh.ErrTimeout.Error()
	timeoutsExpected := o.Timeout > 0

	return errWasTimeout && timeoutsExpected
}

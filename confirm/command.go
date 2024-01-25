package confirm

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	value := o.Default

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle().Padding(1, 0)
	theme.Focused.Title = o.PromptStyle.ToLipgloss()
	theme.Focused.BlurredButton = o.UnselectedStyle.ToLipgloss()
	theme.Focused.FocusedButton = o.SelectedStyle.ToLipgloss()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title(o.Prompt).
				Affirmative(o.Affirmative).
				Negative(o.Negative).
				Value(&value),
		),
	).
		WithTheme(theme).
		WithShowHelp(false).
		Run()

	if err != nil {
		return fmt.Errorf("unable to run confirm: %w", err)
	}

	if value {
		os.Exit(0)
	}
	os.Exit(1)
	return nil
}

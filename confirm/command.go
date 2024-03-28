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
	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle().Margin(0, 1)
	theme.Focused.Title = o.PromptStyle.ToLipgloss()
	theme.Focused.FocusedButton = o.SelectedStyle.ToLipgloss()
	theme.Focused.BlurredButton = o.UnselectedStyle.ToLipgloss()

	var choice bool = o.Default

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Affirmative(o.Affirmative).
				Negative(o.Negative).
				Title(o.Prompt).
				Value(&choice),
		),
	).
		WithTheme(theme).
		WithShowHelp(false).
		Run()

	if err != nil {
		return fmt.Errorf("unable to run confirm: %w", err)
	}

	if !choice {
		os.Exit(1)
	}

	return nil
}

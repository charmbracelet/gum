package confirm

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	m, err := tea.NewProgram(model{
		affirmative:     o.Affirmative,
		negative:        o.Negative,
		confirmation:    o.Default,
		timeout:         o.Timeout,
		hasTimeout:      o.Timeout > 0,
		prompt:          o.Prompt,
		selectedStyle:   o.SelectedStyle.ToLipgloss(),
		unselectedStyle: o.UnselectedStyle.ToLipgloss(),
		promptStyle:     o.PromptStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr)).Run()

	if err != nil {
		return fmt.Errorf("unable to run confirm: %w", err)
	}

	if m.(model).aborted {
		os.Exit(130)
	} else if m.(model).confirmation {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

package confirm

import (
	"os"

	"github.com/charmbracelet/gum/internal/exit"

	tea "github.com/charmbracelet/bubbletea"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	m, err := tea.NewProgram(model{
		affirmative:      o.Affirmative,
		negative:         o.Negative,
		confirmation:     o.Default,
		defaultSelection: o.Default,
		timeout:          o.Timeout,
		hasTimeout:       o.Timeout > 0,
		prompt:           o.Prompt,
		selectedStyle:    o.SelectedStyle.ToLipgloss(),
		unselectedStyle:  o.UnselectedStyle.ToLipgloss(),
		promptStyle:      o.PromptStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return exit.Handle(err, o.Timeout)
	}

	if m.(model).aborted {
		os.Exit(exit.StatusAborted)
	} else if m.(model).confirmation {
		os.Exit(0)
	}

	os.Exit(1)

	return nil
}

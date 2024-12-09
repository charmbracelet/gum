package confirm

import (
	"errors"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/gum/internal/exit"

	tea "github.com/charmbracelet/bubbletea"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	tm, err := tea.NewProgram(model{
		affirmative:      o.Affirmative,
		negative:         o.Negative,
		confirmation:     o.Default,
		defaultSelection: o.Default,
		timeout:          o.Timeout,
		hasTimeout:       o.Timeout > 0,
		keys:             defaultKeymap(o.Affirmative, o.Negative),
		help:             help.New(),
		showHelp:         o.ShowHelp,
		prompt:           o.Prompt,
		selectedStyle:    o.SelectedStyle.ToLipgloss(),
		unselectedStyle:  o.UnselectedStyle.ToLipgloss(),
		promptStyle:      o.PromptStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return err
	}

	m := tm.(model)
	if m.timedOut {
		return exit.ErrTimeout
	}
	if m.aborted {
		return exit.ErrAborted
	}
	if m.confirmation {
		return nil
	}

	return errors.New("not confirmed")
}

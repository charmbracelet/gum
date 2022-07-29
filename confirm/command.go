package confirm

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	m, err := tea.NewProgram(model{
		affirmative:     o.Affirmative,
		negative:        o.Negative,
		confirmation:    o.Default,
		prompt:          o.Prompt,
		selectedStyle:   o.SelectedStyle.ToLipgloss(),
		unselectedStyle: o.UnselectedStyle.ToLipgloss(),
		promptStyle:     o.PromptStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr)).StartReturningModel()

	if err != nil {
		return err
	}

	if m.(model).confirmation {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

	return nil
}

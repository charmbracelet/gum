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
		vertical:        o.Vertical,
		selected:        0,
		prompt:          o.Prompt,
		selectedStyle:   o.SelectedStyle.ToLipgloss(),
		unselectedStyle: o.UnselectedStyle.ToLipgloss(),
		promptStyle:     o.PromptStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr)).StartReturningModel()

	if err != nil {
		return err
	}

	os.Exit(m.(model).selected)

	return nil
}

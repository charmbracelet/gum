package confirm

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/gum/style"
)

// CtrlC Default Return code in case of ctrl-c by user.
const CtrlC = 130

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	var (
		options        []string
		selectedOption selectionType
	)
	// set options
	options = append(options, o.Affirmative)
	if o.Negative != "" {
		options = append(options, o.Negative)
	}
	if o.Canceled != "" && o.WithCancel {
		options = append(options, o.Canceled)
	}

	// what is default
	if !o.Default && o.Negative != "" {
		selectedOption = Negative
	} else {
		selectedOption = Confirmed
	}

	m, err := tea.NewProgram(model{
		options:         options,
		selectedOption:  selectedOption,
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

	switch m.(model).selectedOption {
	case Confirmed:
		os.Exit(0)
	case Negative:
		os.Exit(1)
	case Cancel:

		os.Exit(CtrlC)
	}

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

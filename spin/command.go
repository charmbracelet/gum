package spin

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"

	"github.com/charmbracelet/gum/internal/exit"
)

// Run provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/spinner
func (o Options) Run() error {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())

	s := spinner.New()
	s.Style = o.SpinnerStyle.ToLipgloss()
	s.Spinner = spinnerMap[o.Spinner]
	m := model{
		spinner:    s,
		title:      o.TitleStyle.ToLipgloss().Render(o.Title),
		command:    o.Command,
		align:      o.Align,
		timeout:    o.Timeout,
		hasTimeout: o.Timeout > 0,
	}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	mm, err := p.Run()
	m = mm.(model)

	if err != nil {
		return fmt.Errorf("failed to run spin: %w", err)
	}

	if m.aborted {
		return exit.ErrAborted
	}

	if err != nil {
		return fmt.Errorf("failed to access stdout: %w", err)
	}

	if o.ShowOutput {
		if isTTY {
			_, err := os.Stdout.WriteString(m.stdout)
			if err != nil {
				return fmt.Errorf("failed to write to stdout: %w", err)
			}
		}
	}

	os.Exit(m.status)
	return nil
}

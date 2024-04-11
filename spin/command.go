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
		showOutput: o.ShowOutput && isTTY,
		showError:  o.ShowError,
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

	// If the command succeeds, and we are printing output and we are in a TTY then push the STDOUT we got to the actual
	// STDOUT for piping or other things.
	if m.status == 0 {
		if o.ShowOutput {
			// BubbleTea writes the View() to stderr.
			// If the program is being piped then put the accumulated output in stdout.
			if !isTTY {
				_, err := os.Stdout.WriteString(m.stdout)
				if err != nil {
					return fmt.Errorf("failed to write to stdout: %w", err)
				}
			}
		}
	} else if o.ShowError {
		// Otherwise if we are showing errors and the command did not exit with a 0 status code then push all of the command
		// output to the terminal. This way failed commands can be debugged.
		_, err := os.Stdout.WriteString(m.output)
		if err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
	}

	os.Exit(m.status)
	return nil
}

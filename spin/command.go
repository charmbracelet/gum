package spin

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/style"
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

	// If the command succeeds, and we are printing output and we are in a TTY then push the STDOUT we got to the actual
	// STDOUT for piping or other things.
	if m.status == 0 {
		if o.ShowOutput {
			if isTTY {
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

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

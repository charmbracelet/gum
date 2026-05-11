package spin

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/x/term"
)

// Run provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/spinner
func (o Options) Run() error {
	isOutTTY := term.IsTerminal(os.Stdout.Fd())
	isErrTTY := term.IsTerminal(os.Stderr.Fd())

	s := spinner.New()
	s.Style = o.SpinnerStyle.ToLipgloss()
	s.Spinner = spinnerMap[o.Spinner]
	top, right, bottom, left := style.ParsePadding(o.Padding)
	m := model{
		spinner:    s,
		title:      o.TitleStyle.ToLipgloss().Render(o.Title),
		command:    o.Command,
		align:      o.Align,
		showStdout: (o.ShowOutput || o.ShowStdout) && isOutTTY,
		showStderr: (o.ShowOutput || o.ShowStderr) && isErrTTY,
		showError:  o.ShowError,
		isTTY:      isErrTTY,
		padding:    []int{top, right, bottom, left},
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	tm, err := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
		tea.WithInput(nil),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to run action: %w", err)
	}

	m = tm.(model)
	// If the command succeeds, and we are printing output and we are in a TTY then push the STDOUT we got to the actual
	// STDOUT for piping or other things.
	//nolint:nestif
	if m.err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "%s\n", m.err.Error()); err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
		return exit.ErrExit(1)
	} else if m.status == 0 {
		var output string
		if o.ShowOutput || (o.ShowStdout && o.ShowStderr) {
			output = m.output
		} else if o.ShowStdout {
			output = m.stdout
		} else if o.ShowStderr {
			output = m.stderr
		}
		if output != "" {
			if _, err := os.Stdout.WriteString(output); err != nil {
				return fmt.Errorf("failed to write to stdout: %w", err)
			}
		}
	} else if o.ShowError {
		// Otherwise if we are showing errors and the command did not exit with a 0 status code then push all of the command
		// output to the terminal. This way failed commands can be debugged.
		if _, err := os.Stdout.WriteString(m.output); err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
	}

	return exit.ErrExit(m.status)
}

package spin

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/spinner
func (o Options) Run() error {
	s := spinner.New()
	s.Style = o.SpinnerStyle.ToLipgloss()
	s.Spinner = spinnerMap[o.Spinner]
	m := model{
		spinner: s,
		title:   o.TitleStyle.ToLipgloss().Render(o.Title),
		command: o.Command,
		align:   o.Align,
	}
	if o.TitleFromStdout {
		m.readFromStdin = true
		m.stdin = make(chan string)
	}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	mm, err := p.StartReturningModel()
	m = mm.(model)

	if err != nil {
		return fmt.Errorf("failed to run spin: %w", err)
	}

	if o.ShowOutput {
		fmt.Fprint(os.Stdout, m.stdout)
		fmt.Fprint(os.Stderr, m.stderr)
	}

	if m.aborted {
		return exit.ErrAborted
	}

	os.Exit(m.status)
	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

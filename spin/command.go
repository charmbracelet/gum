package spin

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	}
	p := tea.NewProgram(m)
	mm, err := p.StartReturningModel()
	m = mm.(model)

	if err != nil {
		return err
	}

	if o.ShowOutput {
		fmt.Print(m.output)
	}

	os.Exit(m.status)

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	return style.HideFlags(ctx)
}

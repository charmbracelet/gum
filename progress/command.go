package progress

import (
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Run runs the progress command.
func (o Options) Run() error {
	p := progress.New(
		progress.WithGradient(o.ColorStart, o.ColorEnd),
		progress.WithSpringOptions(o.Frequency, o.Damping),
	)
	m := model{
		progress:  p,
		interval:  o.Interval,
		increment: o.Increment,
	}
	return tea.NewProgram(m, tea.WithOutput(os.Stderr)).Start()
}

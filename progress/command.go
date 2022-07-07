package progress

import (
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Run runs the progress command.
func (o Options) Run() {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithSpringOptions(o.Frequency, o.Damping),
	)
	m := model{
		progress:  p,
		interval:  o.Interval,
		increment: o.Increment,
	}
	_ = tea.NewProgram(m, tea.WithOutput(os.Stderr)).Start()
}

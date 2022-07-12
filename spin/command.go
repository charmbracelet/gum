package spin

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Run provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/spinner
func (o Options) Run() {
	s := spinner.New()
	s.Style = o.SpinnerStyle.ToLipgloss()
	s.Spinner = spinnerMap[o.Spinner]
	m := model{
		spinner: s,
		title:   o.TitleStyle.ToLipgloss().Render(o.Title),
		command: o.Command,
	}
	p := tea.NewProgram(m)
	_ = p.Start()
}

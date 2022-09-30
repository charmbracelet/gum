package pager

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for the viewport bubble.
// https://github.com/charmbracelet/bubbles/viewport
func (o Options) Run() error {
	vp := viewport.New(o.Style.Width, o.Style.Height)
	vp.Style = o.Style.ToLipgloss()
	var err error

	if o.Content == "" {
		o.Content, err = stdin.Read()
		if err != nil {
			return err
		}
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return err
	}
	md, err := renderer.Render(o.Content)
	vp.SetContent(md)

	model := model{
		viewport:  vp,
		helpStyle: o.HelpStyle.ToLipgloss(),
	}
	if err != nil {
		return err
	}

	return tea.NewProgram(model, tea.WithAltScreen()).Start()
}

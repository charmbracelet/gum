package pager

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for the viewport bubble.
// https://github.com/charmbracelet/bubbles/viewport
func (o Options) Run() error {
	vp := viewport.New(o.Style.Width, o.Style.Height)
	vp.Style = o.Style.ToLipgloss()

	if o.Content == "" {
		stdin, err := stdin.Read()
		if err != nil {
			return fmt.Errorf("unable to read stdin")
		}
		if stdin != "" {
			o.Content = stdin
		} else {
			return fmt.Errorf("provide some content to display")
		}
	}

	model := model{
		viewport:        vp,
		helpStyle:       o.HelpStyle.ToLipgloss(),
		content:         o.Content,
		showLineNumbers: o.ShowLineNumbers,
		lineNumberStyle: o.LineNumberStyle.ToLipgloss(),
	}
	err := tea.NewProgram(model, tea.WithAltScreen()).Start()
	if err != nil {
		return fmt.Errorf("unable to start program: %w", err)
	}
	return nil
}

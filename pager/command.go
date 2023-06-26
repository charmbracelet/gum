package pager

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
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
			// Sanitize the input from stdin by removing backspace sequences.
			backspace := regexp.MustCompile(".\x08")
			o.Content = backspace.ReplaceAllString(stdin, "")
		} else {
			return fmt.Errorf("provide some content to display")
		}
	}

	model := model{
		viewport:        vp,
		helpStyle:       o.HelpStyle.ToLipgloss(),
		content:         o.Content,
		origContent:     o.Content,
		showLineNumbers: o.ShowLineNumbers,
		lineNumberStyle: o.LineNumberStyle.ToLipgloss(),
		softWrap:        o.SoftWrap,
		matchStyle:          o.MatchStyle.ToLipgloss(),
		matchHighlightStyle: o.MatchHighlightStyle.ToLipgloss(),
		timeout:         o.Timeout,
		hasTimeout:      o.Timeout > 0,
	}
	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return fmt.Errorf("unable to start program: %w", err)
	}
	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

package pager

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
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
			// Sanitize the input from stdin by removing backspace sequences.
			backspace := regexp.MustCompile(".\x08")
			o.Content = backspace.ReplaceAllString(stdin, "")
		} else {
			return fmt.Errorf("provide some content to display")
		}
	}

	ctx := context.Background()
	if o.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, o.Timeout)
		defer cancel()
	}

	_, err := tea.NewProgram(
		model{
			viewport:            vp,
			helpStyle:           o.HelpStyle.ToLipgloss(),
			content:             o.Content,
			origContent:         o.Content,
			showLineNumbers:     o.ShowLineNumbers,
			lineNumberStyle:     o.LineNumberStyle.ToLipgloss(),
			softWrap:            o.SoftWrap,
			matchStyle:          o.MatchStyle.ToLipgloss(),
			matchHighlightStyle: o.MatchHighlightStyle.ToLipgloss(),
		},
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return exit.ErrTimeout
		}
		if errors.Is(err, tea.ErrInterrupted) {
			return exit.ErrAborted
		}
		return fmt.Errorf("unable to start program: %w", err)
	}
	return nil
}

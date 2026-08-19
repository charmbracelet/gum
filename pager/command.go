package pager

import (
	"fmt"
	"regexp"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/gum/v2/internal/stdin"
	"charm.land/gum/v2/internal/timeout"
)

// Run provides a shell script interface for the viewport bubble.
// https://github.com/charmbracelet/bubbles/viewport
func (o Options) Run() error {
	vp := viewport.New(viewport.WithWidth(o.Style.Width), viewport.WithHeight(o.Style.Height))
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

	m := model{
		viewport:            vp,
		help:                help.New(),
		content:             o.Content,
		origContent:         o.Content,
		showLineNumbers:     o.ShowLineNumbers,
		lineNumberStyle:     o.LineNumberStyle.ToLipgloss(),
		softWrap:            o.SoftWrap,
		matchStyle:          o.MatchStyle.ToLipgloss(),
		matchHighlightStyle: o.MatchHighlightStyle.ToLipgloss(),
		keymap:              defaultKeymap(),
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	_, err := tea.NewProgram(
		m,
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to start program: %w", err)
	}

	return nil
}

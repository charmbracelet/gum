package pager

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
)

// Run provides a shell script interface for the viewport bubble.
// https://github.com/charmbracelet/bubbles/viewport
func (o Options) Run() error {
	vp := viewport.New(o.Style.Width, o.Style.Height)
	vp.Style = o.Style.ToLipgloss()

	if o.Content == "" {
		stdin, err := stdin.Read(stdin.StripANSI(true))
		if err != nil {
			return fmt.Errorf("unable to read stdin")
		}
		if stdin != "" {
			o.Content = stdin
		} else {
			return fmt.Errorf("provide some content to display")
		}
	}

	if o.ShowLineNumbers {
		vp.LeftGutterFunc = viewport.LineNumberGutter(o.LineNumberStyle.ToLipgloss())
	}

	vp.SoftWrap = o.SoftWrap
	vp.FillHeight = o.ShowLineNumbers
	vp.SetContent(o.Content)
	vp.HighlightStyle = o.MatchStyle.ToLipgloss()
	vp.SelectedHighlightStyle = o.MatchHighlightStyle.ToLipgloss()

	m := model{
		viewport:        vp,
		help:            help.New(),
		showLineNumbers: o.ShowLineNumbers,
		lineNumberStyle: o.LineNumberStyle.ToLipgloss(),
		keymap:          defaultKeymap(),
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	w, _ := os.OpenFile("pager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	defer w.Close()
	log.SetOutput(w)
	_, err := tea.NewProgram(
		m,
		// tea.WithAltScreen(),
		tea.WithReportFocus(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to start program: %w", err)
	}

	return nil
}

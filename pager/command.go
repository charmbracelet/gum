package pager

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
)

// Run provides a shell script interface for the viewport bubble.
// https://github.com/charmbracelet/bubbles/viewport
func (o Options) Run() error {
	vp := viewport.New(
		viewport.WithWidth(o.Style.Width),
		viewport.WithHeight(o.Style.Height),
	)
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
		vp.LeftGutterFunc = func(info viewport.GutterContext) string {
			style := o.LineNumberStyle.ToLipgloss()
			if info.Soft {
				return style.Render("     │ ")
			}
			if info.Index >= info.TotalLines {
				return style.Render("   ~ │ ")
			}
			// TODO: handle more lines
			return style.Render(fmt.Sprintf("%4d │ ", info.Index+1))
		}
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

	_, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithReportFocus(),
		tea.WithMouseAllMotion(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to start program: %w", err)
	}

	return nil
}

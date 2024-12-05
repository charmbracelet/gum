package write

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/gum/cursor"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/textarea
func (o Options) Run() error {
	in, _ := stdin.Read()
	if in != "" && o.Value == "" {
		o.Value = strings.ReplaceAll(in, "\r", "")
	}

	a := textarea.New()
	a.Focus()

	a.Prompt = o.Prompt
	a.Placeholder = o.Placeholder
	a.ShowLineNumbers = o.ShowLineNumbers
	a.CharLimit = o.CharLimit

	style := textarea.Style{
		Base:             o.BaseStyle.ToLipgloss(),
		Placeholder:      o.PlaceholderStyle.ToLipgloss(),
		CursorLine:       o.CursorLineStyle.ToLipgloss(),
		CursorLineNumber: o.CursorLineNumberStyle.ToLipgloss(),
		EndOfBuffer:      o.EndOfBufferStyle.ToLipgloss(),
		LineNumber:       o.LineNumberStyle.ToLipgloss(),
		Prompt:           o.PromptStyle.ToLipgloss(),
	}

	a.BlurredStyle = style
	a.FocusedStyle = style
	a.Cursor.Style = o.CursorStyle.ToLipgloss()
	a.Cursor.SetMode(cursor.Modes[o.CursorMode])

	a.SetWidth(o.Width)
	a.SetHeight(o.Height)
	a.SetValue(o.Value)

	p := tea.NewProgram(model{
		textarea:    a,
		header:      o.Header,
		headerStyle: o.HeaderStyle.ToLipgloss(),
		autoWidth:   o.Width < 1,
	}, tea.WithOutput(os.Stderr))
	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run write: %w", err)
	}
	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}

	fmt.Println(m.textarea.Value())
	return nil
}

package write

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/cursor"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/textarea
func (o Options) Run() error {
	in, _ := stdin.Read(stdin.StripANSI(o.StripANSI))
	if in != "" && o.Value == "" {
		o.Value = strings.ReplaceAll(in, "\r", "")
	}

	a := textarea.New()
	a.Focus()

	a.Prompt = o.Prompt
	a.Placeholder = o.Placeholder
	a.ShowLineNumbers = o.ShowLineNumbers
	a.CharLimit = o.CharLimit
	a.MaxHeight = o.MaxLines
	top, right, bottom, left := style.ParsePadding(o.Padding)

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

	a.SetWidth(max(0, o.Width-left-right))
	a.SetHeight(max(0, o.Height-top-bottom))
	a.SetValue(o.Value)

	m := model{
		textarea:    a,
		header:      o.Header,
		headerStyle: o.HeaderStyle.ToLipgloss(),
		autoWidth:   o.Width < 1,
		help:        help.New(),
		showHelp:    o.ShowHelp,
		keymap:      defaultKeymap(),
		padding:     []int{top, right, bottom, left},
	}

	m.textarea.KeyMap.InsertNewline = m.keymap.InsertNewline

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	p := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithReportFocus(),
		tea.WithContext(ctx),
	)
	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run write: %w", err)
	}
	m = tm.(model)
	if !m.submitted {
		return errors.New("not submitted")
	}
	fmt.Println(m.textarea.Value())
	return nil
}

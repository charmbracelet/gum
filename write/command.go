package write

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
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

	style := textarea.StyleState{
		Base:             o.BaseStyle.ToLipgloss(),
		Placeholder:      o.PlaceholderStyle.ToLipgloss(),
		CursorLine:       o.CursorLineStyle.ToLipgloss(),
		CursorLineNumber: o.CursorLineNumberStyle.ToLipgloss(),
		EndOfBuffer:      o.EndOfBufferStyle.ToLipgloss(),
		LineNumber:       o.LineNumberStyle.ToLipgloss(),
		Prompt:           o.PromptStyle.ToLipgloss(),
	}

	a.Styles.Focused = style
	a.Styles.Blurred = style
	a.Styles.Cursor.Color = o.CursorStyle.ToLipgloss().GetForeground()
	a.Styles.Cursor.Blink = o.CursorMode == "blink"
	// TODO: handle cursor hidden
	// a.Cursor.SetMode(cursor.Modes[o.CursorMode])

	a.SetWidth(o.Width)
	a.SetHeight(o.Height)
	a.SetValue(o.Value)

	m := model{
		textarea:    a,
		header:      o.Header,
		headerStyle: o.HeaderStyle.ToLipgloss(),
		autoWidth:   o.Width < 1,
		help:        help.New(),
		showHelp:    o.ShowHelp,
		keymap:      defaultKeymap(),
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

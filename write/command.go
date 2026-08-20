package write

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/gum/v2/cursor"
	"charm.land/gum/v2/internal/stdin"
	"charm.land/gum/v2/internal/timeout"
	"charm.land/gum/v2/style"
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

	taStyles := textarea.StyleState{
		Base:             o.BaseStyle.ToLipgloss(),
		Placeholder:      o.PlaceholderStyle.ToLipgloss(),
		CursorLine:       o.CursorLineStyle.ToLipgloss(),
		CursorLineNumber: o.CursorLineNumberStyle.ToLipgloss(),
		EndOfBuffer:      o.EndOfBufferStyle.ToLipgloss(),
		LineNumber:       o.LineNumberStyle.ToLipgloss(),
		Prompt:           o.PromptStyle.ToLipgloss(),
	}

	styles := a.Styles()
	styles.Focused = taStyles
	styles.Blurred = taStyles
	a.SetStyles(styles)
	cursor.TextArea(&a, o.CursorMode, o.CursorStyle)

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

package input

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/gum/cursor"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() error {
	i := textinput.New()
	if in, _ := stdin.Read(); in != "" && o.Value == "" {
		i.SetValue(in)
	} else {
		i.SetValue(o.Value)
	}

	i.Focus()
	i.Prompt = o.Prompt
	i.Placeholder = o.Placeholder
	i.Width = o.Width
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.Cursor.Style = o.CursorStyle.ToLipgloss()
	i.Cursor.SetMode(cursor.Modes[o.CursorMode])
	i.CharLimit = o.CharLimit

	if o.Password {
		i.EchoMode = textinput.EchoPassword
		i.EchoCharacter = '•'
	}

	p := tea.NewProgram(model{
		textinput:   i,
		aborted:     false,
		header:      o.Header,
		headerStyle: o.HeaderStyle.ToLipgloss(),
		timeout:     o.Timeout,
		hasTimeout:  o.Timeout > 0,
		autoWidth:   o.Width < 1,
	}, tea.WithOutput(os.Stderr))
	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run input: %w", err)
	}
	m := tm.(model)

	if m.aborted {
		return exit.ErrAborted
	}

	fmt.Println(m.textinput.Value())
	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

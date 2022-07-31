package input

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	i.CursorStyle = o.CursorStyle.ToLipgloss()

	if o.Password {
		i.EchoMode = textinput.EchoPassword
		i.EchoCharacter = '•'
	}

	p := tea.NewProgram(model{
		textinput: i,
		aborted:   false,
	}, tea.WithOutput(os.Stderr))
	tm, err := p.StartReturningModel()
	m := tm.(model)

	if m.aborted {
		os.Exit(exit.Aborted)
	}

	fmt.Println(m.textinput.Value())
	return err
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	return style.HideFlags(ctx)
}

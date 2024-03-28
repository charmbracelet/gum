package input

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() error {
	i := textinput.New()
	if o.Value != "" {
		i.SetValue(o.Value)
	} else if in, _ := stdin.Read(); in != "" {
		i.SetValue(in)
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle()
	theme.Focused.TextInput.Cursor = o.CursorStyle.ToLipgloss()
	theme.Focused.TextInput.Placeholder = o.PlaceholderStyle.ToLipgloss()
	theme.Focused.TextInput.Prompt = o.PromptStyle.ToLipgloss()
	theme.Focused.Title = o.HeaderStyle.ToLipgloss()

	var value string
	var echoMode huh.EchoMode

	if o.Password {
		echoMode = huh.EchoModePassword
	} else {
		echoMode = huh.EchoModeNormal
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Prompt(o.Prompt).
				Placeholder(o.Placeholder).
				CharLimit(o.CharLimit).
				EchoMode(echoMode).
				Title(o.Header).
				Value(&value),
		),
	).
		WithShowHelp(false).
		WithTheme(theme).
		Run()

	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

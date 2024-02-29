package input

import (
	"fmt"

	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() error {
	if o.Value == "" {
		o.Value, _ = stdin.Read()
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base.Border(lipgloss.Border{})
	theme.Focused.Title = o.HeaderStyle.ToLipgloss()
	theme.Focused.TextInput.Prompt = o.PromptStyle.ToLipgloss()
	theme.Focused.TextInput.Placeholder = o.PlaceholderStyle.ToLipgloss()
	theme.Focused.TextInput.Cursor = o.CursorStyle.ToLipgloss()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Password(o.Password).
				Title(o.Header).
				Prompt(o.Prompt).
				CharLimit(o.CharLimit).
				Placeholder(o.Placeholder).
				Value(&o.Value),
		),
	).
		WithWidth(o.Width).
		WithShowHelp(false).
		WithTheme(theme).
		Run()
	if err != nil {
		return err
	}

	if o.Value != "" {
		fmt.Println(o.Value)
	}
	return nil
}

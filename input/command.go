package input

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() error {
	var value string
	if o.Value != "" {
		value = o.Value
	} else if in, _ := stdin.Read(); in != "" {
		value = in
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle()
	// theme.Focused.TextInput.Cursor = o.CursorStyle.ToLipgloss()
	theme.Focused.TextInput.Placeholder = o.PlaceholderStyle.ToLipgloss()
	theme.Focused.TextInput.Prompt = o.PromptStyle.ToLipgloss()
	theme.Focused.Title = o.HeaderStyle.ToLipgloss()

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
		WithWidth(o.Width).
		WithTheme(theme).
		WithProgramOptions(tea.WithOutput(os.Stderr)).
		Run()

	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

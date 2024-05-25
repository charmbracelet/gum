package write

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/huh"
)

// Run provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/textarea
func (o Options) Run() error {
	in, _ := stdin.Read()
	if in != "" && o.Value == "" {
		o.Value = strings.ReplaceAll(in, "\r", "")
	}

	var value = o.Value

	theme := huh.ThemeCharm()
	theme.Focused.Base = o.BaseStyle.ToLipgloss()
	theme.Focused.TextInput.Cursor = o.CursorStyle.ToLipgloss()
	theme.Focused.Title = o.HeaderStyle.ToLipgloss()
	theme.Focused.TextInput.Placeholder = o.PlaceholderStyle.ToLipgloss()
	theme.Focused.TextInput.Prompt = o.PromptStyle.ToLipgloss()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(o.Header).
				Placeholder(o.Placeholder).
				CharLimit(o.CharLimit).
				ShowLineNumbers(o.ShowLineNumbers).
				Value(&value),
		),
	).
		WithWidth(o.Width).
		WithHeight(o.Height).
		WithTheme(theme).
		WithShowHelp(false).Run()

	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

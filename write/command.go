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
	if o.Value == "" {
		o.Value, _ = stdin.Read()
	}
	o.Value = strings.ReplaceAll(o.Value, "\r", "")
	theme := huh.ThemeCharm()
	theme.Focused.Title = o.HeaderStyle.ToLipgloss()
	theme.Focused.TextInput.Placeholder = o.PlaceholderStyle.ToLipgloss()
	theme.Focused.TextInput.Cursor = o.CursorStyle.ToLipgloss()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(o.Header).
				ShowLineNumbers(o.ShowLineNumbers).
				CharLimit(o.CharLimit).
				Placeholder(o.Placeholder).
				Value(&o.Value),
		),
	).
		WithShowHelp(false).
		WithWidth(o.Width).
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

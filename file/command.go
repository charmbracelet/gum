package file

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Run is the interface to picking a file.
func (o Options) Run() error {
	if !o.File && !o.Directory {
		return errors.New("at least one between --file and --directory must be set")
	}

	if o.Path == "" {
		o.Path = "."
	}

	path, err := filepath.Abs(o.Path)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle()
	theme.Focused.File = o.FileStyle.ToLipgloss()
	theme.Focused.Directory = o.DirectoryStyle.ToLipgloss()
	theme.Focused.SelectedOption = o.SelectedStyle.ToLipgloss()

	keymap := huh.NewDefaultKeyMap()
	keymap.FilePicker.Open.SetEnabled(false)

	// XXX: These should be file selected specific.
	theme.Focused.TextInput.Placeholder = o.PermissionsStyle.ToLipgloss()
	theme.Focused.TextInput.Prompt = o.CursorStyle.ToLipgloss()

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Picking(true).
				CurrentDirectory(path).
				DirAllowed(o.Directory).
				FileAllowed(o.File).
				Height(o.Height).
				ShowHidden(o.All).
				Value(&path),
		),
	).
		WithTimeout(o.Timeout).
		WithShowHelp(o.ShowHelp).
		WithKeyMap(keymap).
		WithTheme(theme).
		Run()
	if err != nil {
		return exit.Handle(err, o.Timeout)
	}
	fmt.Println(path)
	return nil
}

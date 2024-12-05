package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
)

type keymap filepicker.KeyMap

func defaultKeymap() keymap {
	km := filepicker.DefaultKeyMap()
	km.Down.SetHelp("↓", "down")
	km.Up.SetHelp("↑", "up")
	return keymap(km)
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Up,
			k.Down,
			k.Back,
			k.Open,
			k.Select,
		},
		{
			k.GoToTop,
			k.GoToLast,
			k.PageUp,
			k.PageDown,
		},
	}
}

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		key.NewBinding(key.WithHelp("esc", "close")),
		k.Select,
	}
}

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

	fp := filepicker.New()
	fp.CurrentDirectory = path
	fp.Path = path
	fp.Height = o.Height
	fp.AutoHeight = o.Height == 0
	fp.Cursor = o.Cursor
	fp.DirAllowed = o.Directory
	fp.FileAllowed = o.File
	fp.ShowPermissions = o.Permissions
	fp.ShowSize = o.Size
	fp.ShowHidden = o.All
	fp.Styles = filepicker.DefaultStyles()
	fp.Styles.Cursor = o.CursorStyle.ToLipgloss()
	fp.Styles.Symlink = o.SymlinkStyle.ToLipgloss()
	fp.Styles.Directory = o.DirectoryStyle.ToLipgloss()
	fp.Styles.File = o.FileStyle.ToLipgloss()
	fp.Styles.Permission = o.PermissionsStyle.ToLipgloss()
	fp.Styles.Selected = o.SelectedStyle.ToLipgloss()
	fp.Styles.FileSize = o.FileSizeStyle.ToLipgloss()
	m := model{
		filepicker: fp,
		timeout:    o.Timeout,
		hasTimeout: o.Timeout > 0,
		aborted:    false,
		showHelp:   o.ShowHelp,
		help:       help.New(),
	}

	tm, err := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return fmt.Errorf("unable to pick selection: %w", err)
	}
	m = tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}
	if m.timedOut {
		return exit.ErrTimeout
	}
	if m.selectedPath == "" {
		return errors.New("no file selected")
	}

	fmt.Println(m.selectedPath)
	return nil
}

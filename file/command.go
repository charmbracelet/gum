package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
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

	fp := filepicker.New()
	fp.CurrentDirectory = path
	fp.Path = path
	fp.Height = o.Height
	fp.AutoHeight = o.Height == 0
	fp.Cursor = o.Cursor
	fp.DirAllowed = o.Directory
	fp.FileAllowed = o.File
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

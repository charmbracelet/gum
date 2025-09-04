package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/style"
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
	fp.SetHeight(o.Height)
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
	top, right, bottom, left := style.ParsePadding(o.Padding)
	m := model{
		filepicker:  fp,
		padding:     []int{top, right, bottom, left},
		showHelp:    o.ShowHelp,
		help:        help.New(),
		keymap:      defaultKeymap(),
		headerStyle: o.HeaderStyle.ToLipgloss(),
		header:      o.Header,
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	tm, err := tea.NewProgram(
		&m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to pick selection: %w", err)
	}
	m = tm.(model)
	if m.selectedPath == "" {
		return errors.New("no file selected")
	}

	fmt.Println(m.selectedPath)
	return nil
}

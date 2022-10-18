package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stack"
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

	m := model{
		path:            path,
		cursor:          o.Cursor,
		selected:        0,
		showHidden:      o.All,
		dirAllowed:      o.Directory,
		fileAllowed:     o.File,
		autoHeight:      o.Height == 0,
		height:          o.Height,
		max:             0,
		min:             0,
		selectedStack:   stack.NewStack(),
		minStack:        stack.NewStack(),
		maxStack:        stack.NewStack(),
		cursorStyle:     o.CursorStyle.ToLipgloss().Inline(true),
		symlinkStyle:    o.SymlinkStyle.ToLipgloss().Inline(true),
		directoryStyle:  o.DirectoryStyle.ToLipgloss().Inline(true),
		fileStyle:       o.FileStyle.ToLipgloss().Inline(true),
		permissionStyle: o.PermissionsStyle.ToLipgloss().Inline(true),
		selectedStyle:   o.SelectedStyle.ToLipgloss().Inline(true),
		fileSizeStyle:   o.FileSizeStyle.ToLipgloss().Inline(true),
	}

	tm, err := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return fmt.Errorf("unable to pick selection: %w", err)
	}

	m = tm.(model)

	if m.path == "" {
		os.Exit(1)
	}

	fmt.Println(m.path)

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/internal/tty"
	"github.com/charmbracelet/gum/style"
	"github.com/sahilm/fuzzy"
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

	// Fuzzy-find mode: flat list with fuzzy matching
	if o.Fuzzy {
		return o.runFuzzy(path)
	}

	// Default tree-navigation mode
	return o.runTreePicker(path)
}

func (o Options) runTreePicker(path string) error {
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

func (o Options) runFuzzy(basePath string) error {
	allFiles := files.ListFrom(basePath, o.All, o.Directory)
	if len(allFiles) == 0 {
		return errors.New("no files found")
	}

	i := textinput.New()
	i.Focus()
	i.Prompt = "> "
	i.PromptStyle = o.CursorStyle.ToLipgloss()
	i.Placeholder = "Type to filter files..."

	height := o.Height
	if height == 0 {
		height = 20
	}
	v := viewport.New(0, height)

	matches := matchAll(allFiles)
	top, right, bottom, left := style.ParsePadding(o.Padding)

	fm := fuzzyModel{
		basePath:   basePath,
		allFiles:   allFiles,
		matches:    matches,
		textinput:  i,
		viewport:   &v,
		matchStyle: o.SelectedStyle.ToLipgloss(),
		fileStyle:  o.FileStyle.ToLipgloss(),
		dirStyle:   o.DirectoryStyle.ToLipgloss(),
		padding:    []int{top, right, bottom, left},
		height:     height,
		header:     o.Header,
		showHelp:   o.ShowHelp,
		help:       help.New(),
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	options := []tea.ProgramOption{
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	}
	if height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	tm, err := tea.NewProgram(&fm, options...).Run()
	if err != nil {
		return fmt.Errorf("unable to run fuzzy file picker: %w", err)
	}
	fm = tm.(fuzzyModel)
	if fm.selectedPath == "" {
		return errors.New("no file selected")
	}

	tty.Println(fm.selectedPath)
	return nil
}

func matchAll(choices []string) []fuzzy.Match {
	matches := make([]fuzzy.Match, len(choices))
	for i, c := range choices {
		matches[i] = fuzzy.Match{Str: c, Index: i}
	}
	return matches
}

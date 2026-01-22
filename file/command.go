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

	// Use fuzzy mode if enabled
	if o.Fuzzy {
		return o.runFuzzy(path)
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

// runFuzzy runs the file picker in fuzzy search mode.
func (o Options) runFuzzy(basePath string) error {
	// Collect all files recursively
	files, err := collectFiles(basePath, o.All)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) == 0 {
		return errors.New("no files found in directory")
	}

	// Build the list of file names for fuzzy matching
	filteringNames := make([]string, len(files))
	for i, f := range files {
		filteringNames[i] = f.name
	}

	// Initialize text input
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = o.Prompt
	ti.PromptStyle = o.PromptStyle.ToLipgloss()
	ti.PlaceholderStyle = o.PlaceholderStyle.ToLipgloss()
	ti.Placeholder = o.Placeholder

	// Initialize viewport
	v := viewport.New(0, o.Height)

	// Initial matches (show all)
	matches := matchAllFuzzy(filteringNames)

	top, right, bottom, left := style.ParsePadding(o.Padding)
	m := fuzzyModel{
		textinput:        ti,
		viewport:         &v,
		files:            files,
		filteringNames:   filteringNames,
		matches:          matches,
		cursor:           0,
		header:           o.Header,
		height:           o.Height,
		padding:          []int{top, right, bottom, left},
		showHelp:         o.ShowHelp,
		showPermissions:  o.Permissions,
		showSize:         o.Size,
		dirAllowed:       o.Directory,
		fileAllowed:      o.File,
		basePath:         basePath,
		headerStyle:      o.HeaderStyle.ToLipgloss(),
		matchStyle:       o.MatchStyle.ToLipgloss(),
		indicatorStyle:   o.IndicatorStyle.ToLipgloss(),
		indicator:        o.Indicator,
		cursorStyle:      o.CursorStyle.ToLipgloss(),
		selectedStyle:    o.SelectedStyle.ToLipgloss(),
		directoryStyle:   o.DirectoryStyle.ToLipgloss(),
		fileStyle:        o.FileStyle.ToLipgloss(),
		symlinkStyle:     o.SymlinkStyle.ToLipgloss(),
		permissionsStyle: o.PermissionsStyle.ToLipgloss(),
		fileSizeStyle:    o.FileSizeStyle.ToLipgloss(),
		keymap:           defaultFuzzyKeymap(),
		help:             help.New(),
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	options := []tea.ProgramOption{
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	}
	if o.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	tm, err := tea.NewProgram(&m, options...).Run()
	if err != nil {
		return fmt.Errorf("unable to run fuzzy file picker: %w", err)
	}

	fm := tm.(fuzzyModel)
	if fm.selectedPath == "" {
		return errors.New("no file selected")
	}

	fmt.Println(fm.selectedPath)
	return nil
}

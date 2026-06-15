// Package file provides an interface to pick a file from a folder (tree).
// This file implements the fuzzy find functionality for the file picker.
package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/ordered"
	"github.com/dustin/go-humanize"
	"github.com/rivo/uniseg"
	"github.com/sahilm/fuzzy"
)

// fileEntry represents a file with its full path and display name.
type fileEntry struct {
	path    string
	name    string
	isDir   bool
	mode    fs.FileMode
	size    int64
	symlink string
}

type fuzzyKeymap struct {
	Down   key.Binding
	Up     key.Binding
	Submit key.Binding
	Quit   key.Binding
	Abort  key.Binding
}

func defaultFuzzyKeymap() fuzzyKeymap {
	return fuzzyKeymap{
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+j", "ctrl+n"),
			key.WithHelp("down", "next"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+k", "ctrl+p"),
			key.WithHelp("up", "prev"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
	}
}

// FullHelp implements help.KeyMap.
func (k fuzzyKeymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k fuzzyKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("up/down", "navigate"),
		),
		k.Submit,
		k.Quit,
	}
}

type fuzzyModel struct {
	textinput       textinput.Model
	viewport        *viewport.Model
	files           []fileEntry
	filteringNames  []string
	matches         []fuzzy.Match
	cursor          int
	header          string
	height          int
	padding         []int
	quitting        bool
	selectedPath    string
	showHelp        bool
	showPermissions bool
	showSize        bool
	dirAllowed      bool
	fileAllowed     bool
	basePath        string

	// Styles
	headerStyle      lipgloss.Style
	matchStyle       lipgloss.Style
	indicatorStyle   lipgloss.Style
	indicator        string
	cursorStyle      lipgloss.Style
	selectedStyle    lipgloss.Style
	directoryStyle   lipgloss.Style
	fileStyle        lipgloss.Style
	symlinkStyle     lipgloss.Style
	permissionsStyle lipgloss.Style
	fileSizeStyle    lipgloss.Style

	keymap fuzzyKeymap
	help   help.Model
}

func (m fuzzyModel) Init() tea.Cmd { return textinput.Blink }

func (m fuzzyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, icmd tea.Cmd
	m.textinput, icmd = m.textinput.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.height == 0 || m.height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
		}
		if m.header != "" {
			m.viewport.Height -= lipgloss.Height(m.headerStyle.Render(m.header))
		}
		if m.showHelp {
			m.viewport.Height -= lipgloss.Height(m.helpView())
		}
		m.viewport.Height -= m.padding[0] + m.padding[2]
		m.viewport.Width = msg.Width - m.padding[1] - m.padding[3]
		m.textinput.Width = msg.Width - m.padding[1] - m.padding[3]

	case tea.KeyMsg:
		km := m.keymap
		switch {
		case key.Matches(msg, km.Abort):
			m.quitting = true
			return m, tea.Interrupt
		case key.Matches(msg, km.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Submit):
			if len(m.matches) > 0 && m.cursor < len(m.matches) {
				match := m.matches[m.cursor]
				entry := m.findFileByName(match.Str)
				if entry != nil {
					if (entry.isDir && m.dirAllowed) || (!entry.isDir && m.fileAllowed) {
						m.selectedPath = entry.path
						m.quitting = true
						return m, tea.Quit
					}
				}
			}
		case key.Matches(msg, km.Down):
			m.cursorDown()
		case key.Matches(msg, km.Up):
			m.cursorUp()
		default:
			// Text input changed, update matches
			m.matches = fuzzy.Find(m.textinput.Value(), m.filteringNames)
			if m.textinput.Value() == "" {
				m.matches = matchAllFuzzy(m.filteringNames)
			}
		}
	}

	m.cursor = ordered.Clamp(m.cursor, 0, len(m.matches)-1)
	return m, tea.Batch(cmd, icmd)
}

func (m *fuzzyModel) cursorDown() {
	if len(m.matches) == 0 {
		return
	}
	m.cursor = (m.cursor + 1) % len(m.matches)
	if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.ScrollDown(1)
	}
	if m.cursor < m.viewport.YOffset {
		m.viewport.GotoTop()
	}
}

func (m *fuzzyModel) cursorUp() {
	if len(m.matches) == 0 {
		return
	}
	m.cursor = (m.cursor - 1 + len(m.matches)) % len(m.matches)
	if m.cursor < m.viewport.YOffset {
		m.viewport.ScrollUp(1)
	}
	if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(len(m.matches) - m.viewport.Height)
	}
}

func (m fuzzyModel) findFileByName(name string) *fileEntry {
	for i := range m.files {
		if m.files[i].name == name {
			return &m.files[i]
		}
	}
	return nil
}

func (m fuzzyModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	for i := range m.matches {
		match := m.matches[i]
		entry := m.findFileByName(match.Str)
		if entry == nil {
			continue
		}

		// Indicator for current selection
		if i == m.cursor {
			s.WriteString(m.indicatorStyle.Render(m.indicator))
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(m.indicator)))
		}

		s.WriteString(" ")

		// File info
		var lineStyle lipgloss.Style
		if entry.isDir {
			lineStyle = m.directoryStyle
		} else if entry.symlink != "" {
			lineStyle = m.symlinkStyle
		} else {
			lineStyle = m.fileStyle
		}

		if i == m.cursor {
			lineStyle = m.selectedStyle
		}

		// Build the display line
		var line strings.Builder

		if m.showPermissions {
			line.WriteString(m.permissionsStyle.Render(entry.mode.String()))
			line.WriteString(" ")
		}

		if m.showSize {
			sizeStr := formatSize(entry.size)
			line.WriteString(m.fileSizeStyle.Render(sizeStr))
			line.WriteString(" ")
		}

		// Render name with match highlights
		if len(match.MatchedIndexes) == 0 {
			line.WriteString(lineStyle.Render(entry.name))
		} else {
			var ranges []lipgloss.Range
			for _, rng := range matchedRangesFuzzy(match.MatchedIndexes) {
				start, stop := bytePosToVisibleCharPosFuzzy(match.Str, rng)
				ranges = append(ranges, lipgloss.NewRange(start, stop+1, m.matchStyle))
			}
			line.WriteString(lineStyle.Render(lipgloss.StyleRanges(entry.name, ranges...)))
		}

		if entry.symlink != "" {
			line.WriteString(" -> ")
			line.WriteString(entry.symlink)
		}

		s.WriteString(line.String())
		s.WriteRune('\n')
	}

	m.viewport.SetContent(s.String())

	// Build final view
	header := m.headerStyle.Render(m.header)
	view := m.textinput.View() + "\n" + m.viewport.View()
	if m.showHelp {
		view += m.helpView()
	}
	if m.header != "" {
		return lipgloss.NewStyle().
			Padding(m.padding...).
			Render(header + "\n" + view)
	}

	return lipgloss.NewStyle().
		Padding(m.padding...).
		Render(view)
}

func (m fuzzyModel) helpView() string {
	return "\n\n" + m.help.View(m.keymap)
}

func matchAllFuzzy(options []string) []fuzzy.Match {
	matches := make([]fuzzy.Match, len(options))
	for i, option := range options {
		matches[i] = fuzzy.Match{Str: option}
	}
	return matches
}

func matchedRangesFuzzy(in []int) [][2]int {
	if len(in) == 0 {
		return [][2]int{}
	}
	current := [2]int{in[0], in[0]}
	if len(in) == 1 {
		return [][2]int{current}
	}
	var out [][2]int
	for i := 1; i < len(in); i++ {
		if in[i] == current[1]+1 {
			current[1] = in[i]
		} else {
			out = append(out, current)
			current = [2]int{in[i], in[i]}
		}
	}
	out = append(out, current)
	return out
}

func bytePosToVisibleCharPosFuzzy(str string, rng [2]int) (int, int) {
	bytePos, byteStart, byteStop := 0, rng[0], rng[1]
	pos, start, stop := 0, 0, 0
	gr := uniseg.NewGraphemes(str)
	for byteStart > bytePos {
		if !gr.Next() {
			break
		}
		bytePos += len(gr.Str())
		pos += max(1, gr.Width())
	}
	start = pos
	for byteStop > bytePos {
		if !gr.Next() {
			break
		}
		bytePos += len(gr.Str())
		pos += max(1, gr.Width())
	}
	stop = pos
	return start, stop
}

func formatSize(size int64) string {
	sizeStr := humanize.Bytes(uint64(size)) //nolint:gosec
	// Remove the space between number and unit for compactness
	sizeStr = strings.Replace(sizeStr, " ", "", 1)
	// Right-pad to 7 chars for alignment
	return fmt.Sprintf("%7s", sizeStr)
}

// collectFiles recursively collects all files from the given path.
func collectFiles(basePath string, showHidden bool) ([]fileEntry, error) {
	var files []fileEntry

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Skip the base path itself
		if path == basePath {
			return nil
		}

		name := d.Name()

		// Skip hidden files if not showing them
		if !showHidden && strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil // Skip files we can't get info for
		}

		entry := fileEntry{
			path:  path,
			name:  strings.TrimPrefix(path, basePath+string(os.PathSeparator)),
			isDir: d.IsDir(),
			mode:  info.Mode(),
			size:  info.Size(),
		}

		// Check for symlink
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := filepath.EvalSymlinks(path)
			if err == nil {
				entry.symlink = target
			}
		}

		files = append(files, entry)
		return nil
	})

	return files, err
}


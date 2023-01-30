// Package file provides an interface to pick a file from a folder (tree).
// The user is provided a file manager-like interface to navigate, to
// select a file.
//
// Let's pick a file from the current directory:
//
// $ gum file
// $ gum file .
//
// Let's pick a file from the home directory:
//
// $ gum file $HOME
package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stack"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

const marginBottom = 5

type model struct {
	quitting    bool
	path        string
	files       []os.DirEntry
	showHidden  bool
	dirAllowed  bool
	fileAllowed bool

	selected      int
	selectedStack stack.Stack

	min        int
	max        int
	maxStack   stack.Stack
	minStack   stack.Stack
	height     int
	autoHeight bool

	cursor          string
	cursorStyle     lipgloss.Style
	symlinkStyle    lipgloss.Style
	directoryStyle  lipgloss.Style
	fileStyle       lipgloss.Style
	permissionStyle lipgloss.Style
	selectedStyle   lipgloss.Style
	fileSizeStyle   lipgloss.Style
}

type readDirMsg []os.DirEntry

func readDir(path string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return tea.Quit
		}

		sort.Slice(dirEntries, func(i, j int) bool {
			if dirEntries[i].IsDir() == dirEntries[j].IsDir() {
				return dirEntries[i].Name() < dirEntries[j].Name()
			}
			return dirEntries[i].IsDir()
		})

		if showHidden {
			return readDirMsg(dirEntries)
		}

		var sanitizedDirEntries []fs.DirEntry
		for _, dirEntry := range dirEntries {
			isHidden, _ := IsHidden(dirEntry.Name())
			if isHidden {
				continue
			}
			sanitizedDirEntries = append(sanitizedDirEntries, dirEntry)
		}
		return readDirMsg(sanitizedDirEntries)
	}
}

func (m model) Init() tea.Cmd {
	return readDir(m.path, m.showHidden)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case readDirMsg:
		m.files = msg
	case tea.WindowSizeMsg:
		if m.autoHeight {
			m.height = msg.Height - marginBottom
		}
		m.max = m.height
	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			m.selected = 0
			m.min = 0
			m.max = m.height - 1
		case "G":
			m.selected = len(m.files) - 1
			m.min = len(m.files) - m.height
			m.max = len(m.files) - 1
		case "j", "down":
			m.selected++
			if m.selected >= len(m.files) {
				m.selected = len(m.files) - 1
			}
			if m.selected > m.max {
				m.min++
				m.max++
			}
		case "k", "up":
			m.selected--
			if m.selected < 0 {
				m.selected = 0
			}
			if m.selected < m.min {
				m.min--
				m.max--
			}
		case "J", "pgdown":
			m.selected += m.height
			if m.selected >= len(m.files) {
				m.selected = len(m.files) - 1
			}
			m.min += m.height
			m.max += m.height

			if m.max >= len(m.files) {
				m.max = len(m.files) - 1
				m.min = m.max - m.height
			}
		case "K", "pgup":
			m.selected -= m.height
			if m.selected < 0 {
				m.selected = 0
			}
			m.min -= m.height
			m.max -= m.height

			if m.min < 0 {
				m.min = 0
				m.max = m.min + m.height
			}
		case "ctrl+c", "q":
			m.path = ""
			m.quitting = true
			return m, tea.Quit
		case "backspace", "h", "left":
			m.path = filepath.Dir(m.path)
			if m.selectedStack.Length() > 0 {
				m.selected, m.min, m.max = m.popView()
			} else {
				m.selected = 0
				m.min = 0
				m.max = m.height - 1
			}
			return m, readDir(m.path, m.showHidden)
		case "l", "right", "enter":
			if len(m.files) == 0 {
				break
			}

			f := m.files[m.selected]
			info, err := f.Info()
			if err != nil {
				break
			}
			isSymlink := info.Mode()&fs.ModeSymlink != 0
			isDir := f.IsDir()

			if isSymlink {
				symlinkPath, _ := filepath.EvalSymlinks(filepath.Join(m.path, f.Name()))
				info, err := os.Stat(symlinkPath)
				if err != nil {
					break
				}
				if info.IsDir() {
					isDir = true
				}
			}

			if (!isDir && m.fileAllowed) || (isDir && m.dirAllowed) {
				if msg.String() == "enter" {
					m.path = filepath.Join(m.path, f.Name())
					m.quitting = true
					return m, tea.Quit
				}
			}

			if !isDir {
				break
			}

			m.path = filepath.Join(m.path, f.Name())
			m.pushView()
			m.selected = 0
			m.min = 0
			m.max = m.height - 1
			return m, readDir(m.path, m.showHidden)
		}
	}
	return m, nil
}

func (m model) pushView() {
	m.minStack.Push(m.min)
	m.maxStack.Push(m.max)
	m.selectedStack.Push(m.selected)
}

func (m model) popView() (int, int, int) {
	return m.selectedStack.Pop(), m.minStack.Pop(), m.maxStack.Pop()
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if len(m.files) == 0 {
		return "Bummer. No files found."
	}
	var s strings.Builder

	for i, f := range m.files {
		if i < m.min {
			continue
		}
		if i > m.max {
			break
		}

		var symlinkPath string
		info, _ := f.Info()
		isSymlink := info.Mode()&fs.ModeSymlink != 0
		size := humanize.Bytes(uint64(info.Size()))
		name := f.Name()

		if isSymlink {
			symlinkPath, _ = filepath.EvalSymlinks(filepath.Join(m.path, name))
		}

		if m.selected == i {
			selected := fmt.Sprintf(" %s %"+fmt.Sprint(m.fileSizeStyle.GetWidth())+"s %s", info.Mode().String(), size, name)
			if isSymlink {
				selected = fmt.Sprintf("%s → %s", selected, symlinkPath)
			}
			s.WriteString(m.cursorStyle.Render(m.cursor) + m.selectedStyle.Render(selected))
			s.WriteRune('\n')
			continue
		}

		var style = m.fileStyle
		if f.IsDir() {
			style = m.directoryStyle
		} else if isSymlink {
			style = m.symlinkStyle
		}

		fileName := style.Render(name)
		if isSymlink {
			fileName = fmt.Sprintf("%s → %s", fileName, symlinkPath)
		}
		s.WriteString(fmt.Sprintf("  %s %s %s", m.permissionStyle.Render(info.Mode().String()), m.fileSizeStyle.Render(size), fileName))
		s.WriteRune('\n')
	}

	return s.String()
}

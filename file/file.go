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
<<<<<<< HEAD
	"io/fs"
=======
>>>>>>> next
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
<<<<<<< HEAD
	quitting   bool
	path       string
	files      []os.DirEntry
	showHidden bool
=======
	quitting bool
	path     string
	files    []os.DirEntry
>>>>>>> next

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
<<<<<<< HEAD
	symlinkStyle    lipgloss.Style
=======
>>>>>>> next
	directoryStyle  lipgloss.Style
	fileStyle       lipgloss.Style
	permissionStyle lipgloss.Style
	selectedStyle   lipgloss.Style
	fileSizeStyle   lipgloss.Style
}

type readDirMsg []os.DirEntry

<<<<<<< HEAD
func readDir(path string, showHidden bool) tea.Cmd {
=======
func readDir(path string) tea.Cmd {
>>>>>>> next
	return func() tea.Msg {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return tea.Quit
		}
<<<<<<< HEAD

=======
>>>>>>> next
		sort.Slice(dirEntries, func(i, j int) bool {
			if dirEntries[i].IsDir() == dirEntries[j].IsDir() {
				return dirEntries[i].Name() < dirEntries[j].Name()
			}
			return dirEntries[i].IsDir()
		})
<<<<<<< HEAD

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
=======
		return readDirMsg(dirEntries)
>>>>>>> next
	}
}

func (m model) Init() tea.Cmd {
<<<<<<< HEAD
	return readDir(m.path, m.showHidden)
=======
	return readDir(m.path)
>>>>>>> next
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
<<<<<<< HEAD
			return m, readDir(m.path, m.showHidden)
=======
			return m, readDir(m.path)
>>>>>>> next
		case "l", "right", "enter":
			if len(m.files) == 0 {
				break
			}
			if !m.files[m.selected].IsDir() {
				if msg.String() == "enter" {
					m.path = filepath.Join(m.path, m.files[m.selected].Name())
					m.quitting = true
					return m, tea.Quit
				}
				break
			}
			m.path = filepath.Join(m.path, m.files[m.selected].Name())
<<<<<<< HEAD
			m.pushView()
			m.selected = 0
			m.min = 0
			m.max = m.height - 1
			return m, readDir(m.path, m.showHidden)
=======
			m.pushView(m.selected, m.min, m.max)
			m.selected = 0
			m.min = 0
			m.max = m.height - 1
			return m, readDir(m.path)
>>>>>>> next
		}
	}
	return m, nil
}

<<<<<<< HEAD
func (m model) pushView() {
=======
func (m model) pushView(selected, min, max int) {
>>>>>>> next
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

<<<<<<< HEAD
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
		} else {
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
=======
		info, _ := f.Info()
		size := humanize.Bytes(uint64(info.Size()))
		if m.selected == i {
			s.WriteString(m.cursorStyle.Render(m.cursor) + m.selectedStyle.Render(fmt.Sprintf(" %s %"+fmt.Sprint(m.fileSizeStyle.GetWidth())+"s %s", info.Mode().String(), size, f.Name())))
		} else {
			var style lipgloss.Style
			if f.IsDir() {
				style = m.directoryStyle
			} else {
				style = m.fileStyle
			}

			s.WriteString(fmt.Sprintf("  %s %s %s", m.permissionStyle.Render(info.Mode().String()), m.fileSizeStyle.Render(size), style.Render(f.Name())))
>>>>>>> next
		}
		s.WriteString("\n")
	}

	return s.String()
}

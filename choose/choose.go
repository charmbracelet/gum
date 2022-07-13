// Package choose provides an interface to choose one option from a given list
// of options. The options can be provided as (new-line separated) stdin or a
// list of arguments.
//
// It is different from the filter command as it does not provide a fuzzy
// finding input, so it is best used for smaller lists of options.
//
// Let's pick from a list of gum flavors:
//
//   $ gum choose "Strawberry" "Banana" "Cherry"
//
package choose

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type model struct {
	height           int
	page             int
	cursor           string
	selectedPrefix   string
	unselectedPrefix string
	cursorPrefix     string
	items            []item
	quitting         bool
	index            int
	limit            int
	numSelected      int

	// styles
	cursorStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
}

type item struct {
	text     string
	selected bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "down", "j", "ctrl+n":
			m.index = (m.index + 1) % len(m.items)
			m.page = m.index / m.height
		case "up", "k", "ctrl+p":
			m.index = (m.index - 1 + len(m.items)) % len(m.items)
			m.page = m.index / m.height
		case "right", "l", "ctrl+f":
			if m.index+m.height < len(m.items) {
				m.index += m.height
			} else {
				if m.page < len(m.items)/m.height {
					m.index = len(m.items) - 1
				}
			}
			m.page = m.index / m.height
		case "left", "h", "ctrl+b":
			if m.index-m.height >= 0 {
				m.index -= m.height
			}
			m.page = m.index / m.height
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case " ", "x":
			if m.limit == 1 {
				break // no op
			}

			if m.items[m.index].selected {
				m.items[m.index].selected = false
				m.numSelected--
			} else if m.numSelected < m.limit {
				m.items[m.index].selected = true
				m.numSelected++
			}
		case "enter":
			m.quitting = true
			// Select the item on which they've hit enter if it falls within
			// the limit.
			if m.numSelected < m.limit {
				m.items[m.index].selected = true
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	for i, item := range m.items[clamp(m.page*m.height, 0, len(m.items)):clamp((m.page+1)*m.height, 0, len(m.items))] {
		if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursor))
		} else {
			s.WriteString(strings.Repeat(" ", runewidth.StringWidth(m.cursor)))
		}

		if item.selected {
			s.WriteString(m.selectedItemStyle.Render(m.selectedPrefix + item.text))
		} else if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursorPrefix + item.text))
		} else {
			s.WriteString(m.itemStyle.Render(m.unselectedPrefix + item.text))
		}
		s.WriteRune('\n')
	}

	return s.String()
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// Package choose provides an interface to choose one option from a given list
// of options. The options can be provided as (new-line separated) stdin or a
// list of arguments.
//
// It is different from the filter command as it does not provide a fuzzy
// finding input, so it is best used for smaller lists of options.
//
// Let's pick from a list of gum flavors:
//
// $ gum choose "Strawberry" "Banana" "Cherry"
package choose

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/timeout"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	height           int
	cursor           string
	selectedPrefix   string
	unselectedPrefix string
	cursorPrefix     string
	header           string
	items            []item
	quitting         bool
	index            int
	limit            int
	numSelected      int
	currentOrder     int
	paginator        paginator.Model
	aborted          bool
	timedOut         bool

	// styles
	cursorStyle       lipgloss.Style
	headerStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	hasTimeout        bool
	timeout           time.Duration
}

type item struct {
	text     string
	selected bool
	order    int
}

func (m model) Init() tea.Cmd {
	return timeout.Init(m.timeout, nil)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.quitting = true
			m.timedOut = true
			// If the user hasn't selected any items in a multi-select.
			// Then we select the item that they have pressed enter on. If they
			// have selected items, then we simply return them.
			if m.numSelected < 1 {
				m.items[m.index].selected = true
			}
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	case tea.KeyMsg:
		start, end := m.paginator.GetSliceBounds(len(m.items))
		switch keypress := msg.String(); keypress {
		case "down", "j", "ctrl+j", "ctrl+n":
			m.index++
			if m.index >= len(m.items) {
				m.index = 0
				m.paginator.Page = 0
			}
			if m.index >= end {
				m.paginator.NextPage()
			}
		case "up", "k", "ctrl+k", "ctrl+p":
			m.index--
			if m.index < 0 {
				m.index = len(m.items) - 1
				m.paginator.Page = m.paginator.TotalPages - 1
			}
			if m.index < start {
				m.paginator.PrevPage()
			}
		case "right", "l", "ctrl+f":
			m.index = clamp(m.index+m.height, 0, len(m.items)-1)
			m.paginator.NextPage()
		case "left", "h", "ctrl+b":
			m.index = clamp(m.index-m.height, 0, len(m.items)-1)
			m.paginator.PrevPage()
		case "G", "end":
			m.index = len(m.items) - 1
			m.paginator.Page = m.paginator.TotalPages - 1
		case "g", "home":
			m.index = 0
			m.paginator.Page = 0
		case "a":
			if m.limit <= 1 {
				break
			}
			for i := range m.items {
				if m.numSelected >= m.limit {
					break // do not exceed given limit
				}
				if m.items[i].selected {
					continue
				}
				m.items[i].selected = true
				m.items[i].order = m.currentOrder
				m.numSelected++
				m.currentOrder++
			}
		case "A":
			if m.limit <= 1 {
				break
			}
			for i := range m.items {
				m.items[i].selected = false
				m.items[i].order = 0
			}
			m.numSelected = 0
			m.currentOrder = 0
		case "ctrl+c", "esc":
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case " ", "tab", "x", "ctrl+@":
			if m.limit == 1 {
				break // no op
			}

			if m.items[m.index].selected {
				m.items[m.index].selected = false
				m.numSelected--
			} else if m.numSelected < m.limit {
				m.items[m.index].selected = true
				m.items[m.index].order = m.currentOrder
				m.numSelected++
				m.currentOrder++
			}
		case "enter":
			m.quitting = true
			if m.limit <= 1 && m.numSelected < 1 {
				m.items[m.index].selected = true
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder
	var timeoutStr string

	start, end := m.paginator.GetSliceBounds(len(m.items))
	for i, item := range m.items[start:end] {
		if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursor))
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(m.cursor)))
		}

		if item.selected {
			if m.hasTimeout {
				timeoutStr = timeout.Str(m.timeout)
			}
			s.WriteString(m.selectedItemStyle.Render(m.selectedPrefix + item.text + timeoutStr))
		} else if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursorPrefix + item.text))
		} else {
			s.WriteString(m.itemStyle.Render(m.unselectedPrefix + item.text))
		}
		if i != m.height {
			s.WriteRune('\n')
		}
	}

	if m.paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.height-m.paginator.ItemsOnPage(len(m.items))+1))
		s.WriteString("  " + m.paginator.View())
	}

	if m.header != "" {
		header := m.headerStyle.Render(m.header)
		return lipgloss.JoinVertical(lipgloss.Left, header, s.String())
	}

	return s.String()
}

func clamp(x, low, high int) int {
	if x < low {
		return low
	}
	if x > high {
		return high
	}
	return x
}

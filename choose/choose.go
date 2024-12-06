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

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/timeout"
	"github.com/charmbracelet/lipgloss"
)

func defaultKeymap() keymap {
	return keymap{
		Down: key.NewBinding(
			key.WithKeys("down", "j", "ctrl+j", "ctrl+n"),
			key.WithHelp("↓", "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k", "ctrl+k", "ctrl+p"),
			key.WithHelp("↑", "up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "ctrl+f"),
			key.WithHelp("→", "right"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h", "ctrl+b"),
			key.WithHelp("←", "left"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g", "home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G", "end"),
		),
		ToggleAll: key.NewBinding(
			key.WithKeys("a", "A", "ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
			key.WithDisabled(),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" ", "tab", "x", "ctrl+@"),
			key.WithHelp("x", "toggle"),
			key.WithDisabled(),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter", "ctrl+q"),
			key.WithHelp("enter", "submit"),
		),
	}
}

type keymap struct {
	Down,
	Up,
	Right,
	Left,
	Home,
	End,
	ToggleAll,
	Toggle,
	Abort,
	Submit key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Toggle,
		key.NewBinding(
			key.WithKeys("up", "down", "right", "left"),
			key.WithHelp("↑↓←→", "navigate"),
		),
		k.Submit,
		k.ToggleAll,
	}
}

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
	showHelp         bool
	help             help.Model
	keymap           keymap

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
		km := m.keymap
		switch {
		case key.Matches(msg, km.Down):
			m.index++
			if m.index >= len(m.items) {
				m.index = 0
				m.paginator.Page = 0
			}
			if m.index >= end {
				m.paginator.NextPage()
			}
		case key.Matches(msg, km.Up):
			m.index--
			if m.index < 0 {
				m.index = len(m.items) - 1
				m.paginator.Page = m.paginator.TotalPages - 1
			}
			if m.index < start {
				m.paginator.PrevPage()
			}
		case key.Matches(msg, km.Right):
			m.index = clamp(m.index+m.height, 0, len(m.items)-1)
			m.paginator.NextPage()
		case key.Matches(msg, km.Left):
			m.index = clamp(m.index-m.height, 0, len(m.items)-1)
			m.paginator.PrevPage()
		case key.Matches(msg, km.End):
			m.index = len(m.items) - 1
			m.paginator.Page = m.paginator.TotalPages - 1
		case key.Matches(msg, km.Home):
			m.index = 0
			m.paginator.Page = 0
		case key.Matches(msg, km.ToggleAll):
			if m.limit <= 1 {
				break
			}
			if m.numSelected < len(m.items) && m.numSelected < m.limit {
				m = m.selectAll()
			} else {
				m = m.deselectAll()
			}
		case key.Matches(msg, km.Abort):
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Toggle):
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
		case key.Matches(msg, km.Submit):
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

func (m model) selectAll() model {
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
	return m
}

func (m model) deselectAll() model {
	for i := range m.items {
		m.items[i].selected = false
		m.items[i].order = 0
	}
	m.numSelected = 0
	m.currentOrder = 0
	return m
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

	var parts []string

	if m.header != "" {
		parts = append(parts, m.headerStyle.Render(m.header))
	}
	parts = append(parts, s.String())
	if m.showHelp {
		parts = append(parts, m.help.View(m.keymap))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
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

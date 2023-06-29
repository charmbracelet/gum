// Package filter provides a fuzzy searching text input to allow filtering a
// list of options to select one option.
//
// By default it will list all the files (recursively) in the current directory
// for the user to choose one, but the script (or user) can provide different
// new-line separated options to choose from.
//
// I.e. let's pick from a list of gum flavors:
//
// $ cat flavors.text | gum filter
package filter

import (
	"strings"
	"time"

	"github.com/charmbracelet/gum/timeout"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type model struct {
	textinput             textinput.Model
	viewport              *viewport.Model
	choices               []string
	matches               []fuzzy.Match
	cursor                int
	header                string
	selected              map[string]struct{}
	limit                 int
	numSelected           int
	indicator             string
	selectedPrefix        string
	unselectedPrefix      string
	height                int
	aborted               bool
	quitting              bool
	headerStyle           lipgloss.Style
	matchStyle            lipgloss.Style
	textStyle             lipgloss.Style
	cursorTextStyle       lipgloss.Style
	indicatorStyle        lipgloss.Style
	selectedPrefixStyle   lipgloss.Style
	unselectedPrefixStyle lipgloss.Style
	reverse               bool
	fuzzy                 bool
	sort                  bool
	timeout               time.Duration
	hasTimeout            bool
}

func (m model) Init() tea.Cmd {
	return timeout.Init(m.timeout, nil)
}
func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder
	var lineTextStyle lipgloss.Style

	// For reverse layout, if the number of matches is less than the viewport
	// height, we need to offset the matches so that the first match is at the
	// bottom edge of the viewport instead of in the middle.
	if m.reverse && len(m.matches) < m.viewport.Height {
		s.WriteString(strings.Repeat("\n", m.viewport.Height-len(m.matches)))
	}

	// Since there are matches, display them so that the user can see, in real
	// time, what they are searching for.
	last := len(m.matches) - 1
	for i := range m.matches {
		// For reverse layout, the matches are displayed in reverse order.
		if m.reverse {
			i = last - i
		}
		match := m.matches[i]

		// If this is the current selected index, we add a small indicator to
		// represent it. Otherwise, simply pad the string.
		// The line's text style is set depending on whether or not the cursor
		// points to this line.
		if i == m.cursor {
			s.WriteString(m.indicatorStyle.Render(m.indicator))
			lineTextStyle = m.cursorTextStyle
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(m.indicator)))
			lineTextStyle = m.textStyle
		}

		// If there are multiple selections mark them, otherwise leave an empty space
		if _, ok := m.selected[match.Str]; ok {
			s.WriteString(m.selectedPrefixStyle.Render(m.selectedPrefix))
		} else if m.limit > 1 {
			s.WriteString(m.unselectedPrefixStyle.Render(m.unselectedPrefix))
		} else {
			s.WriteString(" ")
		}

		// For this match, there are a certain number of characters that have
		// caused the match. i.e. fuzzy matching.
		// We should indicate to the users which characters are being matched.
		mi := 0
		var buf strings.Builder
		for ci, c := range match.Str {
			// Check if the current character index matches the current matched
			// index. If so, color the character to indicate a match.
			if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
				// Flush text buffer.
				s.WriteString(lineTextStyle.Render(buf.String()))
				buf.Reset()

				s.WriteString(m.matchStyle.Render(string(c)))
				// We have matched this character, so we never have to check it
				// again. Move on to the next match.
				mi++
			} else {
				// Not a match, buffer a regular character.
				buf.WriteRune(c)
			}
		}
		// Flush text buffer.
		s.WriteString(lineTextStyle.Render(buf.String()))

		// We have finished displaying the match with all of it's matched
		// characters highlighted and the rest filled in.
		// Move on to the next match.
		s.WriteRune('\n')
	}

	m.viewport.SetContent(s.String())

	// View the input and the filtered choices
	header := m.headerStyle.Render(m.header)
	if m.reverse {
		view := m.viewport.View() + "\n" + m.textinput.View()
		if m.header != "" {
			return lipgloss.JoinVertical(lipgloss.Left, view, header)
		}

		return view
	}

	view := m.textinput.View() + "\n" + m.viewport.View()
	if m.header != "" {
		return lipgloss.JoinVertical(lipgloss.Left, header, view)
	}
	return view
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.quitting = true
			m.aborted = true
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)

	case tea.WindowSizeMsg:
		if m.height == 0 || m.height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
		}

		// Make place in the view port if header is set
		if m.header != "" {
			m.viewport.Height = m.viewport.Height - lipgloss.Height(m.headerStyle.Render(m.header))
		}
		m.viewport.Width = msg.Width
		if m.reverse {
			m.viewport.YOffset = clamp(0, len(m.matches), len(m.matches)-m.viewport.Height)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+n", "ctrl+j", "down":
			m.CursorDown()
		case "ctrl+p", "ctrl+k", "up":
			m.CursorUp()
		case "tab":
			if m.limit == 1 {
				break // no op
			}
			m.ToggleSelection()
			m.CursorDown()
		case "shift+tab":
			if m.limit == 1 {
				break // no op
			}
			m.ToggleSelection()
			m.CursorUp()
		case "ctrl+@":
			if m.limit == 1 {
				break // no op
			}
			m.ToggleSelection()
		default:
			m.textinput, cmd = m.textinput.Update(msg)

			// yOffsetFromBottom is the number of lines from the bottom of the
			// list to the top of the viewport. This is used to keep the viewport
			// at a constant position when the number of matches are reduced
			// in the reverse layout.
			var yOffsetFromBottom int
			if m.reverse {
				yOffsetFromBottom = max(0, len(m.matches)-m.viewport.YOffset)
			}

			// A character was entered, this likely means that the text input has
			// changed. This suggests that the matches are outdated, so update them.
			if m.fuzzy {
				if m.sort {
					m.matches = fuzzy.Find(m.textinput.Value(), m.choices)
				} else {
					m.matches = fuzzy.FindNoSort(m.textinput.Value(), m.choices)
				}
			} else {
				m.matches = exactMatches(m.textinput.Value(), m.choices)
			}

			// If the search field is empty, let's not display the matches
			// (none), but rather display all possible choices.
			if m.textinput.Value() == "" {
				m.matches = matchAll(m.choices)
			}

			// For reverse layout, we need to offset the viewport so that the
			// it remains at a constant position relative to the cursor.
			if m.reverse {
				maxYOffset := max(0, len(m.matches)-m.viewport.Height)
				m.viewport.YOffset = clamp(0, maxYOffset, len(m.matches)-yOffsetFromBottom)
			}
		}
	}

	// It's possible that filtering items have caused fewer matches. So, ensure
	// that the selected index is within the bounds of the number of matches.
	m.cursor = clamp(0, len(m.matches)-1, m.cursor)
	return m, cmd
}

func (m *model) CursorUp() {
	if m.reverse {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
		if len(m.matches)-m.cursor <= m.viewport.YOffset {
			m.viewport.SetYOffset(len(m.matches) - m.cursor - 1)
		}
	} else {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
		if m.cursor < m.viewport.YOffset {
			m.viewport.SetYOffset(m.cursor)
		}
	}
}

func (m *model) CursorDown() {
	if m.reverse {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
		if len(m.matches)-m.cursor > m.viewport.Height+m.viewport.YOffset {
			m.viewport.LineDown(1)
		}
	} else {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
		if m.cursor >= m.viewport.YOffset+m.viewport.Height {
			m.viewport.LineDown(1)
		}
	}
}

func (m *model) ToggleSelection() {
	if _, ok := m.selected[m.matches[m.cursor].Str]; ok {
		delete(m.selected, m.matches[m.cursor].Str)
		m.numSelected--
	} else if m.numSelected < m.limit {
		m.selected[m.matches[m.cursor].Str] = struct{}{}
		m.numSelected++
	}
}

func matchAll(options []string) []fuzzy.Match {
	matches := make([]fuzzy.Match, len(options))
	for i, option := range options {
		matches[i] = fuzzy.Match{Str: option}
	}
	return matches
}

func exactMatches(search string, choices []string) []fuzzy.Match {
	matches := fuzzy.Matches{}
	for i, choice := range choices {
		search = strings.ToLower(search)
		matchedString := strings.ToLower(choice)

		index := strings.Index(matchedString, search)
		if index >= 0 {
			matchedIndexes := []int{}
			for s := range search {
				matchedIndexes = append(matchedIndexes, index+s)
			}
			matches = append(matches, fuzzy.Match{
				Str:            choice,
				Index:          i,
				MatchedIndexes: matchedIndexes,
			})
		}
	}

	return matches
}

//nolint:unparam
func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

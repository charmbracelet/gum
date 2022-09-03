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

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/sahilm/fuzzy"
)

type model struct {
	textinput             textinput.Model
	viewport              *viewport.Model
	choices               []string
	matches               []fuzzy.Match
	cursor                int
	selected              map[string]struct{}
	limit                 int
	numSelected           int
	indicator             string
	selectedPrefix        string
	unselectedPrefix      string
	height                int
	aborted               bool
	quitting              bool
	matchStyle            lipgloss.Style
	textStyle             lipgloss.Style
	indicatorStyle        lipgloss.Style
	selectedPrefixStyle   lipgloss.Style
	unselectedPrefixStyle lipgloss.Style
}

func (m model) Init() tea.Cmd { return nil }
func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Since there are matches, display them so that the user can see, in real
	// time, what they are searching for.
	for i, match := range m.matches {
		// If this is the current selected index, we add a small indicator to
		// represent it. Otherwise, simply pad the string.
		if i == m.cursor {
			s.WriteString(m.indicatorStyle.Render(m.indicator))
		} else {
			s.WriteString(strings.Repeat(" ", runewidth.StringWidth(m.indicator)))
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
		var mi = 0
		for ci, c := range match.Str {
			// Check if the current character index matches the current matched
			// index. If so, color the character to indicate a match.
			if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
				s.WriteString(m.matchStyle.Render(string(c)))
				// We have matched this character, so we never have to check it
				// again. Move on to the next match.
				mi++
			} else {
				// Not a match, simply show the character, unstyled.
				s.WriteString(m.textStyle.Render(string(c)))
			}
		}

		// We have finished displaying the match with all of it's matched
		// characters highlighted and the rest filled in.
		// Move on to the next match.
		s.WriteRune('\n')
	}

	m.viewport.SetContent(s.String())

	// View the input and the filtered choices
	return m.textinput.View() + "\n" + m.viewport.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.height == 0 || m.height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
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
			m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
			if m.cursor >= m.viewport.YOffset+m.viewport.Height {
				m.viewport.LineDown(1)
			}
		case "ctrl+p", "ctrl+k", "up":
			m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
			if m.cursor < m.viewport.YOffset {
				m.viewport.SetYOffset(m.cursor)
			}
		case "tab":
			if m.limit == 1 {
				break // no op
			}

			// Tab is used to toggle selection of current item in the list
			if _, ok := m.selected[m.matches[m.cursor].Str]; ok {
				delete(m.selected, m.matches[m.cursor].Str)
				m.numSelected--
			} else if m.numSelected < m.limit {
				m.selected[m.matches[m.cursor].Str] = struct{}{}
				m.numSelected++
			}

			// Go down by one line
			m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
			if m.cursor >= m.viewport.YOffset+m.viewport.Height {
				m.viewport.LineDown(1)
			}
		default:
			m.textinput, cmd = m.textinput.Update(msg)

			// A character was entered, this likely means that the text input
			// has changed. This suggests that the matches are outdated, so
			// update them, with a fuzzy finding algorithm provided by
			// https://github.com/sahilm/fuzzy
			m.matches = fuzzy.Find(m.textinput.Value(), m.choices)

			// If the search field is empty, let's not display the matches
			// (none), but rather display all possible choices.
			if m.textinput.Value() == "" {
				m.matches = matchAll(m.choices)
			}
		}
	}

	// It's possible that filtering items have caused fewer matches. So, ensure
	// that the selected index is within the bounds of the number of matches.
	m.cursor = clamp(0, len(m.matches)-1, m.cursor)
	return m, cmd
}

func matchAll(options []string) []fuzzy.Match {
	var matches = make([]fuzzy.Match, len(options))
	for i, option := range options {
		matches[i] = fuzzy.Match{Str: option}
	}
	return matches
}

func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

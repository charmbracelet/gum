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

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/ordered"
	"github.com/rivo/uniseg"
	"github.com/sahilm/fuzzy"
)

func defaultKeymap() keymap {
	return keymap{
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+j", "ctrl+n"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+k", "ctrl+p"),
		),
		NDown: key.NewBinding(
			key.WithKeys("j"),
		),
		NUp: key.NewBinding(
			key.WithKeys("k"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
		),
		ToggleAndNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle"),
			key.WithDisabled(),
		),
		ToggleAndPrevious: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "toggle"),
			key.WithDisabled(),
		),
		Toggle: key.NewBinding(
			key.WithKeys("ctrl+@"),
			key.WithHelp("ctrl+@", "toggle"),
			key.WithDisabled(),
		),
		ToggleAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
			key.WithDisabled(),
		),
		FocusInSearch: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		FocusOutSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "blur search"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter", "ctrl+q"),
			key.WithHelp("enter", "submit"),
		),
	}
}

type keymap struct {
	FocusInSearch,
	FocusOutSearch,
	Down,
	Up,
	NDown,
	NUp,
	Home,
	End,
	ToggleAndNext,
	ToggleAndPrevious,
	ToggleAll,
	Toggle,
	Abort,
	Quit,
	Submit key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↓↑", "navigate"),
		),
		k.FocusInSearch,
		k.FocusOutSearch,
		k.ToggleAndNext,
		k.ToggleAll,
		k.Submit,
	}
}

type model struct {
	textinput             textinput.Model
	viewport              *viewport.Model
	choices               map[string]string
	filteringChoices      []string
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
	padding               []int
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
	showHelp              bool
	keymap                keymap
	help                  help.Model
	strict                bool
	submitted             bool
}

func (m model) Init() tea.Cmd { return textinput.Blink }

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

		styledOption := m.choices[match.Str]
		if len(match.MatchedIndexes) == 0 {
			// No matches, just render the text.
			s.WriteString(lineTextStyle.Render(styledOption))
			s.WriteRune('\n')
			continue
		}

		var ranges []lipgloss.Range
		for _, rng := range matchedRanges(match.MatchedIndexes) {
			// ansi.Cut is grapheme and ansi sequence aware, we match against a ansi.Stripped string, but we might still have graphemes.
			// all that to say that rng is byte positions, but we need to pass it down to ansi.Cut as char positions.
			// so we need to adjust it here:
			start, stop := bytePosToVisibleCharPos(match.Str, rng)
			ranges = append(ranges, lipgloss.NewRange(start, stop+1, m.matchStyle))
		}

		s.WriteString(lineTextStyle.Render(lipgloss.StyleRanges(styledOption, ranges...)))

		// We have finished displaying the match with all of it's matched
		// characters highlighted and the rest filled in.
		// Move on to the next match.
		s.WriteRune('\n')
	}

	m.viewport.SetContent(s.String())

	// View the input and the filtered choices
	header := m.headerStyle.Render(m.header)
	if m.reverse {
		view := m.viewport.View()
		if m.header != "" {
			view += "\n" + header
		}
		view += "\n" + m.textinput.View()
		if m.showHelp {
			view += m.helpView()
		}
		return lipgloss.NewStyle().
			Padding(m.padding...).
			Render(view)
	}

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

func (m model) helpView() string {
	return "\n\n" + m.help.View(m.keymap)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, icmd tea.Cmd
	m.textinput, icmd = m.textinput.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.height == 0 || m.height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
		}
		// Include the header in the height calculation.
		if m.header != "" {
			m.viewport.Height = m.viewport.Height - lipgloss.Height(m.headerStyle.Render(m.header))
		}
		// Include the help in the total height calculation.
		if m.showHelp {
			m.viewport.Height = m.viewport.Height - lipgloss.Height(m.helpView())
		}
		m.viewport.Height = m.viewport.Height - m.padding[0] - m.padding[2]
		m.viewport.Width = msg.Width - m.padding[1] - m.padding[3]
		m.textinput.Width = msg.Width - m.padding[1] - m.padding[3]
		if m.reverse {
			m.viewport.YOffset = ordered.Clamp(len(m.matches)-m.viewport.Height, 0, len(m.matches))
		}
	case tea.KeyMsg:
		km := m.keymap
		switch {
		case key.Matches(msg, km.FocusInSearch):
			m.textinput.Focus()
		case key.Matches(msg, km.FocusOutSearch):
			m.textinput.Blur()
		case key.Matches(msg, km.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Abort):
			m.quitting = true
			return m, tea.Interrupt
		case key.Matches(msg, km.Submit):
			m.quitting = true
			m.submitted = true
			return m, tea.Quit
		case key.Matches(msg, km.Down, km.NDown):
			m.CursorDown()
		case key.Matches(msg, km.Up, km.NUp):
			m.CursorUp()
		case key.Matches(msg, km.Home):
			m.cursor = 0
			m.viewport.GotoTop()
		case key.Matches(msg, km.End):
			m.cursor = len(m.choices) - 1
			m.viewport.GotoBottom()
		case key.Matches(msg, km.ToggleAndNext):
			if m.limit == 1 {
				break // no op
			}
			m.ToggleSelection()
			m.CursorDown()
		case key.Matches(msg, km.ToggleAndPrevious):
			if m.limit == 1 {
				break // no op
			}
			m.ToggleSelection()
			m.CursorUp()
		case key.Matches(msg, km.Toggle):
			if m.limit == 1 {
				break // no op
			}
			m.ToggleSelection()
		case key.Matches(msg, km.ToggleAll):
			if m.limit <= 1 {
				break
			}
			if m.numSelected < len(m.matches) && m.numSelected < m.limit {
				m = m.selectAll()
			} else {
				m = m.deselectAll()
			}
		default:
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
			var choices []string
			if !m.strict {
				choices = append(choices, m.textinput.Value())
			}
			choices = append(choices, m.filteringChoices...)
			if m.fuzzy {
				if m.sort {
					m.matches = fuzzy.Find(m.textinput.Value(), choices)
				} else {
					m.matches = fuzzy.FindNoSort(m.textinput.Value(), choices)
				}
			} else {
				m.matches = exactMatches(m.textinput.Value(), choices)
			}

			// If the search field is empty, let's not display the matches
			// (none), but rather display all possible choices.
			if m.textinput.Value() == "" {
				m.matches = matchAll(m.filteringChoices)
			}

			// For reverse layout, we need to offset the viewport so that the
			// it remains at a constant position relative to the cursor.
			if m.reverse {
				maxYOffset := max(0, len(m.matches)-m.viewport.Height)
				m.viewport.YOffset = ordered.Clamp(len(m.matches)-yOffsetFromBottom, 0, maxYOffset)
			}
		}
	}

	m.keymap.FocusInSearch.SetEnabled(!m.textinput.Focused())
	m.keymap.FocusOutSearch.SetEnabled(m.textinput.Focused())
	m.keymap.NUp.SetEnabled(!m.textinput.Focused())
	m.keymap.NDown.SetEnabled(!m.textinput.Focused())
	m.keymap.Home.SetEnabled(!m.textinput.Focused())
	m.keymap.End.SetEnabled(!m.textinput.Focused())

	// It's possible that filtering items have caused fewer matches. So, ensure
	// that the selected index is within the bounds of the number of matches.
	m.cursor = ordered.Clamp(m.cursor, 0, len(m.matches)-1)
	return m, tea.Batch(cmd, icmd)
}

func (m *model) CursorUp() {
	if len(m.matches) == 0 {
		return
	}
	if m.reverse { //nolint:nestif
		m.cursor = (m.cursor + 1) % len(m.matches)
		if len(m.matches)-m.cursor <= m.viewport.YOffset {
			m.viewport.ScrollUp(1)
		}
		if len(m.matches)-m.cursor > m.viewport.Height+m.viewport.YOffset {
			m.viewport.SetYOffset(len(m.matches) - m.viewport.Height)
		}
	} else {
		m.cursor = (m.cursor - 1 + len(m.matches)) % len(m.matches)
		if m.cursor < m.viewport.YOffset {
			m.viewport.ScrollUp(1)
		}
		if m.cursor >= m.viewport.YOffset+m.viewport.Height {
			m.viewport.SetYOffset(len(m.matches) - m.viewport.Height)
		}
	}
}

func (m *model) CursorDown() {
	if len(m.matches) == 0 {
		return
	}
	if m.reverse { //nolint:nestif
		m.cursor = (m.cursor - 1 + len(m.matches)) % len(m.matches)
		if len(m.matches)-m.cursor > m.viewport.Height+m.viewport.YOffset {
			m.viewport.ScrollDown(1)
		}
		if len(m.matches)-m.cursor <= m.viewport.YOffset {
			m.viewport.GotoTop()
		}
	} else {
		m.cursor = (m.cursor + 1) % len(m.matches)
		if m.cursor >= m.viewport.YOffset+m.viewport.Height {
			m.viewport.ScrollDown(1)
		}
		if m.cursor < m.viewport.YOffset {
			m.viewport.GotoTop()
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

func (m model) selectAll() model {
	for i := range m.matches {
		if m.numSelected >= m.limit {
			break // do not exceed given limit
		}
		if _, ok := m.selected[m.matches[i].Str]; ok {
			continue
		}
		m.selected[m.matches[i].Str] = struct{}{}
		m.numSelected++
	}
	return m
}

func (m model) deselectAll() model {
	m.selected = make(map[string]struct{})
	m.numSelected = 0
	return m
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

func matchedRanges(in []int) [][2]int {
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

func bytePosToVisibleCharPos(str string, rng [2]int) (int, int) {
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

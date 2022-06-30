package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/seashell/internal/stdin"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/coral"
	"github.com/sahilm/fuzzy"
)

var defaultIgnore = []string{
	".git",
}

type model struct {
	textinput textinput.Model
	choices   []string
	matches   []fuzzy.Match
	selected  int
	indicator string
	height    int
}

func (m model) Init() tea.Cmd { return nil }
func (m model) View() string {
	var s strings.Builder

	// Since there are matches, display them so that the user can see, in real
	// time, what they are searching for.
	for i, match := range m.matches {

		// If this is the current selected index, we add a small indicator to
		// represent it. Otherwise, simply pad the string.
		if i == m.selected {
			s.WriteString(m.textinput.PromptStyle.Render(m.indicator) + " ")
		} else {
			s.WriteString(strings.Repeat(" ", runewidth.StringWidth(m.indicator)) + " ")
		}

		// For this match, there are a certain number of characters that have
		// caused the match. i.e. fuzzy matching.
		// We should indicate to the users which characters are being matched.
		var mi = 0
		for ci, c := range match.Str {
			// Check if the current character index matches the current matched
			// index. If so, color the character to indicate a match.
			if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
				s.WriteString(m.textinput.PromptStyle.Render(string(c)))
				// We have matched this character, so we never have to check it
				// again. Move on to the next match.
				mi++
			} else {
				// Not a match, simply show the character, unstyled.
				s.WriteRune(c)
			}
		}

		// We have finished displaying the match with all of it's matched
		// characters highlighted and the rest filled in.
		// Move on to the next match.
		s.WriteRune('\n')
	}

	tv := m.textinput.View()
	results := lipgloss.NewStyle().MaxHeight(m.height - lipgloss.Height(tv)).Render(s.String())
	// View the input and the filtered choices
	return tv + "\n" + results
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "enter":
			return m, tea.Quit
		case "ctrl+n":
			m.selected = clamp(0, len(m.matches)-1, m.selected+1)
		case "ctrl+p":
			m.selected = clamp(0, len(m.matches)-1, m.selected-1)
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
	m.selected = clamp(0, len(m.matches)-1, m.selected)
	return m, cmd
}

type options struct {
	prompt      *string
	placeholder *string
	width       *int
	accentColor *string
	indicator   *string
}

// Cmd returns the input command
func Cmd() *coral.Command {
	var opts options

	var cmd = &coral.Command{
		Use:   "search",
		Short: "Search searches for an item in a list of options. Uses fuzzy matching to filter.",
		RunE: func(cmd *coral.Command, args []string) error {
			ti := textinput.New()

			// Flags + Options
			ti.Prompt = *opts.prompt
			ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(*opts.accentColor))
			ti.Placeholder = *opts.placeholder
			ti.Width = *opts.width

			ti.Focus()

			input, err := stdin.Read()
			if err != nil {
				return err
			}
			if input == "" {
				// No input, let's assume they want to filter files in the
				// current directory.
				filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
					if info.IsDir() || shouldIgnore(path) {
						return nil
					}

					input += path + "\n"
					return nil
				},
				)
			}
			choices := strings.Split(string(input), "\n")

			p := tea.NewProgram(model{
				textinput: ti,
				choices:   choices,
				matches:   matchAll(choices),
				indicator: *opts.indicator,
			}, tea.WithOutput(os.Stderr))

			m, err := p.StartReturningModel()
			mm := m.(model)

			if len(mm.matches) > mm.selected && mm.selected >= 0 {
				fmt.Println(mm.matches[mm.selected].Str)
			}

			return err
		},
	}

	opts = options{
		prompt:      cmd.Flags().String("prompt", "> ", "Prompt to display"),
		placeholder: cmd.Flags().String("placeholder", "Enter a value...", "Placeholder value"),
		width:       cmd.Flags().Int("width", 20, "Input width"),
		accentColor: cmd.Flags().String("accent-color", "#FF06B7", "Accent color for the prompt, indicator, and matches"),
		indicator:   cmd.Flags().String("indicator", "•", "Character to use to indicate the selected match"),
	}

	return cmd
}

func matchAll(options []string) []fuzzy.Match {
	var matches []fuzzy.Match
	for _, option := range options {
		matches = append(matches, fuzzy.Match{Str: option})
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

func shouldIgnore(e string) bool {
	for _, a := range defaultIgnore {
		if strings.HasPrefix(e, a) {
			return true
		}
	}
	return false
}

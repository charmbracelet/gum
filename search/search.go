package search

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/coral"
	"github.com/sahilm/fuzzy"
)

const (
	indicator = "•"
)

type model struct {
	textinput textinput.Model
	options   []string
	matches   []fuzzy.Match
	selected  int
}

func (m model) Init() tea.Cmd { return nil }
func (m model) View() string {
	var s string
	for i, match := range m.matches {
		if i == m.selected {
			s += m.textinput.PromptStyle.Render(indicator)
		} else {
			s += " "
		}
		s += " "
		for ci, c := range match.Str {
			included := false
			for _, mi := range match.MatchedIndexes {
				if ci == mi {
					included = true
					s += m.textinput.PromptStyle.Render(string(c))
					break
				}
			}
			if !included {
				s += string(c)
			}
		}
		s += "\n"
	}

	return m.textinput.View() + "\n" + s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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
			m.matches = fuzzy.Find(m.textinput.Value(), m.options)
			if m.textinput.Value() == "" {
				m.matches = matchAll(m.options)
			}
		}
	}

	m.selected = clamp(0, len(m.matches)-1, m.selected)
	return m, cmd
}

type options struct {
	prompt      *string
	placeholder *string
	width       *int
	accentColor *string
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

			ti.SetCursorMode(textinput.CursorStatic)
			ti.Focus()

			input, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			options := strings.Split(string(input), "\n")

			p := tea.NewProgram(model{
				textinput: ti,
				options:   options,
				matches:   matchAll(options),
			})

			m, err := p.StartReturningModel()
			mm := m.(model)

			if len(mm.matches) > mm.selected {
				fmt.Println(mm.matches[mm.selected].Str)
			}

			return err
		},
	}

	opts = options{
		prompt:      cmd.Flags().String("prompt", "> ", "Prompt to display"),
		placeholder: cmd.Flags().String("placeholder", "Enter a value...", "Placeholder value"),
		width:       cmd.Flags().Int("width", 20, "Input width"),
		accentColor: cmd.Flags().String("color", "#FF06B7", "Accent color for the prompt, indicator, and matches"),
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

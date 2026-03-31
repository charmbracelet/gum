package file

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type fuzzyModel struct {
	basePath     string
	allFiles     []string
	matches      []fuzzy.Match
	cursor       int
	textinput    textinput.Model
	viewport     *viewport.Model
	matchStyle   lipgloss.Style
	fileStyle    lipgloss.Style
	dirStyle     lipgloss.Style
	selectedPath string
	quitting     bool
	submitted    bool
	padding      []int
	height       int
	header       string
	showHelp     bool
	help         help.Model
}

func (m fuzzyModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m fuzzyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textinput.Width = msg.Width - 4
		h := msg.Height - 3
		if m.header != "" {
			h--
		}
		if m.showHelp {
			h -= 2
		}
		if h < 1 {
			h = 1
		}
		m.viewport.Height = h
		m.viewport.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.quitting = true
			return m, tea.Interrupt
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if len(m.matches) > 0 && m.cursor < len(m.matches) {
				m.selectedPath = filepath.Join(m.basePath, m.matches[m.cursor].Str)
				m.submitted = true
				m.quitting = true
			}
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "ctrl+p"))):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "ctrl+n"))):
			if m.cursor < len(m.matches)-1 {
				m.cursor++
			}
			return m, nil
		}
	}

	prevValue := m.textinput.Value()
	m.textinput, cmd = m.textinput.Update(msg)

	if m.textinput.Value() != prevValue {
		query := m.textinput.Value()
		if query == "" {
			m.matches = matchAll(m.allFiles)
		} else {
			m.matches = fuzzy.Find(query, m.allFiles)
		}
		m.cursor = 0
	}

	return m, cmd
}

func (m fuzzyModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	if m.header != "" {
		b.WriteString(m.header + "\n")
	}

	b.WriteString(m.textinput.View() + "\n")

	maxVisible := m.height
	if maxVisible == 0 {
		maxVisible = 20
	}

	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}

	end := start + maxVisible
	if end > len(m.matches) {
		end = len(m.matches)
	}

	for i := start; i < end; i++ {
		match := m.matches[i]
		line := match.Str
		indicator := "  "

		if i == m.cursor {
			indicator = "> "
			line = m.renderMatch(match)
		} else if len(match.MatchedIndexes) > 0 {
			line = m.renderMatch(match)
		}

		if i == m.cursor {
			b.WriteString(m.matchStyle.Render(indicator + line))
		} else {
			b.WriteString(indicator + line)
		}
		b.WriteString("\n")
	}

	if len(m.matches) == 0 {
		b.WriteString("  No matches\n")
	}

	info := fmt.Sprintf("  %d/%d", len(m.matches), len(m.allFiles))
	b.WriteString(lipgloss.NewStyle().Faint(true).Render(info))

	if m.showHelp {
		b.WriteString("\n" + m.help.View(fuzzyKeymap{}))
	}

	return lipgloss.NewStyle().Padding(m.padding...).Render(b.String())
}

func (m fuzzyModel) renderMatch(match fuzzy.Match) string {
	text := match.Str
	if len(match.MatchedIndexes) == 0 {
		return text
	}

	highlighted := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	var result strings.Builder
	matchSet := make(map[int]bool)
	for _, idx := range match.MatchedIndexes {
		matchSet[idx] = true
	}
	for i, ch := range text {
		if matchSet[i] {
			result.WriteString(highlighted.Render(string(ch)))
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

type fuzzyKeymap struct{}

func (fuzzyKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("↑↓"), key.WithHelp("↑↓", "navigate")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

func (fuzzyKeymap) FullHelp() [][]key.Binding { return nil }

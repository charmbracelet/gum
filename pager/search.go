package pager

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/gum/internal/utils"
)

type search struct {
	active       bool
	input        textinput.Model
	query        *regexp.Regexp
	lastMatchLoc int
}

func (s *search) new() {
	input := textinput.New()
	input.Placeholder = "search"
	input.Prompt = "/"
	s.input = input
}

func (s *search) Begin() {
	s.new()
	s.active = true
	s.input.Focus()
}

// Execute find all lines in the model with a match
func (s *search) Execute(m *model) {
	defer s.Done()
	if s.input.Value() == "" {
		s.query = nil
		return
	}

	s.query = regexp.MustCompile(s.input.Value())
	matches := utils.UniqueStrings(s.query.FindAllString(m.content, -1))
	for _, match := range matches {
		replacement := m.matchStyle.Render(match)
		m.content = strings.ReplaceAll(m.content, match, replacement)
	}
}

func (s *search) Done() {
	s.active = false
}

func (s *search) NextMatch(m *model) {
	// Check that we are within bounds.
	if s.query == nil {
		return
	}

	// Find the string to highlight.
	nextMatch := s.query.FindString(m.content[s.lastMatchLoc:])
	if nextMatch == "" {
		// Start the search from the beginning of the document.
		s.lastMatchLoc = 0
		m.viewport.SetYOffset(0)
		return
	}
	m.content = m.content[:s.lastMatchLoc] + strings.Replace(m.content[s.lastMatchLoc:], nextMatch, m.matchHighlightStyle.Render(nextMatch), 1)

	// Update the postion of the last found match.
	nextMatchI := s.query.FindStringIndex(m.content[s.lastMatchLoc:])
	s.lastMatchLoc += nextMatchI[1]

	// Update the viewport position.
	line := 0
	for i, c := range m.content {
		if c == '\n' {
			line++
		}
		if i == s.lastMatchLoc {
			break
		}
	}

	// Only update if the match is not within the viewport
	if line > m.viewport.YOffset+m.viewport.VisibleLineCount()-1 {
		m.viewport.SetYOffset(line)
	}
}

func (s *search) PrevMatch(m *model) {
}

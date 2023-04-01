package pager

import (
	"fmt"
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
	prevMatch    string
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
	query := regexp.MustCompile(fmt.Sprintf("(%s)", s.query.String()))
	m.content = query.ReplaceAllString(m.content, m.matchStyle.Render("$1"))
}

func (s *search) Done() {
	s.active = false
	s.lastMatchLoc = 0
	s.prevMatch = ""
}

func (s *search) NextMatch(m *model) {
	// Check that we are within bounds.
	if s.query == nil {
		return
	}

	// Removed last highlight.
	if s.prevMatch != "" {
		leftPadding, rightPadding := utils.LipglossLengthPadding(s.prevMatch, m.matchHighlightStyle)
		metastring := regexp.QuoteMeta(m.matchHighlightStyle.Render(s.query.String()))
		query := regexp.MustCompile("(" + metastring[:leftPadding+1] + ")(" + s.query.String() + ")(" + metastring[len(metastring)-rightPadding-1:] + ")")
		m.content = query.ReplaceAllString(m.content, "$2")
	}
	// Find the string to highlight.
	nextMatch := s.query.FindString(m.content[s.lastMatchLoc:])
	s.prevMatch = nextMatch
	if nextMatch == "" {
		// Start the search from the beginning of the document.
		s.lastMatchLoc = 0
		m.viewport.GotoTop()
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
	if line > m.viewport.YOffset+m.viewport.VisibleLineCount()-1 || line < m.viewport.YOffset {
		m.viewport.SetYOffset(line)
	}
}

func (s *search) PrevMatch(m *model) {
	// Check that we are within bounds.
	if s.query == nil {
		return
	}

	// Removed last highlight.
	if s.prevMatch != "" {
		leftPadding, rightPadding := utils.LipglossLengthPadding(s.prevMatch, m.matchHighlightStyle)
		metastring := regexp.QuoteMeta(m.matchHighlightStyle.Render(s.query.String()))
		query := regexp.MustCompile("(" + metastring[:leftPadding+1] + ")(" + s.query.String() + ")(" + metastring[len(metastring)-rightPadding-1:] + ")")
		m.content = query.ReplaceAllString(m.content, "$2")
	}
	// Find the string to highlight.
	var i int
	var nextMatch string
	for i = 0; s.query.FindString(m.content[i:s.lastMatchLoc]) != ""; i++ {
		nextMatch = s.query.FindString(m.content[:s.lastMatchLoc])
	}
	s.prevMatch = nextMatch
	if nextMatch == "" {
		// Start the search from the beginning of the document.
		s.lastMatchLoc = m.viewport.TotalLineCount()
		m.viewport.GotoTop()
		return
	}
	// m.content = strings.Replace(m.content[:i], nextMatch, m.matchHighlightStyle.Render(nextMatch), 1) + m.content[i:]

	// Update the postion of the last found match.
	for i = 0; s.query.FindString(m.content[i:s.lastMatchLoc]) != ""; i++ {
	}

	s.lastMatchLoc = i

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
	if line > m.viewport.YOffset+m.viewport.VisibleLineCount()-1 || line < m.viewport.YOffset {
		m.viewport.SetYOffset(line)
	}
}

func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

package pager

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/gum/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

type search struct {
	active           bool
	input            textinput.Model
	query            *regexp.Regexp
	matchIndex       int
	matchLipglossStr string
	matchString      string
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

// Execute find all lines in the model with a match.
func (s *search) Execute(m *model) {
	defer s.Done()
	if s.input.Value() == "" {
		s.query = nil
		return
	}

	s.query = regexp.MustCompile(s.input.Value())
	query := regexp.MustCompile(fmt.Sprintf("(%s)", s.query.String()))
	m.content = query.ReplaceAllString(m.content, m.matchStyle.Render("$1"))

	// Recompile the regex to match the an replace the highlights.
	leftPad, _ := utils.LipglossPadding(m.matchStyle)
	matchingString := regexp.QuoteMeta(m.matchStyle.Render()[:leftPad]) + s.query.String() + regexp.QuoteMeta(m.matchStyle.Render()[leftPad:])
	s.query = regexp.MustCompile(matchingString)
}

func (s *search) Done() {
	s.active = false

	// To account for the first match is always executed.
	s.matchIndex = -1
}

func (s *search) NextMatch(m *model) {
	// Check that we are within bounds.
	if s.query == nil {
		return
	}

	// Remove Previvous highlight.
	m.content = strings.Replace(m.content, s.matchLipglossStr, s.matchString, 1)

	// Highlight the next match.
	allMatches := s.query.FindAllStringIndex(m.content, -1)
	if len(allMatches) == 0 {
		return
	}

	s.matchIndex = (s.matchIndex + 1) % len(allMatches)
	match := allMatches[s.matchIndex]
	lhs := m.content[:match[0]]
	rhs := m.content[match[0]:]
	s.matchString = m.content[match[0]:match[1]]
	s.matchLipglossStr = m.matchHighlightStyle.Render(s.matchString)
	m.content = lhs + strings.Replace(rhs, m.content[match[0]:match[1]], s.matchLipglossStr, 1)

	// Update the viewport position.
	line := 0
	for i, c := range softWrapEm(m.content, m.maxWidth, m.softWrap) {
		if c == '\n' {
			line++
		}
		if i == match[0]+len(s.matchLipglossStr) {
			break
		}
	}

	// Only update if the match is not within the viewport.
	if line > m.viewport.YOffset-1+m.viewport.VisibleLineCount()-1 || line < m.viewport.YOffset {
		m.viewport.SetYOffset(line)
	}
}

func (s *search) PrevMatch(m *model) {
	// Check that we are within bounds.
	if s.query == nil {
		return
	}

	// Remove Previvous highlight.
	m.content = strings.Replace(m.content, s.matchLipglossStr, s.matchString, 1)

	// Highlight the previous match.
	allMatches := s.query.FindAllStringIndex(m.content, -1)
	if len(allMatches) == 0 {
		return
	}

	s.matchIndex = (s.matchIndex - 1) % len(allMatches)
	if s.matchIndex < 0 {
		s.matchIndex = 0
	}

	match := allMatches[s.matchIndex]
	lhs := m.content[:match[0]]
	rhs := m.content[match[0]:]
	s.matchString = m.content[match[0]:match[1]]
	s.matchLipglossStr = m.matchHighlightStyle.Render(s.matchString)
	m.content = lhs + strings.Replace(rhs, m.content[match[0]:match[1]], s.matchLipglossStr, 1)

	// Update the viewport position.
	line := 0
	for i, c := range softWrapEm(m.content, m.maxWidth, m.softWrap) {
		if c == '\n' {
			line++
		}
		if i == match[0]+len(s.matchLipglossStr) {
			break
		}
	}

	// Only update if the match is not within the viewport.
	if line > m.viewport.YOffset-1+m.viewport.VisibleLineCount()-1 || line < m.viewport.YOffset {
		m.viewport.SetYOffset(line)
	}
}

func softWrapEm(str string, maxWidth int, softWrap bool) string {
	var text strings.Builder
	for _, line := range strings.Split(str, "\n") {
		for softWrap && lipgloss.Width(line) > maxWidth {
			truncatedLine := utils.LipglossTruncate(line, maxWidth)
			text.WriteString(truncatedLine)
			text.WriteString("\n")
			line = strings.Replace(line, truncatedLine, "", 1)
		}
		text.WriteString(utils.LipglossTruncate(line, maxWidth))
		text.WriteString("\n")
	}

	return text.String()
}

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

	leftPadding, _ := utils.LipglossLengthPadding(s.prevMatch, m.matchHighlightStyle)
	allMatches := s.query.FindAllSubmatchIndex([]byte(m.content), -1)
	for i, subm := range allMatches {
		if subm[1]+leftPadding == s.lastMatchLoc {
			// Highliht the current match.
			m.content = m.content[:allMatches[i+1][0]] + strings.Replace(m.content[allMatches[i+1][0]:], m.content[allMatches[i+1][0]:allMatches[i+1][1]], m.matchHighlightStyle.Render(m.content[allMatches[i+1][0]:allMatches[i+1][1]]), 1)
			s.lastMatchLoc = allMatches[i+1][1] + leftPadding
			break
		}
		if i == len(allMatches)-1 {
			m.content = m.content[:allMatches[0][0]] + strings.Replace(m.content[allMatches[0][0]:], m.content[allMatches[0][0]:allMatches[0][1]], m.matchHighlightStyle.Render(m.content[allMatches[0][0]:allMatches[0][1]]), 1)
			s.lastMatchLoc = allMatches[0][1] + leftPadding
			break
		}
	}

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

	var prev []int
	leftPadding, _ := utils.LipglossLengthPadding(s.prevMatch, m.matchHighlightStyle)
	allMatches := s.query.FindAllSubmatchIndex([]byte(m.content), -1)
	for i, subm := range allMatches {
		if prev != nil && subm[1]+leftPadding == s.lastMatchLoc {
			// Highliht the current match.
			m.content = m.content[:prev[0]] + strings.Replace(m.content[prev[0]:], m.content[prev[0]:prev[1]], m.matchHighlightStyle.Render(m.content[prev[0]:prev[1]]), 1)
			s.lastMatchLoc = prev[1] + leftPadding
			break
		}

		// If reaching this at the end of the loop we have looked through all matches.
		if i == len(allMatches)-1 {
			m.content = m.content[:subm[0]] + strings.Replace(m.content[subm[0]:], m.content[subm[0]:subm[1]], m.matchHighlightStyle.Render(m.content[subm[0]:subm[1]]), 1)
			s.lastMatchLoc = subm[1] + leftPadding
		}
		prev = subm
	}

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

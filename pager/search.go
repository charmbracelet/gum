package pager

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

type search struct {
	Active   bool
	Input    textinput.Model
	Matches  []int
	CurMatch int
}

func (s *search) New() {
	input := textinput.New()
	input.Placeholder = "search"
	input.Prompt = "/"
	s.Input = input
}

func (s *search) Begin() {
	s.New()
	s.Matches = s.Matches[0:0]
	s.Active = true
	s.Input.Focus()
}

func (s *search) Execute(m *model) {
	defer s.Done()
	if s.Input.Value() == "" {
		return
	}

	queryRe := regexp.MustCompile(s.Input.Value())
	for i, line := range strings.Split(m.content, "\n") {
		if queryRe.Match([]byte(line)) {
			s.Matches = append(s.Matches, i)
		}
	}

	matches := unique(queryRe.FindAllString(m.content, -1))
	for _, match := range matches {
		replacement := m.matchStyle.Render(match)
		m.content = strings.ReplaceAll(m.content, match, replacement)
	}
}

func unique(strings []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, s := range strings {
		if _, uniq := keys[s]; !uniq {
			keys[s] = true
			list = append(list, s)
		}
	}
	return list
}

func (s *search) Done() {
	s.Active = false
	s.CurMatch = 0
}

func (s *search) NextMatch(m *model) {
	switch {
	case len(s.Matches) <= 0:
		return
	case s.CurMatch == len(s.Matches)-1:
		(*m).viewport.GotoTop()
		s.CurMatch = 0
	case (*m).viewport.AtBottom():
		s.CurMatch++
	default:
		for i, match := range s.Matches {
			if match > (*m).viewport.YOffset {
				s.CurMatch = i
				break
			}
		}
	}

	m.viewport.SetYOffset(m.search.Matches[s.CurMatch])
}

func (s *search) PrevMatch(m *model) {
	switch {
	case len(s.Matches) <= 0:
		return
	case s.CurMatch == 0:
		(*m).viewport.GotoBottom()
		s.CurMatch = len(s.Matches) - 1
	case (*m).viewport.AtBottom():
		s.CurMatch--
	default:
		for i := len(s.Matches) - 1; i >= 0; i-- {
			if s.Matches[i] < (*m).viewport.YOffset {
				s.CurMatch = i
				break
			}
		}
	}

	m.viewport.SetYOffset(m.search.Matches[s.CurMatch])
}

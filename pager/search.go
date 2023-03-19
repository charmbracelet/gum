package pager

import (
	"github.com/charmbracelet/bubbles/textinput"
	"math"
	"regexp"
	"strings"
)

type Search struct {
	Active   bool
	Input    textinput.Model
	Matches  []int
	CurMatch int
}

func (s *Search) New() {
	input := textinput.New()
	input.Placeholder = "Search"
	input.Prompt = "/"
	s.Input = input
}

func (s *Search) Begin() {
	s.New()
	s.Matches = s.Matches[0:0]
	s.Active = true
	s.Input.Focus()
}

func (s *Search) Execute(m *model) {
	defer s.Done()
	if s.Input.Value() == "" {
		return
	}

	queryRe := regexp.MustCompile(s.Input.Value())
	//matches := queryRe.FindAllIndex([]byte(m.content), -1)
	for i, line := range strings.Split(m.content, "\n") {
		if queryRe.Match([]byte(line)) {
			s.Matches = append(s.Matches, i)
		}
	}
}

func (s *Search) Done() {
	s.Active = false
	s.CurMatch = 0

}

func (s *Search) NextMatch(m *model) {
	if len(s.Matches) <= 0 {
		return
	}

	pos := m.search.Matches[s.findNext(&m)] - m.viewport.YOffset
	m.viewport.LineDown(pos)
}

func (s *Search) PrevMatch(m *model) {
	if len(s.Matches) <= 0 {
		return
	}

	pos := m.viewport.YOffset - m.search.Matches[s.findPrev(&m)]
	m.viewport.LineUp(int(math.Abs(float64(pos))))
}

func (s *Search) findNext(m **model) int {
	if s.CurMatch == len(s.Matches)-1 {
		(*m).viewport.GotoTop()
		s.CurMatch = 0
		return 0
	}

	for i, match := range s.Matches {
		if match > (*m).viewport.YOffset {
			s.CurMatch = i
			return i
		}
	}

	return 0
}

func (s *Search) findPrev(m **model) int {
	if s.CurMatch == 0 {
		(*m).viewport.GotoBottom()
		s.CurMatch = len(s.Matches) - 1
		return s.CurMatch
	}

	for i := len(s.Matches) - 1; i >= 0; i-- {
		if s.Matches[i] < (*m).viewport.YOffset {
			s.CurMatch = i
			return i
		}
	}

	return len(s.Matches) - 1
}

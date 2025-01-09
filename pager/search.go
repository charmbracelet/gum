package pager

import (
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type search struct {
	visible    bool
	navigating bool
	input      textinput.Model
}

func (s *search) new() {
	input := textinput.New()
	input.Placeholder = "search"
	input.Prompt = "/"
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	s.input = input
}

func (s *search) Show(w int) {
	s.new()
	s.visible = true
	s.input.Width = w
	s.input.Focus()
}

// Execute find all lines in the model with a match.
func (s *search) Execute(content string) [][]int {
	if s.input.Value() == "" {
		s.navigating = false
		s.visible = false
		return nil
	}

	s.navigating = true
	s.visible = false
	query, err := regexp.Compile(s.input.Value())
	if err != nil {
		return nil
	}
	return query.FindAllStringIndex(content, -1)
}

func (s *search) Done() {
	s.visible = false
	s.navigating = false
}

package filter

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func newResizeTestModel(itemCount int) model {
	filteringChoices := make([]string, itemCount)
	choices := make(map[string]string, itemCount)
	for i := range filteringChoices {
		filteringChoices[i] = fmt.Sprintf("item-%02d", i)
		choices[filteringChoices[i]] = filteringChoices[i]
	}
	ti := textinput.New()
	ti.Focus()
	vp := viewport.New(80, 20)
	return model{
		textinput:        ti,
		viewport:         &vp,
		choices:          choices,
		filteringChoices: filteringChoices,
		matches:          matchAll(filteringChoices),
		selected:         make(map[string]struct{}),
		limit:            1,
		showHelp:         true,
		header:           "Pick one",
		keymap:           defaultKeymap(),
		help:             help.New(),
		indicator:        ">",
		padding:          []int{1, 1, 1, 1},
	}
}

func TestFilterResizeThenCursorUpNoPanic(t *testing.T) {
	m := newResizeTestModel(30)
	m.cursor = 0

	var updated tea.Model
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
	m = updated.(model)

	// Shrink terminal so viewport height goes negative after chrome subtraction.
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 6})
	m = updated.(model)
	if m.viewport.Height < 0 {
		t.Fatalf("viewport height must be non-negative after resize, got %d", m.viewport.Height)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(model)

	_ = m.View()
	_ = m.viewport.View()
}

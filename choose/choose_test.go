package choose

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestToggleWithSpace ensures that pressing space toggles the highlighted
// item in a multi-select choose prompt. Regression test for #1143: the
// keymap previously bound the literal " " rune, but bubbletea v2's
// KeyPressMsg.String() reports the spacebar as "space".
func TestToggleWithSpace(t *testing.T) {
	km := defaultKeymap()
	km.Toggle.SetEnabled(true)

	m := model{
		items:  []item{{text: "one"}, {text: "two"}},
		limit:  2,
		keymap: km,
	}

	got, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	updated := got.(model)

	if !updated.items[0].selected {
		t.Fatalf("expected item to be toggled on after pressing space")
	}
	if updated.numSelected != 1 {
		t.Fatalf("expected numSelected to be 1, got %d", updated.numSelected)
	}
}

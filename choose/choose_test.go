package choose

import (
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func TestChooseToggleKeybinding(t *testing.T) {
	km := defaultKeymap()
	km.Toggle.SetEnabled(true)

	testKeys := []struct {
		name string
		msg  tea.KeyPressMsg
	}{
		{
			name: "space key",
			msg:  tea.KeyPressMsg{Code: tea.KeySpace},
		},
		{
			name: "literal space character",
			msg:  tea.KeyPressMsg{Code: ' ', Text: " "},
		},
		{
			name: "x key",
			msg:  tea.KeyPressMsg{Code: 'x', Text: "x"},
		},
		{
			name: "tab key",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab},
		},
	}

	for _, tc := range testKeys {
		t.Run(tc.name, func(t *testing.T) {
			if !key.Matches(tc.msg, km.Toggle) {
				t.Errorf("expected %s (%v) to match Toggle binding", tc.name, tc.msg)
			}
		})
	}
}

func TestChooseToggleModel(t *testing.T) {
	km := defaultKeymap()
	km.Toggle.SetEnabled(true)

	m := model{
		limit: 2,
		items: []item{
			{text: "one"},
			{text: "two"},
		},
		keymap: km,
	}

	spaceMsg := tea.KeyPressMsg{Code: tea.KeySpace}

	// Press space to toggle first item
	newModel, _ := m.Update(spaceMsg)
	m = newModel.(model)
	if !m.items[0].selected {
		t.Errorf("expected item 0 to be selected after pressing space")
	}
	if m.numSelected != 1 {
		t.Errorf("expected numSelected to be 1, got %d", m.numSelected)
	}

	// Press space again to deselect first item
	newModel, _ = m.Update(spaceMsg)
	m = newModel.(model)
	if m.items[0].selected {
		t.Errorf("expected item 0 to be deselected after pressing space again")
	}
	if m.numSelected != 0 {
		t.Errorf("expected numSelected to be 0, got %d", m.numSelected)
	}
}

package choose

import (
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
)

// buildModel builds a model with multi-select enabled
// (limit > 1 enables the Toggle binding, mirroring command.go).
func buildModel(t *testing.T, items, limit int) model {
	t.Helper()
	pager := paginator.New()
	pager.SetTotalPages((items + 4) / 5)
	pager.PerPage = 5
	ms := make([]item, items)
	for i := range ms {
		ms[i] = item{text: string(rune('a' + i))}
	}
	km := defaultKeymap()
	if limit > 1 {
		km.Toggle.SetEnabled(true)
	}
	return model{
		height:     5,
		items:      ms,
		index:      0,
		limit:      limit,
		paginator:  pager,
		help:       help.New(),
		keymap:     km,
	}
}

func TestSpaceTogglesItem(t *testing.T) {
	m := buildModel(t, 3, 5)
	if m.items[m.index].selected {
		t.Fatal("item already selected before the space press")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.items[m.index].selected {
		t.Fatal("space does not toggle the highlighted item in multi-select")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.items[m.index].selected {
		t.Fatal("space should toggle the item back off")
	}
}

func TestOtherToggleKeysStillWork(t *testing.T) {
	for _, keystroke := range []struct {
		name string
		code rune
	}{
		{"x", 'x'},
		{"tab", tea.KeyTab},
	} {
		t.Run(keystroke.name, func(t *testing.T) {
			m := buildModel(t, 2, 5)
			m.Update(tea.KeyPressMsg{Code: keystroke.code})
			if !m.items[m.index].selected {
				t.Fatalf("%s should toggle the highlighted item", keystroke.name)
			}
		})
	}
}

func TestSpaceDoesNotToggleInSingleSelect(t *testing.T) {
	m := buildModel(t, 3, 1)
	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.items[m.index].selected {
		t.Fatal("space must not select in single-select (limit == 1 guard)")
	}
}
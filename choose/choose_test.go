package choose

import (
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
)

// buildModel собирает модель с включённым мульти-выбором
// (limit > 1 включает привязку Toggle, как в command.go).
func buildModel(t *testing.T, items int, limit int) model {
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
		t.Fatal("туccа: пункт уже выбран до нажатия")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.items[m.index].selected {
		t.Fatal("регрессия: Space не переключает пункт в мульти-выборе")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.items[m.index].selected {
		t.Fatal("Space должен переключать туда-обратно (unselect)")
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
				t.Fatalf("%s должен переключать пункт", keystroke.name)
			}
		})
	}
}

func TestSpaceDoesNotToggleInSingleSelect(t *testing.T) {
	m := buildModel(t, 3, 1)
	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.items[m.index].selected {
		t.Fatal("в single-select Space не должен выбирать (guard limit==1)")
	}
}
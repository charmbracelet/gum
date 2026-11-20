package table

import (
	"strings"
	"testing"

	bubbles "charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

func TestWindowSizeSetsTableWidth(t *testing.T) {
	m := model{
		table: bubbles.New(
			bubbles.WithColumns([]bubbles.Column{{Title: "NAME", Width: 10}}),
			bubbles.WithRows([]bubbles.Row{{"Alpha"}}),
			bubbles.WithFocused(true),
		),
		padding: []int{1, 2, 1, 3},
	}

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	got := updated.(model).table.Width()
	if want := 75; got != want {
		t.Fatalf("table width = %d, want %d", got, want)
	}
	if view := updated.(model).View().Content; !strings.Contains(view, "Alpha") {
		t.Fatalf("table view does not contain row: %q", view)
	}
}

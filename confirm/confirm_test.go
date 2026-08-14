package confirm

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWindowSizeMsgClearsScreen(t *testing.T) {
	m := model{
		prompt:      "do you confirm?",
		affirmative: "yes",
		negative:    "no",
		padding:     []int{0, 0, 0, 0},
	}

	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	next := updated.(model)

	if next.width != 80 {
		t.Fatalf("expected window width to be stored, got %d", next.width)
	}
	if cmd == nil {
		t.Fatal("expected ClearScreen command on resize")
	}
	if reflect.ValueOf(cmd).Pointer() != reflect.ValueOf(tea.ClearScreen).Pointer() {
		t.Fatalf("expected tea.ClearScreen command, got %#v", cmd)
	}
}

func TestContentWidthAccountsForPadding(t *testing.T) {
	m := model{
		width:   100,
		padding: []int{1, 2, 3, 4},
	}
	if got := m.contentWidth(); got != 94 {
		t.Fatalf("expected content width 94, got %d", got)
	}
}

package confirm

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestSubmitKeys(t *testing.T) {
	tests := map[string]tea.KeyPressMsg{
		"enter": {Code: tea.KeyEnter},
		// An Enter typed before the terminal is switched to raw mode is
		// translated to a line feed by the tty (ICRNL) and arrives as ctrl+j.
		"ctrl+j": {Code: 'j', Mod: tea.ModCtrl},
	}
	for name, msg := range tests {
		t.Run(name, func(t *testing.T) {
			m := model{
				affirmative:  "Yes",
				negative:     "No",
				confirmation: true,
				keys:         defaultKeymap("Yes", "No"),
			}
			if got := msg.String(); got != name {
				t.Fatalf("key string = %q, want %q", got, name)
			}
			tm, cmd := m.Update(msg)
			if cmd == nil {
				t.Fatal("expected a quit command")
			}
			got := tm.(model)
			if !got.quitting {
				t.Error("expected the prompt to quit")
			}
			if !got.confirmation {
				t.Error("expected the default (affirmative) selection to be submitted")
			}
		})
	}
}

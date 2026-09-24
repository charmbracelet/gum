package cursor

import "testing"

func TestModes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"blink exists", "blink", true},
		{"hide exists", "hide", true},
		{"static exists", "static", true},
		{"unknown does not exist", "unknown", false},
	}

	if len(Modes) != 3 {
		t.Fatalf("expected 3 modes, got %d", len(Modes))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := Modes[tt.key]
			if ok != tt.expected {
				t.Errorf("Modes[%q] expected existence=%v, got %v", tt.key, tt.expected, ok)
			}
		})
	}
}

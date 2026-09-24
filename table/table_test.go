package table

import "testing"

func TestNumLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{5, 1},
		{9, 1},
		{10, 2},
		{50, 2},
		{99, 2},
		{100, 3},
		{500, 3},
		{999, 3},
		{1000, 4},
		{9999, 4},
		{10000, 5},
		{-1, 1},
		{-10, 2},
		{-100, 3},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := numLen(tt.input)
			if got != tt.expected {
				t.Errorf("numLen(%d) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

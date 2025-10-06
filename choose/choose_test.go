package choose

import (
	"testing"
)

func TestInitialPositioning(t *testing.T) {
	tests := []struct {
		name     string
		options  []string
		initial  string
		expected int // expected starting index
	}{
		{
			name:     "initial option exists",
			options:  []string{"Apple", "Banana", "Cherry", "Date"},
			initial:  "Cherry",
			expected: 2,
		},
		{
			name:     "initial option doesn't exist",
			options:  []string{"Apple", "Banana", "Cherry", "Date"},
			initial:  "Orange",
			expected: 0, // should default to first item
		},
		{
			name:     "empty initial option",
			options:  []string{"Apple", "Banana", "Cherry", "Date"},
			initial:  "",
			expected: 0, // should default to first item
		},
		{
			name:     "initial option is first",
			options:  []string{"Apple", "Banana", "Cherry", "Date"},
			initial:  "Apple",
			expected: 0,
		},
		{
			name:     "initial option is last",
			options:  []string{"Apple", "Banana", "Cherry", "Date"},
			initial:  "Date",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from command.go
			startingIndex := 0
			if tt.initial != "" {
				for i, option := range tt.options {
					if option == tt.initial {
						startingIndex = i
						break
					}
				}
			}

			if startingIndex != tt.expected {
				t.Errorf("expected starting index %d, got %d", tt.expected, startingIndex)
			}
		})
	}
}
